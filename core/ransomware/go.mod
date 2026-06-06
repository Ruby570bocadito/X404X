module github.com/ruby570bocadito/x404x/core/ransomware

go 1.25.0

require (
	github.com/ruby570bocadito/x404x/core/crypto v0.0.0
	github.com/ruby570bocadito/x404x/shared/types v0.0.0
	golang.org/x/crypto v0.48.0
)

require golang.org/x/sys v0.42.0 // indirect

replace (
	github.com/ruby570bocadito/x404x/core/crypto => ../crypto
	github.com/ruby570bocadito/x404x/shared/types => ../../shared/types
)
