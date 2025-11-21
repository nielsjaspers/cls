> [!IMPORTANT]
> `cls` has been deprecated for a simpler alternative, requiring much less setup
>
> Visit [`clifs`](https://github.com/nielsjaspers/clifs) instead!


# CLS (Command Line Share)

CLS (Command Line Share) is a lightweight command-line tool written in Go for transferring files between a client and a server on a local network. It's simple, fast, and efficient, making it a great solution for local file sharing.

## Features

- **Send Files:** Quickly share files from a client to the server.
- **List Files:** View all files available on the server.
- **Retrieve Files:** Download files from the server to the client.

## Installation

### Clone the Repository

1. Clone the repository:
   ```bash
   git clone https://github.com/nielsjaspers/cls.git
   cd cls
   ```

### Generate TLS/SSL Keys and Self-Sign Certificates

2. Create a directory to store the server certificate and configuration:

   ```bash
   mkdir server-crt
   cd server-crt
   ```

3. Create a file named `server.cnf` in the `server-crt` directory with the following content:

   ```bash
   [req]
   default_md = sha256
   prompt = no
   req_extensions = v3_ext
   distinguished_name = req_distinguished_name

   [req_distinguished_name]
   CN = localhost

   [v3_ext]
   keyUsage = critical,digitalSignature,keyEncipherment
   extendedKeyUsage = critical,serverAuth
   subjectAltName = DNS:localhost
   ```

4. Run the `keygen.sh` script at the project root to generate and sign the certificates used by the client and server:

   ```
   ./keygen.sh
   ```

### Build the Binaries

5) Build the client and server binaries separately:

   ```bash
   go build -o client ./cmd/cls/client
   go build -o server ./cmd/cls/server
   ```

6) The resulting `client` and `server` binaries are your executables for client and server operations, respectively.

## Usage

### Server

To start the server and set up the directory for incoming files:

```bash
./server path <path/to/folder>
```

This will configure the server to save all received files in the specified folder.

#### Start the Server Using Docker

To build and run a Docker container containing the server, you can use the provided shell script:

```bash
./dockersetup.sh
```

This script will handle building the Docker image and starting the container. Ensure you have Docker installed and running before executing the script.

### Accessing the Server on a Local Network (WIP)

If you want to make the server accessible on a local network, you can use a free Cloudflare tunnel to access the server with any device running the client. Documentation for this feature is still a work in progress.

### Client

#### Share a File

To share a file with the server:

```bash
./client share <path/to/file>
```

#### List Files

To get a list of all files available on the server:

```bash
./client list
```

#### Retrieve a File

To download a specific file from the server and download it to a local path:

```bash
./client get <remote/file> <local/path>
```

## Example

### Server:

Start the server and set the folder for incoming files:

```bash
./server path /home/user/shared_files
```

### Client:

1. Share a file:

   ```bash
   ./client share /home/user/documents/report.pdf
   ```

2. List files on the server:

   ```bash
   ./client list
   ```

3. Retrieve a file from the server and save it to `home/user/`:

   ```bash
   ./client get report.pdf home/user/
   ```

## Requirements

- Go 1.23.3 or higher

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

## License

This project is licensed under the Apache-2.0 License. See the `LICENSE` file for details.
