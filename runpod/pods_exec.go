package runpod

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
)

type ErrExecOnSSH struct {
	StdErr string
}

func (e *ErrExecOnSSH) Error() string {
	return e.StdErr
}

type ErrSshDisabled struct {
	PodId string
}

func (e *ErrSshDisabled) Error() string {
	return fmt.Sprintf("ssh is not enabled on pod with id '%s'; create your pods with tcp port 22 open to enable remote command execution", e.PodId)
}

func (p *Pod) Exec(cmd string, privateKeyPath string) (stdOut string, err error) {

	var sshIp string
	var sshPort string

	for _, port := range p.Runtime.Ports {
		if port.Type == "tcp" && port.PrivatePort == 22 {
			sshPort = strconv.Itoa(port.PublicPort)
			sshIp = port.IP
		}
	}

	if sshIp == "" {
		return "", &ErrSshDisabled{PodId: p.Id}
	}

	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("unable to read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("unable to parse private key: %w", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //❗️❗️❗️ TODO Use secure host key verification in production
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", sshIp, sshPort), sshConfig)
	if err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	out, stdErr := session.CombinedOutput(cmd)
	if stdErr != nil {
		return string(out), &ErrExecOnSSH{StdErr: stdErr.Error()}
	}

	return string(out), nil
}
