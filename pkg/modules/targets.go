package modules

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

/*
Modules can take as parameter a "target":
1. One IP
2. One Host
3. CIDR (e.g. 192.168.0.0/24)
4. File containing multiple IPs, Hosts, CIDRs

This file aims at handling all of them.

Current choice: modules handle only one target at a time, and the core multiplies module calls.
*/

// Given by the user
type Targets struct {
	// Could be an ip, ~~cidr~~, or file 
	Host string `help:"Host of the target (e.g., IP, hostname)"`
	
	// File with multiple hosts
	HostFile string `help:"Path to a file that contains multiple hosts"`

	// CIDR string
}

// Given by the program to each module
type ModuleTarget struct {
	// Could be an IP or a hostname
	Host string
}

// Returns the resolved IP address of Host:
// - If Host is already an IP, it returns it directly.
// - Otherwise, it resolves the hostname using DNS.
func (mt ModuleTarget) ResolveToIP() (string, error) {
	// Try parsing as IP first
	ip := net.ParseIP(mt.Host)
	if ip != nil {
		return ip.String(), nil
	}

	// Otherwise resolve as hostname
	ips, err := net.LookupIP(mt.Host)
	if err != nil || len(ips) == 0 {
		return "", fmt.Errorf("failed to resolve host '%s': %w", mt.Host, err)
	}

	// Return the first IPv4 if possible, else the first IP
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String(), nil
		}
	}

	return ips[0].String(), nil
}

// ResolveToDomain resolves the Host to a domain name.
// - If Host is already a hostname, returns it.
// - If Host is an IP, performs a reverse DNS lookup.
func (mt ModuleTarget) ResolveToDomain() (string, error) {
	ip := net.ParseIP(mt.Host)
	if ip == nil {
		// Already a hostname
		return mt.Host, nil
	}

	// Reverse DNS lookup for IP
	names, err := net.LookupAddr(mt.Host)
	if err != nil || len(names) == 0 {
		return "", fmt.Errorf("failed to resolve IP '%s' to domain: %w", mt.Host, err)
	}

	// Return first name without trailing dot
	name := names[0]
	if len(name) > 0 && name[len(name)-1] == '.' {
		name = name[:len(name)-1]
	}

	return name, nil
}

// When a user gives a host, host file, etc., we have to return a list of hosts
func GetTargets(targets Targets) ([]ModuleTarget, error) {
	var outTargets []ModuleTarget	

	// If hostfile is defined, ignore the rest and import the file 
	if targets.HostFile != "" {
		file, err := os.Open(targets.HostFile)
		if err != nil {
			return nil, err
		}
		var fileCloseErr error
		defer func() {
			if cerr := file.Close(); cerr != nil {
				fileCloseErr = cerr
			}
		}()
	
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			
			if line != "" {
				outTargets = append(outTargets, ModuleTarget{Host: line})
			}
		}
	
		return outTargets, fileCloseErr
	}

	// TODO: Modify to expand CIDR?
	
	// Fallback, no file given, need to try the host itself
	if targets.Host == "" {
		return nil, fmt.Errorf("Given host is empty")
	}

	outTargets = []ModuleTarget{
		{Host: targets.Host},
	}

	return outTargets, nil
}
