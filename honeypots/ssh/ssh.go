package ssh

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/ejcx/honeypotd/honeypots"
	"github.com/ejcx/honeypotd/notification/twilio"
	"golang.org/x/crypto/ssh"
)

type SSHPot struct {
}

func (p *SSHPot) Run(h *honeypots.HoneyPot) error {
	// Public key authentication is done by comparing
	// the public key of a received connection
	// with the entries in the authorized_keys file.
	authorizedKeysBytes, err := ioutil.ReadFile("authorized_keys")
	if err != nil {
		log.Fatalf("Failed to load authorized_keys, err: %v", err)
	}

	authorizedKeysMap := map[string]bool{}
	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			log.Fatal(err)
		}

		authorizedKeysMap[string(pubKey.Marshal())] = true
		authorizedKeysBytes = rest
	}

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
		// Remove to disable password auth.
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			twilio.Notify(fmt.Sprintf("SSH password auth connection from %s\n", c.RemoteAddr()))
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},

		// Remove to disable public key auth.
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			twilio.Notify(fmt.Sprintf("SSH key auth connection from %s\n", c.RemoteAddr()))
			if authorizedKeysMap[string(pubKey.Marshal())] {
				return &ssh.Permissions{
					// Record the public key used for authentication.
					Extensions: map[string]string{
						"pubkey-fp": ssh.FingerprintSHA256(pubKey),
					},
				}, nil
			}
			return nil, fmt.Errorf("unknown public key for %q", c.User())
		},
	}

	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		log.Fatal("failed to listen for connection: ", err)
	}
	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Println("failed to accept incoming connection: ", err)
			continue
		}

		// Before use, a handshake must be performed on the incoming
		// net.Conn.
		conn, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			log.Println("failed to handshake: ", err)
			continue
		}
		if conn.Permissions != nil {
			log.Printf("logged in with key %s", conn.Permissions.Extensions["pubkey-fp"])
		}

		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)

		// Service the incoming Channel channel.
		for newChannel := range chans {
			// Channels have a type, depending on the application level
			// protocol intended. In the case of a shell, the type is
			// "session" and ServerShell may be used to present a simple
			// terminal interface.
			if newChannel.ChannelType() != "session" {
				newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				continue
			}
			channel, _, err := newChannel.Accept()
			if err != nil {
				log.Println("Could not accept channel: %v", err)
				continue
			}
			channel.Close()

			// // Sessions have out-of-band requests such as "shell",
			// // "pty-req" and "env".  Here we handle only the
			// // "shell" request.
			// go func(in <-chan *ssh.Request) {
			// 	for req := range in {
			// 		req.Reply(req.Type == "shell", nil)
			// 	}
			// }(requests)

			// term := terminal.NewTerminal(channel, "> ")

			// go func() {
			// 	defer
			// 	for {
			// 		line, err := term.ReadLine()
			// 		if err != nil {
			// 			break
			// 		}
			// 		fmt.Println(line)
			// 	}
			// }()
		}
	}
	return nil
}
