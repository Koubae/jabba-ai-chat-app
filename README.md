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

* 2) Install dependencies

```bash
go mod tidy
```

* 3) Start

```bash
# Run chat-identity
make run-chat-identity
```
