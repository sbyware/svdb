# svq

service query engine

## Installation

### Windows

#### x64

```bash
curl -L https://github.com/sbyware/svq/releases/download/latest/svq-windows-x64.exe -o svq.exe
```

#### ARM

```bash
curl -L https://github.com/sbyware/svq/releases/download/latest/svq-windows-arm64.exe -o svq.exe
```

### Linux

#### x64

```bash
curl -L https://github.com/sbyware/svq/releases/download/latest/svq-linux-x64 -o ~/.local/bin/svq
```

#### ARM

```bash
curl -L https://github.com/sbyware/svq/releases/download/latest/svq-linux-arm64 -o ~/.local/bin/svq
```

### macOS

#### Intel (x64)

```bash
curl -L https://github.com/sbyware/svq/releases/download/latest/svq-macos-x64 -o ~/.local/bin/svq
```

#### Apple Silicon (ARM)

```bash
curl -L https://github.com/sbyware/svq/releases/download/latest/svq-macos-arm64 -o ~/.local/bin/svq
```

## Usage

```bash
svq [-p <port,numbers,comma-separated>] [-m <regex-to-match>]
```


Example:

```bash
svq -p 80,443 -m "http"
```

## Configuration

The `.svq` file (copied to ~/.svq) contains the list of services and their respective ports. The file is a JSON array with the following structure:

```json
[
  {
    "name": "http",
    "port": 80
  },
  {
    "name": "https",
    "port": 443
  }
]
```

## the `-select` flag

The `-select` flag allows you to filter the output of `svq` by a specific field. The flag takes a string as an argument and will filter the output based on the value of the field.

#### Get all services with port 80 or 443
```bash
svq -p 80,443 -select port
```

#### Get all services with port 80 or 443
```bash
svq -p 80,443 -select name
```

## the `-j` flag, with `jq`

You can use `jq` to filter the output of `svq` when using the `-j` flag.


#### Get all services with port 80 or 443
```bash
svq -j -p 80,443 | jq '.[] | select(.name == "http")'
```

#### Get all ports for services that match "http"
```bash
svq -j -m "http" | jq -r '.[].port'
```

#### Get just the first description for services on port 443
```bash
svq -j -p 443 | jq -r '.[].description' | head -n 1
```

## License

MIT
