# Go Honeypot Server Project

This project includes an HTTP server, an SSH server, and an FTP server, all implemented in Go. It also includes a SQLite database for logging.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (version 1.16 or later recommended)
- SQLite3

### Installing

1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Run `go build` to compile the project.
4. Run `./main` (or `main.exe` on Windows) to start the servers.

## Usage

Once the servers are running, they will start listening for connections:

- The HTTP server starts on port 8080.
- The SSH server starts on specified ports.
- The FTP server starts on specified ports.

All login attempts to the SSH and FTP servers are logged and rejected. The HTTP server logs all requests to a SQLite database.

## Contributing

Please read `CONTRIBUTING.md` for details on our code of conduct, and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the `LICENSE.md` file for details.
