package ransomware

// config_default.go — DefaultRansomwareConfig returns a config populated with
// safe defaults. It was referenced by the test suite but was removed during a
// cleanup commit; recreated here for the audit.

// DefaultRansomwareConfig returns a RansomwareConfig with production-safe
// defaults (simulation disabled, non-empty network targets) as expected by the
// test suite (cloud_id_test, destruction_test, raas_test, c2_channels_test).
func DefaultRansomwareConfig() *RansomwareConfig {
	cfg := NewScannerConfig()

	// Cloud profiles expected by TestIMDSv2BypassConfig.
	cfg.AWSProfile = "default"
	cfg.AzureProfile = "default"
	cfg.GCPProject = "default"

	// Network targets expected by TestHardwareKillConfig.
	cfg.TargetSubnet = "10.0.0.0/24"
	cfg.C2Endpoint = "127.0.0.1:8443"
	cfg.GatewayIP = "10.0.0.1"

	return cfg
}
