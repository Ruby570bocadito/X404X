module github.com/ruby570bocadito/x404x

go 1.24

require (
	github.com/charmbracelet/bubbletea v1.2.4
	github.com/charmbracelet/lipgloss v1.0.0
	github.com/ruby570bocadito/x404x/internal/agent v0.0.0
	github.com/ruby570bocadito/x404x/internal/api v0.0.0
	github.com/ruby570bocadito/x404x/internal/appstate v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/config v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/logger v0.0.0
	github.com/ruby570bocadito/x404x/pkg/shared/types v0.0.0
	github.com/spf13/cobra v1.8.1
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/x/ansi v0.4.5 // indirect
	github.com/charmbracelet/x/term v0.2.1 // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/sync v0.9.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/ruby570bocadito/x404x/internal/agent => ../../core/agent
	github.com/ruby570bocadito/x404x/internal/api => ../../core/api
	github.com/ruby570bocadito/x404x/internal/appstate => ../../core/appstate
	github.com/ruby570bocadito/x404x/internal/orchestrator => ../../core/orchestrator
	github.com/ruby570bocadito/x404x/pkg/shared/config => ../../shared/config
	github.com/ruby570bocadito/x404x/pkg/shared/logger => ../../shared/logger
	github.com/ruby570bocadito/x404x/pkg/shared/types => ../../shared/types
)
