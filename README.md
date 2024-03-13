# Chroma CLI in Golang

This is an experiment in building CLI experiences for ChromaDB developers.

## Getting started

Installation:

```bash
go install github.com/amikos-tech/chroma-cli
```


Add a server:

```bash
chroma server add my-chroma -H localhost -p 8000
```

Create a collection:

```bash
chroma create my-collection
```

List collections:

```bash
chroma ls
```

Delete a collection:

```bash
chroma rm my-collection
```

Clone a collection:

```bash
chroma cp my-collection my-collection-clone
```

> Note: The source collection must have some records in it.

## Commands to support

- ✅ Add Server (host, port) - `chroma server add <server-alias> -h <host> -p <port> -o`
- ✅ List Servers - `chroma server ls`
- ✅ Remove Server - `chroma server rm <server-id>`
- ✅ Switch Server, Tenant or Database - `chroma use -s -t -d`
- ✅ List Collections - `chroma ls` or `chroma c/collection ls`
- ✅ Create Collection - `chroma create <collection-name>` or `chroma c/collection create <collection-name> -e -d`
- ✅ Delete Collection - `chroma remove <collection-name>` or `chroma c/collection rm <collection-name>`
- 🚫 Copy Collection - `chroma copy <collection-name> <new-collection-name>` or `chroma c/collection cp <collection-name> <new-collection-name>`
  or `chroma c cp <collection-name> <new-collection-name>` (remote to local or local to remote will be supported in the
  near future)
- 🚫 List Documents - `chroma docs ls <collection-name>` (using bubblegum interactive tables)
- ✅ App version (via -ldflags) - `chroma --version`
- 🚫 App help
- ✅ Chroma version
- 🚫 Chroma docs
- 🚫 Chroma help

Interactive mode - a mode where you can interact with the server using GUI based interface.


Example config file:

```yaml
active_db: default_database
active_server: test1
active_tenant: default_tenant
servers:
    local:
        host: localhost
        port: "8000"
    myserver:
        database: mydb
        host: 10.10.10.1
        port: 9011
        secure: false
        tenant: my_tenant
    test1:
        database: default_database
        host: localhost
        port: 8000
        secure: false
        tenant: default_tenant
```

### Usage

```bash
make build # or go build/ go install
./chroma server add test1 -h localhost -p 8000 -o
```

