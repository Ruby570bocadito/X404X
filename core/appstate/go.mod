module github.com/ruby570bocadito/x404x/core/appstate

go 1.23

require (
	github.com/ruby570bocadito/x404x/core/agent v0.0.0
	github.com/ruby570bocadito/x404x/core/orchestrator v0.0.0
	github.com/ruby570bocadito/x404x/shared/config v0.0.0
	github.com/ruby570bocadito/x404x/shared/logger v0.0.0
	github.com/ruby570bocadito/x404x/shared/types v0.0.0
)

replace (
	github.com/ruby570bocadito/x404x/core/agent => ../agent
	github.com/ruby570bocadito/x404x/core/orchestrator => ../orchestrator
	github.com/ruby570bocadito/x404x/shared/config => ../../shared/config
	github.com/ruby570bocadito/x404x/shared/logger => ../../shared/logger
	github.com/ruby570bocadito/x404x/shared/types => ../../shared/types
)
