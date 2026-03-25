package tunnel

import (
	"db-sync/config"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

type Tunnel struct {
	SSHHost    string // the remote server IP/hostname
	SSHPort    int    // usually 22
	SSHUser    string // SSH user
	LocalPort  int    // random free port on your machine
	RemotePort int    // MySQL port on remote, usually 3306
	KeyPath    string // Path to the SSH keys
	// internal fields
	listener net.Listener
	client   *ssh.Client
}

var tunnel *Tunnel

// tunnel := &Tunnel{
// 	SSHHost: config.AppConfig.Source.SSH.Host,
// 	SSHPort: config.AppConfig.Source.SSH.Port,
// 	SSHUser: config.AppConfig.Source.SSH.User,
// 	KeyPath: config.AppConfig.Source.SSH.KeyPath,
// 	RemotePort: config.AppConfig.Source.Port,
// 	LocalPort: 0,
// }

// Open establishes the tunnel, returns the local port to connect to
func (t *Tunnel) Open() error {

	// Load private key
	key, err := os.ReadFile(t.KeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}

	// parse private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// 2. Connect to the remote SSH server (destination)
	sshAddr := fmt.Sprintf("%s:%d", t.SSHHost, t.SSHPort)
	// hostKeyCallback, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")) // for production
	clientConfig := &ssh.ClientConfig{
		User: t.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

	}

	// connect to SSH
	t.client, err = ssh.Dial("tcp", sshAddr, clientConfig)
	if err != nil {
		return fmt.Errorf("ssh dial to %s: %w", sshAddr, err)
	}

	// 3. Open a local TCP listener on the source machine
	localAddr := fmt.Sprintf("%s:%d", config.AppConfig.Source.Host, t.LocalPort)
	t.listener, err = net.Listen("tcp", localAddr)
	if err != nil {
		t.client.Close()
		return fmt.Errorf("failed to listen locally on %s: %w", localAddr, err)
	}

	// If using port 0, capture assigned port
	if t.LocalPort == 0 {
		t.LocalPort = t.listener.Addr().(*net.TCPAddr).Port
	}

	// 5. Accept loop — each connection gets its own goroutine
	go t.acceptLoop()

	return nil
}

// acceptLoop forwards every new local connection through the SSH tunnel
func (t *Tunnel) acceptLoop() {
	for {
		localConn, err := t.listener.Accept()
		if err != nil {
			// listener was closed intentionally — stop quietly
			return
		}
		go t.forward(localConn)
	}
}

// forward pipes one local connection ↔ one remote connection

// Now you have two open connections:
// 1. localConn — our program talking to 127.0.0.1:XXXX on our machine
// 2. remoteConn — the SSH client talking to MySQL on the remote server

func (t *Tunnel) forward(localConn net.Conn) {
	defer localConn.Close()

	remoteAddr := fmt.Sprintf("127.0.0.1:%d", t.RemotePort)
	remoteConn, err := t.client.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("tunnel: remote dial %s: %v", remoteAddr, err)
		return
	}
	defer remoteConn.Close()

	// Bidirectional copy — stop both directions when either side closes
	done := make(chan struct{}, 2)
	go func() {
		io.Copy(remoteConn, localConn) // your program → MySQL
		done <- struct{}{}
	}()
	go func() {
		io.Copy(localConn, remoteConn) // MySQL → your program
		done <- struct{}{}
	}()
	<-done // wait for the first half to finish, then defer closes both
}

// Close tears down the tunnel cleanly
func (t *Tunnel) Close() error {
	if t.listener != nil {
		t.listener.Close()
	}
	if t.client != nil {
		t.client.Close()
	}
	return nil
}
