# Chroma CLI

## Installation

```bash
```

## Usage

### Add Server

Arguments:

- `server-alias` - Alias for the server

Flags:

- `-h` or `--host` - Hostname of the server
- `-p` or `--port` - Port of the server
- `-o` or `--overwrite` - Overwrite existing server alias
- `-s` or `--secure` - Use secure connection (TLS)
- `-t` or `--tenant` - Tenant Name to use by default
- `-d` or `--database` - Database Name to use by default

```bash
chroma server add <server-alias> -h <host> -p <port> -o -s -t <tenant> -d <database>
```

### Switch Server

Arguments:

- `server-alias` - Alias for the server

Flags:

- `-t` or `--tenant` - Changes the default tenant
- `-d` or `--database` - Changes the default database
- `-r` or `--defaults` - Uses the default tenant and database (this is mutually exclusive with `-t` and `-d` flags)

```bash
chroma switch <server-alias>
```

With shorthands:

```bash
chroma sw <server-alias> -t <tenant> -d <database>
```

With defaults:

```bash
chroma sw <server-alias> -r
```

### List Collections

List collection will use the currently active server, tenant and database.

```bash
chroma ls
```

or 

```bash
chroma c ls # c is an alias for `collection`
```
