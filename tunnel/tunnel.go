package tunnel

import (
	"net"

	"golang.org/x/crypto/ssh"
)

type Tunnel struct {
	SSHHost    string // the remote server IP/hostname
	SSHPort    int    // usually 22
	SSHUser    string // SSH user
	LocalPort  int    // random free port on your machine
	RemotePort int    // MySQL port on remote, usually 3306
	// internal fields
	listener net.Listener
	client   *ssh.Client
}

// Open establishes the tunnel, returns the local port to connect to
func (t *Tunnel) Open() error {}

// Close tears down the tunnel cleanly
func (t *Tunnel) Close() error {}
