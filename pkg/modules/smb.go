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
