module github.com/ruby570bocadito/x404x/internal/dispatch
go 1.25.0
require (
	github.com/ruby570bocadito/x404x/internal/registry v0.0.0
	github.com/ruby570bocadito/x404x/shared/types v0.0.0
)
replace (
	github.com/ruby570bocadito/x404x/internal/registry => ../registry
	github.com/ruby570bocadito/x404x/shared/types => ../../shared/types
)
