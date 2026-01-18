package modules

import (
	"github.com/mvo5/libsmbclient-go"
)

// Callback for static auth
func SMBAuth(domain, username, password string) libsmbclient.AuthCallback {
	return func(serverName, shareName string) (domain, username, password string) {
		return domain, username, password
	}
}

// SetupSMBAuth setups either kerberos or user/pass depending on the flag
func SetupSMBAuth(client *libsmbclient.Client, creds Credentials) error {
	hasKerberos, err := IsKerberosAvailale(creds)
	if err != nil {
		return err
	}

	if hasKerberos {
		client.SetUseKerberos()
	} else {
		// FIXME: Not working with empty password with SMB server
		client.SetAuthCallback(SMBAuth(
			creds.Domain, creds.Username, creds.Password,
		))
	}

	return nil
}
