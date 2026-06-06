module github.com/ruby570bocadito/x404x/core/orchestrator

go 1.22

require (
	github.com/ruby570bocadito/x404x/shared/config v0.0.0
	github.com/ruby570bocadito/x404x/shared/logger v0.0.0
	github.com/ruby570bocadito/x404x/shared/types v0.0.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/ruby570bocadito/x404x/shared/config => ../../shared/config
	github.com/ruby570bocadito/x404x/shared/logger => ../../shared/logger
	github.com/ruby570bocadito/x404x/shared/types => ../../shared/types
)
