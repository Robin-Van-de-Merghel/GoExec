package modules

type Module interface {
	// Configure a module
	Configure(any) error
	// Run a module n-th time (n = Target number)
	// TODO: Later make an "algorithm" function
	Run() ([]any, error)
	// Run a module once
	RunOnce(ModuleTarget) (any, error)
}

type ModuleMetadata struct {
	// MUST be unique to avoid collisions
	UniqueName string
	// Present the module, to display while listing modules
	PresentMessages string
	// Labels to filter modules (e.g. "smb")
	Labels []string
}
