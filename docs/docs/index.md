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
chroma use <server-alias>
```

With defaults:

```bash
chroma use <server-alias> -r
```

### List Collections

List collection will use the currently active server, tenant and database.

!!! note "Server Alias"

    Specify -s/--alias flag to use a different server.

```bash
chroma list
```

or shorthand:

```bash
chroma ls
```

or

```bash
chroma c ls # c is an alias for `collection`
```

### Create Collection

```bash
chroma create <collection-name> \
  -s <alias> \
  -p/--space <distance_functiom> \
  -m/--m <hnsw:M> \
  -u/--construction-ef <hnsw:efConstruction> \
  -f/--search-ef <hnsw:search_ef> \
  -b/--batch-size <hnsw:batch_size> \
  -k/--sync-threshold <hnsw:sync_threshold> \
  -n/--threads <hnsw:threads> \
  -r/--resize-factor <hnsw:resize_factor> \
  --ensure <create_if_not_exist>
```

### Clone Collection

```bash
chroma clone <collection-name> <target-collection>
```

```bash
chroma clone <collection-name> <target-collection>\
  -s <alias> \
  -p/--space <distance_functiom> \
  -m/--m <hnsw:M> \
  -u/--construction-ef <hnsw:efConstruction> \
  -f/--search-ef <hnsw:search_ef> \
  -b/--batch-size <hnsw:batch_size> \
  -k/--sync-threshold <hnsw:sync_threshold> \
  -n/--threads <hnsw:threads> \
  -r/--resize-factor <hnsw:resize_factor> \
  --embedding-function/-e <embedding-function>
```

All flags are optional and applied to the target collection.

### Delete Collection

```bash
chroma delete/rm <collection-name>
```

### Create Tenant

!!! note "Server Alias"

    Use `-a` or `--alias` flag to specify a server alias for the tenant.

```bash
chroma tenant create <tenant-name>
```

or shorthand:

```bash
chroma c t <tenant-name>
```

### Create Database

!!! note "Server Alias"

    Use `-a` or `--alias` flag to specify a server alias for the database.

By default, if no tenant is specified the database is created in the default tenant (`default_tenant`). To specify a
tenant use `-t` or `--tenant` flag.

```bash
chroma database create <database-name>
```

or shorthand:

```bash
chroma c d <database-name>
```

Create database in a specific tenant:

```bash
chroma c d <database-name> -t <tenant-name>
```

### Version

App version `chroma --version`

Chroma server version `chroma version -s/--alias <server-alias>`
