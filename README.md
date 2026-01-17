# Textcat 
Textcat 2.x.x is no longer supported, [article about that here](https://zion8992.github.io/updates/textcat-v3/) <br>
This repository it only textcat *server*, client is available [here](https://github.com/zion8992/textcat-telesto-client).

## Running
To run textcat, download the binaries from the *releases* tab.

Building from source:

1. install go `1.25`
2. build for your system with: `go build -o textcat cmd/main.go`

## Configuring
`-serverName`: set server name (default: "textcat server")
`-serverDesc`: set server description (default: "textcat server")
`-maxSessions`: set server maximum sessions (default: 10)
`-maxMessages`: set maximum messages to cache (default: 5000)

example usage
```
./textcat -serverName "my name" -serverDesc "cool description" -maxSessions 10 -maxMessages 10000
```
