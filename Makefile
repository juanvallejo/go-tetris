all:
	mkdir -p bin && go build -o bin/tictactoe cmd/main.go
