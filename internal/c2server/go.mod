module github.com/ruby570bocadito/x404x/internal/c2server

go 1.24

require (
	github.com/ruby570bocadito/x404x/internal/appstate v0.0.0
	github.com/ruby570bocadito/x404x/pkg/proto/gen/agent v0.0.0
	github.com/ruby570bocadito/x404x/pkg/proto/gen/c2 v0.0.0
	github.com/ruby570bocadito/x404x/pkg/proto/gen/common v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/config v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/logger v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/types v0.0.0
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.36.11
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/ruby570bocadito/x404x/internal/appstate => ../appstate
	github.com/ruby570bocadito/x404x/internal/crypto => ../crypto
	github.com/ruby570bocadito/x404x/pkg/proto/gen/agent => ../../pkg/proto/gen/agent
	github.com/ruby570bocadito/x404x/pkg/proto/gen/c2 => ../../pkg/proto/gen/c2
	github.com/ruby570bocadito/x404x/pkg/proto/gen/common => ../../pkg/proto/gen/common
	github.com/ruby570bocadito/x404x/pkg/shared/config => ../../pkg/shared/config
	github.com/ruby570bocadito/x404x/pkg/shared/logger => ../../pkg/shared/logger
	github.com/ruby570bocadito/x404x/pkg/shared/types => ../../pkg/shared/types
)
