# svdb

the service database (tcp & ucp).

## Installation

Copy `svdb` binary to a directory in your PATH. Place the `svdb.json` file in this same directory or the tool will not be able to find it.

## Usage

```bash
svdb [-p <port,numbers,comma-separated>] [-m <regex-to-match>]
```


Example:

```bash
svdb -p 80,443 -m "http"
```

## Configuration

The `.svdb` file (copied to ~/.svdb) contains the list of services and their respective ports. The file is a JSON array with the following structure:

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

The `-select` flag allows you to filter the output of `svdb` by a specific field. The flag takes a string as an argument and will filter the output based on the value of the field.

#### Get all services with port 80 or 443
```bash
svdb -p 80,443 -select port
```

#### Get all services with port 80 or 443
```bash
svdb -p 80,443 -select name
```

## the `-j` flag, with `jq`

You can use `jq` to filter the output of `svdb` when using the `-j` flag.


#### Get all services with port 80 or 443
```bash
svdb -j -p 80,443 | jq '.[] | select(.name == "http")'
```

#### Get all ports for services that match "http"
```bash
svdb -j -m "http" | jq -r '.[].port'
```

#### Get just the first description for services on port 443
```bash
svdb -j -p 443 | jq -r '.[].description' | head -n 1
```

## License

MIT
