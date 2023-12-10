package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

func StartSSHServer() {
	listener, err := net.Listen("tcp", "localhost:20022")
	if err != nil {
		log.Fatalf("Failed to listen on port 22: %v", err)
	}

	log.Println("Starting SSH server on port 22...")

	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection: %v", err)
			continue
		}

		go handleSSHConnection(nConn)
	}
}

var ErrUnauthorized = errors.New("unauthorized")

func handleSSHConnection(nConn net.Conn) {
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			log.Printf("Attempted login: user=%s, password=%s", c.User(), string(pass))
			return nil, ErrUnauthorized // Reject all logins, but log them
		},
	}

	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Printf("Failed to load private key: %v", err)
		return
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Printf("Failed to parse private key: %v", err)
		return
	}

	config.AddHostKey(private)

	// Accept an incoming SSH connection
	conn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Printf("Failed to handshake: %v", err)
		return
	}

	// Log the connection details
	log.Printf("New SSH connection from %s (%s)", conn.RemoteAddr(), conn.ClientVersion())

	go handleChannels(chans)
	go ssh.DiscardRequests(reqs)
}

func handleChannels(chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("Could not accept channel: %v", err)
			continue
		}

		go handleChannelRequests(requests)
		go handleChannelSession(channel)
	}
}

func handleChannelRequests(requests <-chan *ssh.Request) {
	for req := range requests {
		switch req.Type {
		case "shell":
			if len(req.Payload) == 0 {
				req.Reply(true, nil)
			}
		case "pty-req":
			req.Reply(true, nil)
		}
	}
}

func handleChannelSession(channel ssh.Channel) {
	// Session logging
	sessionStartTime := time.Now()
	shell := "$ "
	channel.Write([]byte(shell))

	var sessionCommands []string

	// Map of commands to responses
	commandResponses := map[string]string{
		"ls":     "file1.txt\nfile2.txt\n",
		"pwd":    "/home/user\n",
		"whoami": "user\n",
		// Add more commands and responses here
	}

	for {
		data := make([]byte, 256)
		n, err := channel.Read(data)
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read from channel: %v", err)
			}
			break
		}

		// Log the command and simulate a response
		command := string(data[:n])
		sessionCommands = append(sessionCommands, command)
		log.Printf("Command run: %s", command)
		response, ok := commandResponses[command]
		if !ok {
			response = "command not found\n"
		}
		channel.Write([]byte(response))
		channel.Write([]byte(shell))
	}

	channel.Close()

	// Log the session details
	sessionEndTime := time.Now()
	log.Printf("Session started at %s and ended at %s", sessionStartTime, sessionEndTime)
	log.Printf("Commands run: %v", sessionCommands)
}
