package drivers

import (
	"github.com/k1LoW/tbls/schema"
)

// Driver is the common interface for database drivers.
type Driver interface {
	Analyze(*schema.Schema) error
	Info() (*schema.Driver, error)
}

// ConfigurableDriver is the interface for drivers that can accept configuration.
type ConfigurableDriver interface {
	Driver
	SetLogicalNameConfig(delimiter string, fallbackToName bool)
}

// Option is the type for change Config.
type Option func(Driver) error
