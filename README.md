# Introduction

This project aims at reproducing [@Pennyw0rth](https://github.com/Pennyw0rth)'s NetExec project in Go.

# Examples

These examples are here to provide you a small doc.

Listing tags:

```bash
➜  GoExec git:(main) ✗ ./GoExec -T # or --list-tags                                                                             
Available tags:
  - smb
```

Listing modules from a specific tag:

```bash
➜  GoExec git:(main) ✗ ./GoExec smb -L # or --list-modules
Matching modules:
  - list-shares (tags: SMB)
```

Listing all modules:

```bash
➜  GoExec git:(main) ✗ ./GoExec -L    
Matching modules:
  - list-shares (tags: SMB)
```

Getting options for a module:

```bash
➜  GoExec git:(main) ✗ ./GoExec  smb -M list-shares --options                                        
Module: list-shares
Description: List-shares is a SMB module aiming at getting first information about SMB shares such as: open shares, rights, and more.
Tags: SMB
Arguments:
 - Host:  (string)
 - Username: Username to login with (string)
 - Password: Password to login with (string)
```

> [!NOTE]
> Host, Username and Password are shared structures, and GoExec flatten it so that we can use it easily.

Running a module:

```bash
➜  GoExec git:(main) ✗ ./GoExec smb -M list-shares --Host localhost --Username admin --Password admin 
2026/01/15 12:55:30 INFO Starting module module=list-shares
...
2026/01/15 12:55:30 INFO Module executed successfully module=list-shares duration=1.070833ms result=Success
```

