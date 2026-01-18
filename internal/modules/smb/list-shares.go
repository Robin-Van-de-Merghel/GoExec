package smb

import (
	"fmt"

	"github.com/GoExec/pkg/modules"
	"github.com/mvo5/libsmbclient-go"
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
	Labels:          []string{"SMB"},
}

type ModuleInput struct {
	// Target, either an IP or a Host
	Target modules.ModuleTarget

	// Credentials
	// TODO: Change for generic
	Credentials modules.Credentials
}

type Module struct {
	Input ModuleInput
}

func (m *Module) Configure(input any) (error, string) {
	mi, ok := input.(ModuleInput)
	if !ok {
		return fmt.Errorf("invalid input type"), ""
	}
	m.Input = mi
	return nil, "configured"
}

func (m *Module) Run() (error, string) {
	client := libsmbclient.New()
	err := modules.SetupSMBAuth(client, m.Input.Credentials)
	if err != nil {
		return err, ""
	}

	dh, err := client.Opendir(fmt.Sprintf("smb://%s", m.Input.Target.Host))
	if err != nil {
		return err, ""
	}

	defer dh.Close()
	for {
		dirent, err := dh.Readdir()
		if err != nil {
			break
		}
		fmt.Println(dirent)
	}

	return nil, "Success"
}
