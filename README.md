jabba-ai-chat-app
=================

The ChatApp of jabba-ai

* [jabba-ai](https://github.com/Koubae/jabba-ai)


### QuickStart

* 1) Install [air-verse/air](https://github.com/air-verse/air) globally

```bash
go install github.com/air-verse/air@latest

# Make sure that GOPATH and GOROOT is in your PATH
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin
```

* 2) Initialize (this will perform some installations)

```bash
make init
```

* 3) Start

```bash
# Run chat-identity
make run-chat-identity
```

### Services

| Service             | Port    |
|---------------------|---------|
| `chat-identity`     | `20000` |
| `chat-orchestrator` | `20001` |
| `chat-session`      | `20002` |


chat-identity
-------------

### Generate Pub/Sub Keys for JWT Auth

```bash
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -out public.pem

## Better naming
openssl genrsa -out cert_private.pem 2048
openssl rsa -in cert_private.pem -pubout -out cert_public.pem

```
