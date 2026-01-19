package smb

import (
	"fmt"
	"log/slog"

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
	Targets modules.Targets

	// Credentials
	// TODO: Change for generic
	Credentials modules.Credentials
}

type Module struct {
	Input ModuleInput
}

func (m *Module) Configure(input any) error {
	mi, ok := input.(ModuleInput)
	if !ok {
		return fmt.Errorf("invalid input type")
	}
	m.Input = mi
	return nil
}

func (m *Module) Run() error {
	targets, err := modules.GetTargets(m.Input.Targets)
	if err != nil {
		return err
	}

	for _, el := range targets {
		err := m.RunOnce(el)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: Refactor to avoid having access to the targets from a Module object
func (m *Module) RunOnce(target modules.ModuleTarget) error {
	client := libsmbclient.New()
	err := modules.SetupSMBAuth(client, m.Input.Credentials)
	if err != nil {
		return err
	}

	dh, err := client.Opendir(
		fmt.Sprintf("smb://%s/", target.Host),
	)
	if err != nil {
		return err
	}
	defer dh.Closedir()

	for {
		dirent, err := dh.Readdir()
		if err != nil {
			break
		}
		slog.Debug("Found a share", "share_name", dirent.Name, "share_comment", dirent.Comment)
	}

	return nil
}
