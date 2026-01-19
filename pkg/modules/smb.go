package modules

import (
	"github.com/mvo5/libsmbclient-go"
)

var (
	smbDomain   string
	smbUsername string
	smbPassword string
)

func authCallbackFunc(serverName, shareName string) (string, string, string) {
	return smbDomain, smbUsername, smbPassword
}

func SetupSMBAuth(client *libsmbclient.Client, creds Credentials) error {
	hasKerberos, err := IsKerberosAvailale(creds)
	if err != nil {
		return err
	}
	if hasKerberos {
		client.SetUseKerberos()
	} else {
		smbDomain = creds.Domain
		smbUsername = creds.Username
		smbPassword = creds.Password
		client.SetAuthCallback(authCallbackFunc)
	}
	return nil
}
