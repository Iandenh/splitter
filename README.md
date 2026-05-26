# Splitter

Proxy incoming requests to multiple upstreams. Handy for fanning out webhooks to multiple targets simultaneously.

## Features

* **Simple:** Single binary, no heavy dependencies.
* **Configurable:** Easy-to-use YAML configuration.
* **Efficient:** Written in Go for high performance and low resource usage.

## Installation

### Prerequisites
Ensure you have [Go](https://go.dev/doc/install) installed on your system.

### Install via `go install`
The quickest way to install `splitter` is using the Go toolchain. This will download and compile the binary into your `$GOPATH/bin` directory.

```bash
go install github.com/Iandenh/splitter@latest
```

> **Note:** Ensure your `$(go env GOPATH)/bin` directory is added to your system's `$PATH` so you can run the `splitter` command globally from your terminal.

## Configuration

`splitter` relies on a YAML configuration file to define the upstream targets. 

Use the provided `example-config.yaml` as a starting point or use `splitter init`. Here is an example of what your `splitter.yaml` might look like:

```yaml
# splitter.yaml
port: 8080
upstreams:
  - https://api.target-one.com/webhook
  - https://api.target-two.com/webhook
```

## Usage

Start the proxy server by pointing it to your configuration file:

```bash
# will automatically look for splitter.yaml
splitter start

splitter start --config config.yaml
```

Once running, any HTTP request sent to the splitter (e.g., `http://localhost:8080`) will be duplicated and proxied to all defined upstreams.

---

### Building from Source (Local Development)

If you prefer to clone the repository and build it manually:

```bash
git clone https://github.com/Iandenh/splitter.git
cd splitter
go build -o splitter main.go
./splitter start --config config.yaml
```
