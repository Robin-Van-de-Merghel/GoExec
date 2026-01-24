package core

import (
	listshares "github.com/GoExec/internal/modules/smb/list-shares"
	"github.com/GoExec/pkg/modules"
)

type ModuleEntry struct {
	Metadata modules.ModuleMetadata
	Factory  func() modules.Module
}

var AllModules = []ModuleEntry{
	{
		Metadata: listshares.Metadata,
		Factory: func() modules.Module {
			return &listshares.ListSharesModule{}
		},
	},
}
