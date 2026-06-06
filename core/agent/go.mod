module github.com/ruby570bocadito/x404x/core/agent

go 1.25.0

require (
	github.com/ruby570bocadito/x404x/core/crypto v0.0.0
	github.com/ruby570bocadito/x404x/core/proto/gen/agent v0.0.0
	github.com/ruby570bocadito/x404x/core/ransomware v0.0.0
	github.com/ruby570bocadito/x404x/shared/config v0.0.0
	github.com/ruby570bocadito/x404x/shared/logger v0.0.0
	github.com/ruby570bocadito/x404x/shared/types v0.0.0
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/ruby570bocadito/x404x/core/crypto => ../crypto
	github.com/ruby570bocadito/x404x/core/proto/gen/agent => ../proto/gen/agent
	github.com/ruby570bocadito/x404x/core/ransomware => ../ransomware
	github.com/ruby570bocadito/x404x/shared/config => ../../shared/config
	github.com/ruby570bocadito/x404x/shared/logger => ../../shared/logger
	github.com/ruby570bocadito/x404x/shared/types => ../../shared/types
)
