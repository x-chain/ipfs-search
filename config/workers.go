package config

// import (
// 	"time"
// )

type Workers struct {
	// RetryWait       time.Duration     `yaml:"retry_wait"` // Time to wait between failed http requests.
	// StartupDelay     time.Duration `yaml:"startup_delay"`
	HashWorkers      int `yaml:"hash_workers" env:"HASH_WORKERS"`
	FileWorkers      int `yaml:"file_workers" env:"FILE_WORKERS"`
	DirectoryWorkers int `yaml:"directory_workers" env:"DIRECTORY_WORKERS"`
}

func WorkersDefaults() Workers {
	return Workers{
		// StartupDelay:     time.Duration(100 * time.Millisecond),
		HashWorkers:      70,
		FileWorkers:      120,
		DirectoryWorkers: 70,
	}
}