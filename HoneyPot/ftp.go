package main

import (
	"fmt"
	"log"
	"net"
)

func StartFTPServer() {
	listener, err := net.Listen("tcp", "0.0.0.0:21")
	if err != nil {
		log.Fatalf("Failed to listen on port 21: %v", err)
	}

	log.Println("Starting FTP server on port 21...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// TODO: Implement FTP protocol handling here
	// Read and write data to the connection using conn.Read() and conn.Write() respectively
	// Handle FTP commands and responses according to the FTP protocol specification

	// Example: Send FTP welcome message
	conn.Write([]byte("220 Welcome to the FTP server\r\n"))

	// Example: Read FTP command from the client
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection:", err.Error())
		return
	}
	command := string(buffer[:n])

	// Example: Handle FTP command
	switch command {
	case "USER":
		// Handle USER command
		conn.Write([]byte("331 User name okay, need password\r\n"))
	case "PASS":
		// Handle PASS command
		conn.Write([]byte("230 User logged in, proceed\r\n"))
	default:
		// Handle unsupported command
		conn.Write([]byte("500 Syntax error, command unrecognized\r\n"))
	}
}
