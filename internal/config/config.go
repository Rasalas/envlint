package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

const defaultConfigFile = ".envlint.toml"

// Config represents the .envlint.toml configuration.
type Config struct {
	Example  string   `toml:"example"`
	EnvFiles []string `toml:"envFiles"`
	Rules    Rules    `toml:"rules"`
}

// Rules holds validation rule settings.
type Rules struct {
	RequireAll  bool     `toml:"requireAll"`
	NoExtra     bool     `toml:"noExtra"`
	StrictURLs  bool     `toml:"strictUrls"`
	StrictPorts bool     `toml:"strictPorts"`
	Required    KeyList  `toml:"required"`
	Ignore      KeyList  `toml:"ignore"`
}

// KeyList holds a list of key names.
type KeyList struct {
	Keys []string `toml:"keys"`
}

// Default returns the default configuration.
func Default() Config {
	return Config{
		Example:  ".env.example",
		EnvFiles: []string{".env"},
		Rules: Rules{
			StrictURLs:  true,
			StrictPorts: true,
		},
	}
}

// Load reads .envlint.toml from the current directory.
func Load() (Config, error) {
	return LoadFrom(defaultConfigFile)
}

// LoadFrom reads config from a specific path.
func LoadFrom(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
