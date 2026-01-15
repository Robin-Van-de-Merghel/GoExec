package modules

type Module interface {
	// Configure a module
	Configure(any) (error, string)
	// Run a module
	Run() (error, string)
}

type ModuleMetadata struct {
	// MUST be unique to avoid collisions
	UniqueName string
	// Present the module, to display while listing modules
	PresentMessages string
	// Labels to filter modules (e.g. "smb")
	Labels []string
}
