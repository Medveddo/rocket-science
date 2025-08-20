package migrator

type Migrator interface {
	Up() error
	Down() error
	Status() error
	Version() error
}
