package config

import (
	"github.com/ipfs-search/ipfs-search/components/sniffer"
	"time"
)

// Sniffer is configuration pertaining to the sniffer
type Sniffer struct {
	LastSeenExpiration time.Duration `yaml:"lastseen_expiration"`
	LastSeenPruneLen   int           `yaml:"lastseen_prunelen"`
	LoggerTimeout      time.Duration `yaml:"logger_timeout"`
	BufferSize         uint          `yaml:"buffer_size"`
}

// SnifferConfig returns component-specific configuration from the canonical central configuration.
func (c *Config) SnifferConfig() *sniffer.Config {
	cfg := sniffer.Config(c.Sniffer)
	return &cfg
}

// SnifferDefaults returns the defaults for component configuration, based on the component-specific configuration.
func SnifferDefaults() Sniffer {
	return Sniffer(*sniffer.DefaultConfig())
}
