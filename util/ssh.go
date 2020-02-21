package util

import (
	"bufio"
	"io"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	*ssh.Client
	password string
	output   chan byte
}

func NewSshClient(addr, user, password string) (*SshClient, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}
	return &SshClient{conn, password, make(chan byte, 1024)}, nil
}

func (conn *SshClient) SendCommands(cmds ...string) error {
	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return err
	}
	in, err := session.StdinPipe()
	if err != nil {
		return err
	}
	out, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	go func(in io.WriteCloser, out io.Reader, output chan byte) {
		var (
			err error
			b   byte
			bs  = make([]byte, 0)
			r   = bufio.NewReader(out)
		)
		for {
			if b, err = r.ReadByte(); err != nil {
				break
			}
			output <- b
			if b == byte('\n') {
				bs = make([]byte, 0)
				continue
			}
			bs = append(bs, b)
			if err := conn.inputSudoPass(string(bs), in); err != nil {
				break
			}
		}
	}(in, out, conn.output)

	cmd := strings.Join(cmds, "; ")
	if _, err = session.Output(cmd); err != nil {
		return err
	}
	return nil
}

func (conn *SshClient) GetOutput() <-chan byte {
	return conn.output
}

func (conn *SshClient) CloseClient() {
	close(conn.output)
	_ = conn.Close()
}

func (conn *SshClient) inputSudoPass(line string, in io.WriteCloser) error {
	if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
		if _, err := in.Write([]byte(conn.password + "\n")); err != nil {
			return err
		}
	}
	return nil
}
