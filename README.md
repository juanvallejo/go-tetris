go-tictactoe
=========

Simple tic-tac-toe game written in Go, using an opengl-based graphics library.

## Building

### Requirements

go-tetris requires the [faiface/pixel](https://github.com/faiface/pixel) dependency.

```
go get -u github.com/faiface/pixel
```

Additionally, please ensure you meet the requirements from that repo:
https://github.com/faiface/pixel#requirements

### Compiling

```
make
```

Alternatively, build with:

```
mkdir -p ./bin && go build -o bin/tictactoe ./cmd/main.go
```

### Running

```
./bin/tictactoe
```
