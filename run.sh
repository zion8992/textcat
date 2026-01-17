echo "$1 - max sessions"
echo "$2 - max messages"
echo ""

go run ./cmd/ -serverName "development server" -serverDesc "this server is in development" -maxSessions $1 -maxMessages $2
