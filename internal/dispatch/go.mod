module github.com/ruby570bocadito/x404x/internal/dispatch
go 1.24
require (
	github.com/ruby570bocadito/x404x/internal/registry v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/types v0.0.0
)
replace (
	github.com/ruby570bocadito/x404x/internal/registry => ../registry
	github.com/ruby570bocadito/x404x/pkg/shared/types => ../../pkg/shared/types
)
