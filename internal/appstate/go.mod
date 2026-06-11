module github.com/ruby570bocadito/x404x/internal/appstate

go 1.24

require (
	github.com/ruby570bocadito/x404x/internal/agent v0.0.0
	github.com/ruby570bocadito/x404x/internal/orchestrator v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/config v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/logger v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/types v0.0.0
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/sys v0.26.0 // indirect
	modernc.org/libc v1.72.3 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.51.0 // indirect
)

replace (
	github.com/ruby570bocadito/x404x/internal/agent => ../agent
	github.com/ruby570bocadito/x404x/internal/orchestrator => ../orchestrator
	github.com/ruby570bocadito/x404x/pkg/shared/config => ../../pkg/shared/config
	github.com/ruby570bocadito/x404x/pkg/shared/logger => ../../pkg/shared/logger
	github.com/ruby570bocadito/x404x/pkg/shared/types => ../../pkg/shared/types
)
require github.com/ruby570bocadito/x404x/internal/dispatch v0.0.0
replace github.com/ruby570bocadito/x404x/internal/dispatch => ../../internal/dispatch
