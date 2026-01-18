package modules

import (
	"fmt"
	"os"
)

/*

Credentials in an active directory can be multiple things:
1. Username/Password (+ guest <=> empty)
2. Kerberos
3. Hash

This module supports all of them.
*/

type Credentials struct {
	UsernamePassAuth
	KerberosAuth
}

// Classic username/password
type UsernamePassAuth struct {
	Domain   string `help:"Domain of the user"`
	Username string `help:"Username to login with"`
	Password string `help:"Password to login with"`
}

// Kerberos Auth
type KerberosAuth struct {
	Kerberos bool `help:"Use kerberos auth, requires KRB5CCNAME env variable"`
}

// TODO: Hash
// type HashAuth struct {}

// IsKerberosAvailale checks if Kerberos is asked, if so the env var MUST be defined
func IsKerberosAvailale(creds Credentials) (bool, error) {
	if creds.Kerberos {
		// For kerberos, KRB5CCNAME is required in the environment variables
		if os.Getenv("KRB5CCNAME") == "" {
			return false, fmt.Errorf("To use kerberos, you need to define KRB5CCNAME")
		}

		return true, nil
	}

	return false, nil
}
