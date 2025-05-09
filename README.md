# 🚀 Haxor Client

Haxor Client is a client application for the Haxorport service, allowing you to create HTTP and TCP tunnels to expose local services to the internet.

## ✨ Features

- 🌐 **HTTP/HTTPS Tunnels**: Expose local web services with custom subdomains, supporting both HTTP and HTTPS protocols
- 🔌 **TCP Tunnels**: Expose local TCP services with remote ports
- 🔒 **Authentication**: Protect tunnels with basic or header authentication
- ⚙️ **Configuration**: Easily manage configuration through CLI
- 🔄 **Automatic Reconnection**: Connections will automatically reconnect if disconnected

## 🏗️ Architecture

Haxor Client is built with a hexagonal architecture (ports and adapters) that separates business domain from technical infrastructure. This architecture enables:

1. **Separation of Concerns**: Business domain is separated from technical details
2. **Testability**: Components can be tested separately
3. **Flexibility**: Infrastructure implementations can be replaced without changing the business domain

Project structure:

```
haxor-client/
├── cmd/                    # Command-line interface
├── internal/               # Internal code
│   ├── domain/             # Domain layer
│   │   ├── model/          # Domain models
│   │   └── port/           # Ports (interfaces)
│   ├── application/        # Application layer
│   │   └── service/        # Services
│   ├── infrastructure/     # Infrastructure layer
│   │   ├── config/         # Configuration implementation
│   │   ├── transport/      # Communication implementation
│   │   └── logger/         # Logger implementation
│   └── di/                 # Dependency injection
├── scripts/                # Build and run scripts
└── main.go                 # Entry point
```

## 💻 Installation

### 🚀 Easy Installation (All OS)

Use the automated installer script that supports Linux, macOS, and Windows (via WSL):

```bash
# Download and run the installer
curl -sSL https://raw.githubusercontent.com/alwanandri2712/haxorport-go-client/main/install.sh | bash
```

The installer script will:
- 🔎 Automatically detect your OS
- 📚 Install required dependencies
- 💿 Compile and install haxorport
- 📝 Create default configuration

After installation, you can immediately use the `haxorport` command.

### 🔧 Manual Installation

#### 📂 From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/alwanandri2712/haxorport-go-client.git
   cd haxorport-go-client
   ```

2. Build the application:
   ```bash
   # Make sure Go is installed
   go build -o bin/haxor main.go
   ```

3. (Optional) Move the binary to a directory in your PATH:
   ```bash
   # Linux/macOS
   sudo cp bin/haxor /usr/local/bin/
   
   # Windows (PowerShell Admin)
   Copy-Item .\bin\haxor.exe -Destination "$env:ProgramFiles\haxorport\"
   ```

#### 💾 From Binary

1. Download the latest binary from [releases](https://github.com/alwanandri2712/haxorport-go-client/releases)
2. Extract and move it to a directory in your PATH

## 💬 Usage

### ⚙️ Configuration

Before using Haxor Client, you need to set up the configuration:

#### 🔑 Getting an Auth Token

To obtain an auth-token, you must first register at:

**[https://haxorport.online/Register](https://haxorport.online/Register)**

After registering and logging in, you can find your auth-token in your account dashboard.

#### 📝 Setting Up Configuration

```
haxor config set server_address control.haxorport.online
haxor config set control_port 443
haxor config set auth_token your-auth-token
haxor config set tls_enabled true
```

Or use the easier method with the command:

```
./build.sh config
```

To view the current configuration:

```
haxor config show
```

### 🌐 HTTP Tunnel

Create an HTTP tunnel for a local web service:

```
haxor http --port 8080 --subdomain myapp
```

With basic authentication:

```
haxor http --port 8080 --subdomain myapp --auth basic --username user --password pass
```

With header authentication:

```
haxor http --port 8080 --subdomain myapp --auth header --header "X-API-Key" --value "secret-key"
```

### 🔒 HTTPS Tunnel

Haxorport now supports HTTPS tunnels automatically with a reverse connection architecture. When the client connects to the server, the server detects whether the request comes via HTTP or HTTPS and forwards the request to the client through a WebSocket connection. The client then makes a request to the local service and sends the response back to the server.

Advantages of the reverse connection architecture:

1. **No SSH tunnel required**: You don't need to set up an SSH tunnel to access local services
2. **Automatic URL replacement**: Local URLs in HTML responses are automatically replaced with tunnel URLs
3. **HTTPS support**: Access local services via HTTPS without configuring TLS on the local service
4. **Custom subdomains**: Use easy-to-remember subdomains to access local services

To use an HTTPS tunnel:

1. Ensure the haxorport server is correctly configured to support HTTPS
2. Run the client by specifying the local port and subdomain:
   ```
   haxor http --port 8080 --subdomain myapp
   ```
3. Access your service via HTTPS:
   ```
   https://myapp.haxorport.online
   ```

All links and references in your web pages will be automatically modified to use the tunnel URL, ensuring that navigation on the website works correctly.

### 🔌 TCP Tunnel

Haxorport supports TCP tunnels that allow you to expose local TCP services (such as SSH, databases, or other services) to the internet. TCP tunnels work by forwarding connections from a remote port on the Haxorport server to a local port on your machine.

Create a TCP tunnel for a local TCP service:

```
haxor tcp --port 22 --remote-port 2222
```

If `--remote-port` is not specified, the server will assign a remote port automatically.

Advantages of Haxorport TCP tunnels:

1. **Secure Access**: Access local TCP services from anywhere without opening ports in your firewall
2. **Multi-Protocol Support**: Supports all TCP-based protocols (SSH, MySQL, PostgreSQL, Redis, etc.)
3. **Integrated Authentication**: Uses the same authentication system as HTTP/HTTPS tunnels
4. **Usage Limits**: Control the number of tunnels based on user subscription

Examples of TCP tunnel usage:

- **🔑 SSH Server**:
  ```
  haxor tcp --port 22 --remote-port 2222
  # Access: ssh user@haxorport.online -p 2222
  ```

- **💾 MySQL Database**:
  ```
  haxor tcp --port 3306 --remote-port 3306
  # Access: mysql -h haxorport.online -P 3306 -u user -p
  ```

- **💾 PostgreSQL Database**:
  ```
  haxor tcp --port 5432 --remote-port 5432
  # Access: psql -h haxorport.online -p 5432 -U user -d database
  ```

### 📝 Adding Tunnels to Configuration

You can add tunnels to the configuration for later use:

```
haxor config add-tunnel --name web --type http --port 8080 --subdomain myapp
haxor config add-tunnel --name ssh --type tcp --port 22 --remote-port 2222
```

## 👨‍💻 Development

### 📚 Prerequisites

- Go 1.21 or newer
- Git

### 🔧 Development Setup

1. Clone the repository:
   ```
   git clone https://github.com/alwanandri2712/haxorport-go-client.git
   cd haxorport-go-client
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Run the application in development mode:
   ```
   ./scripts/run.sh
   ```

### 💻 Code Structure

- **Domain Layer**: Contains domain models and ports (interfaces)
- **Application Layer**: Contains services that implement use cases
- **Infrastructure Layer**: Contains concrete implementations of ports
- **CLI Layer**: Contains command-line interface using Cobra
- **DI Layer**: Contains container for dependency injection

## 🔧 Troubleshooting

### 📉 Reducing Debug Output

If you see too many INFO log messages when running the application, you can change the log level to `warn` as follows:

```bash
# Edit configuration file
sudo nano /etc/haxorport/config.yaml  # For Linux
nano ~/.haxorport/config/config.yaml  # For Windows (WSL)
nano ~/Library/Preferences/haxorport/config.yaml  # For macOS
```

Change the line `log_level: info` to `log_level: warn`, then save the file.

Or use the following command to change the log level automatically:

```bash
# For Linux
sudo sed -i 's/log_level:.*/log_level: warn/g' /etc/haxorport/config.yaml

# For macOS
sed -i '' 's/log_level:.*/log_level: warn/g' ~/Library/Preferences/haxorport/config.yaml

# For Windows (WSL)
sed -i 's/log_level:.*/log_level: warn/g' ~/.haxorport/config/config.yaml
```

## 📃 License

MIT License
