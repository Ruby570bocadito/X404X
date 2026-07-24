module github.com/ruby570bocadito/x404x/internal/ransomware

go 1.25.0

require (
	github.com/ruby570bocadito/x404x/internal/crypto v0.0.0
	golang.org/x/crypto v0.48.0
)

require golang.org/x/sys v0.42.0 // indirect

replace (
	github.com/ruby570bocadito/x404x/internal/crypto => ../crypto
	github.com/ruby570bocadito/x404x/pkg/shared/types => ../../pkg/shared/types
)
