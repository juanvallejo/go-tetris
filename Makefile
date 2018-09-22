all:
	mkdir -p bin && go build -o bin/tetris cmd/main.go
