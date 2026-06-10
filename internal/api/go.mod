module github.com/ruby570bocadito/x404x/internal/api

go 1.24

require (
	github.com/gorilla/websocket v1.5.3
	github.com/ruby570bocadito/x404x/internal/appstate v0.0.0
	github.com/ruby570bocadito/x404x/internal/orchestrator v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/config v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/logger v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/types v0.0.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/ruby570bocadito/x404x/internal/appstate => ../appstate
	github.com/ruby570bocadito/x404x/internal/orchestrator => ../orchestrator
	github.com/ruby570bocadito/x404x/pkg/shared/config => ../../shared/config
	github.com/ruby570bocadito/x404x/pkg/shared/logger => ../../shared/logger
	github.com/ruby570bocadito/x404x/pkg/shared/types => ../../shared/types
)
