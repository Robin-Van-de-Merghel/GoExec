package modules

/*

Credentials in an active directory can be multiple things:
1. Username/Password (+ guest <=> empty)
2. Kerberos
3. Hash

This module supports all of them.
*/

// TODO: Implement the choice
type Credentials interface{}

// Classic username/password
type UsernamePassAuth struct {
	Domain   string `help:"Domain of the user"`
	Username string `help:"Username to login with"`
	Password string `help:"Password to login with"`
}

// TODO: Kerberos
// type KerberosAuth struct {}

// TODO: Hash
// type HashAuth struct {}
