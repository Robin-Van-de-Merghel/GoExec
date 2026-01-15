package core

import (
	"github.com/GoExec/internal/modules/smb"
	"github.com/GoExec/pkg/modules"
)

type ModuleEntry struct {
	Metadata modules.ModuleMetadata
	Factory  func() modules.Module
}

var AllModules = []ModuleEntry{
	{
		Metadata: smb.Metadata,
		Factory: func() modules.Module {
			return &smb.Module{}
		},
	},
}
