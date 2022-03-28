package migration

type Command interface {
	Up() error
	UpTo(limit uint) error
	Down() error
	DownTo(limit int) error
	Drop() error
	Refresh() error

	GoTo(version uint) error
	Force(version uint) error
	CurrentVersion() (uint, bool, error)
	PrintUsageInfo()
}
