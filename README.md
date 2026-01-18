# Introduction

This project aims at reproducing [@Pennyw0rth](https://github.com/Pennyw0rth)'s NetExec project in Go.

It acts for me as a personal training in go as well as understanding how AD vulnerabilities works, by coding it.

# First step on GoExec

GoExec is made to be modular as well as easy to use. As later described in [Modules](#Modules), GoExec uses tags. It helps sorting modules.

The global syntax is the following:

```bash
./GoExec [tag] -M [module_name] <flags>
```

For example, to launch the `list-shares` module, you can:

```bash
./GoExec smb -M list-shares <flags>
```

> [!NOTE]
> It's close to what does NetExec, but the host is moved into a flag.

We said above that we use tags. You can list all tags using the `-T` flag:

```bash
./GoExec -T
```

To list all modules, there are two ways: either you know which tag you want to search for, either you don't:

```bash
./GoExec -L // List all modules

./GoExec smb -L // List all modules with the "smb" tag
```

To get the options of a module, you can use `--options`:

```bash
➜  GoExec git:(master) ✗ ./GoExec  smb -M list-shares --options
Module: list-shares
Description: List-shares is a SMB module aiming at getting first information about SMB shares such as: open shares, rights, and more.
Tags: SMB
Arguments:
 - Host:  (string)
 - Domain: Domain of the user (string)
 - Username: Username to login with (string)
 - Password: Password to login with (string)
 - Kerberos: Use kerberos auth, requires KRB5CCNAME env variable (bool)
```

And finally to launch a module:

```bash
➜  GoExec git:(main) ✗ ./GoExec smb -M list-shares --Host localhost --Username admin --Password admin 
2026/01/15 12:55:30 INFO Starting module module=list-shares
...
2026/01/15 12:55:30 INFO Module executed successfully module=list-shares duration=1.070833ms result=Success
```

# Design

## Modules

Modules are at the center of GoExec. They must follow some strict rules, and export the necessary for the user to play with.

```go
type Module interface {
	// Configure a module
	Configure(any) (error, string)
	// Run a module
	Run() (error, string)
}

type ModuleMetadata struct {
	// MUST be unique to avoid collisions
	UniqueName string
	// Present the module, to display while listing modules
	PresentMessages string
	// Labels to filter modules (e.g. "smb")
	Labels []string
}
```

Labels are a way to sort modules in categories. Where NetExec enforces users to chose protocols, GoExec is more flexible and allows a more granular selection. For example:

```bash
// (these examples does not exist)
./GoExec shares -L // List all shares-related modules
./GoExec RCE -L // List all RCEs
```

## Arguments

By using reflection, GoExec allows using structures as arguments before make it flat. For example:

```go
type Credentials struct {
	UsernamePassAuth
	KerberosAuth
}

// Classic username/password
type UsernamePassAuth struct {
	Domain   string
	Username string
	Password string
}

// Kerberos Auth
type KerberosAuth struct {
	Kerberos bool
}
```

GoExec in [flags.go](./internal/core/flags.go) will extend all structs, and provide the following flags: `--Username`, `--Password`, `--Kerberos`.
