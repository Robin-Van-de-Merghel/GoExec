package listshares

import (
	"fmt"
	"log/slog"

	"github.com/GoExec/pkg/modules"
	libsmbclient "github.com/robin-van-de-merghel/libsmbclient-go/pkg/bindings"
)

func (m *ListSharesModule) Run() ([]any, error) {
	targets, err := modules.GetTargets(m.Input.Targets)
	if err != nil {
		return nil, err
	}

	for _, el := range targets {
		_, err := m.RunOnce(el)
		if err != nil {
			return nil, err
		}
	}

	// FIXME: Return Share object
	return nil, nil
}

// TODO: Refactor to avoid having access to the targets from a Module object
func (m *ListSharesModule) RunOnce(target modules.ModuleTarget) (any, error) {
	client := libsmbclient.New()
	err := modules.SetupSMBAuth(client, m.Input.Credentials)
	if err != nil {
		return nil, err
	}

	dh, err := client.Opendir(
		fmt.Sprintf("smb://%s/", target.Host),
	)
	if err != nil {
		return nil, err
	}
	defer dh.Closedir()

	for {
		dirent, err := dh.Readdir()
		if err != nil {
			break
		}
		slog.Debug("Found a share", "share_name", dirent.Name, "share_comment", dirent.Comment)
	}

	// FIXME: Return Share object
	return nil, nil
}
