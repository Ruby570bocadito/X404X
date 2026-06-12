module github.com/ruby570bocadito/x404x/internal/agent

go 1.24.0

require (
	github.com/ruby570bocadito/x404x/internal/crypto v0.0.0
	github.com/ruby570bocadito/x404x/internal/ransomware v0.0.0
	github.com/ruby570bocadito/x404x/pkg/proto/gen/agent v0.0.0
	github.com/ruby570bocadito/x404x/pkg/proto/gen/bridge v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/config v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/logger v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/types v0.0.0
	golang.org/x/crypto v0.48.0
	google.golang.org/grpc v1.64.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/ruby570bocadito/x404x/internal/crypto => ../crypto
	github.com/ruby570bocadito/x404x/internal/ransomware => ../ransomware
	github.com/ruby570bocadito/x404x/pkg/proto/gen/agent => ../../pkg/proto/gen/agent
	github.com/ruby570bocadito/x404x/pkg/proto/gen/bridge => ../../pkg/proto/gen/bridge
	github.com/ruby570bocadito/x404x/pkg/shared/config => ../../pkg/shared/config
	github.com/ruby570bocadito/x404x/pkg/shared/logger => ../../pkg/shared/logger
	github.com/ruby570bocadito/x404x/pkg/shared/types => ../../pkg/shared/types
)
