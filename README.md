# go-mainrun

[![Go Reference](https://pkg.go.dev/badge/github.com/payfazz/go-mainrun.svg)](https://pkg.go.dev/github.com/payfazz/go-mainrun)

Utility for main package

## How to use

```go
func main() { mainrun.Run(run) }

func run(ctx context.Context) error {

  // ...

  return nil
}
```

The `ctx` passed to run will be cancelled if the program caught os signal
(graceful shutdown).

The returned `error` will be printed to `stderr` and the program will be exit
with exit code 1
