package listshares

import (
	"fmt"

	"github.com/GoExec/pkg/modules"
)

/*

Goals:

1. Connect to a SMB server
2. List shares
	a. Name
	b. Read/Write
	c. Content and / or stats

*/

var Metadata = modules.ModuleMetadata{
	UniqueName:      "list-shares",
	PresentMessages: `List-shares is a SMB module aiming at getting first information about SMB shares such as: open shares, rights, and more.`,
	Labels:          []string{"SMB", "Shares", "Low-Privilege"},
}

type ModuleInput struct {
	// Target, either an IP or a Host
	Targets modules.Targets

	// Credentials
	// TODO: Change for generic
	Credentials modules.Credentials
}

type ListSharesModule struct {
	Input ModuleInput
}

type ShareRights uint8

type Share struct {
	Name        string
	Host        string
	ShareRights ShareRights
}

func (m *ListSharesModule) Configure(input any) error {
	mi, ok := input.(ModuleInput)
	if !ok {
		return fmt.Errorf("invalid input type")
	}
	m.Input = mi
	return nil
}

func (sr ShareRights) String() string {
	switch sr {
	default:
		return ""
	}
}
