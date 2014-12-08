package migration

import (
	"errors"
	"github.com/elwinar/rambler/configuration"
)

// Driver is the interface used by the program to interact with the migration
// table in database
type Driver interface {
	MigrationTableExists() (bool, error)
	CreateMigrationTable() error
	ListAppliedMigrations() ([]uint64, error)
	ListAvailableMigrations() ([]uint64, error)
}

// DriverConstructor is the function type used to create drivers
type DriverConstructor func(configuration.Environment) (Driver, error)

// The various errors returned by the package
var (
	ErrDriverAlreadyRegistered = errors.New("driver already registered")
	ErrDriverNotRegistered     = errors.New("driver not registered")
)

var (
	constructors map[string]DriverConstructor
)

func init() {
	constructors = make(map[string]DriverConstructor)
}

// Register register a constructor for a driver
func RegisterDriver(name string, constructor DriverConstructor) error {
	return registerDriver(name, constructor, constructors)
}

func registerDriver(name string, constructor DriverConstructor, constructors map[string]DriverConstructor) error {
	if _, found := constructors[name]; found {
		return ErrDriverAlreadyRegistered
	}

	constructors[name] = constructor
	return nil
}

// Get initialize a driver from the given name and options
func GetDriver(env configuration.Environment) (Driver, error) {
	return getDriver(env, constructors)
}

func getDriver(env configuration.Environment, constructors map[string]DriverConstructor) (Driver, error) {
	constructor, found := constructors[env.Driver]
	if !found {
		return nil, ErrDriverNotRegistered
	}

	return constructor(env)
}
