# Chroma CLI in Golang

This is an experiment in building CLI experiences for ChromaDB developers.

## Commands to support

- ✅ Add Server (host, port) - `chroma server add <server-alias> -h <host> -p <port> -o`
- ✅ List Servers - `chroma server ls`
- ✅ Remove Server - `chroma server rm <server-id>`
- ✅ Switch Server, Tenant or Database - `chroma use -s -t -d`
- 🚫 List Collections - `chroma ls` or `chroma c/collection ls`
- 🚫 Create Collection - `chroma create <collection-name>` or `chroma c/collection create <collection-name> -e -d`
- 🚫 Delete Collection - `chroma remove <collection-name>` or `chroma c/collection rm <collection-name>`
- 🚫 Copy Collection - `chroma copy <collection-name> <new-collection-name>` or `chroma c/collection cp <collection-name> <new-collection-name>`
  or `chroma c cp <collection-name> <new-collection-name>` (remote to local or local to remote will be supported in the
  near future)
- 🚫 List Documents - `chroma docs ls <collection-name>` (using bubblegum interactive tables)

Interactive mode - a mode where you can interact with the server using GUI based interface.