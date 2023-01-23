// Package mainrun.
package mainrun

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/win-t/go-errors/v2"
	"github.com/win-t/go-errors/v2/trace"
	"github.com/win-t/go-typedcontext"
)

type osSignal os.Signal

// Return nil if graceful shutdown is not requested yet, otherwise return the signal
func Interrupted(ctx context.Context) os.Signal {
	if s, ok := typedcontext.Get[*osSignal](ctx); ok {
		return *s
	}
	return nil
}

// Run f
//
// The function Run never return.
//
// ctx passed to f will be canceled when graceful shutdown is requested,
// if f returned error or panic, then log it and run os.Exit(1), otherwise run os.Exit(0).
func Run(f func(ctx context.Context) error) {
	exitCode := 1
	defer func() { os.Exit(exitCode) }()

	var s osSignal
	ctx, cancel := context.WithCancel(typedcontext.New(context.Background(), &s))
	defer cancel()

	go func() {
		defer cancel()
		c := make(chan os.Signal, 1)
		signal.Notify(c, getInterruptSigs()...)
		s = <-c
		signal.Stop(c)
	}()

	err := errors.Catch(func() error { return f(ctx) })
	if err == nil {
		exitCode = 0
		return
	}

	onError.Lock()
	defer onError.Unlock()

	if onError.fn != nil {
		exitCode = onError.fn(err)
		return
	}

	fmt.Fprintln(os.Stderr, errors.FormatWithFilter(err,
		func(l trace.Location) bool { return !l.InPkg("github.com/win-t/go-mainrun") },
	))
}

// Go run the f function in new go routine, and return chan to get the value returned by f
func Go(f func() error) <-chan error {
	ch := make(chan error, 1)
	go func() { ch <- errors.Catch(f) }()
	return ch
}

// Result of [Go2]
type Go2Result[Result any] struct {
	Result Result
	Error  error
}

// similar with [Go] but returning some value instead of just error
func Go2[Result any](f func() (Result, error)) <-chan Go2Result[Result] {
	ch := make(chan Go2Result[Result])
	go func() {
		r, err := errors.Catch2(f)
		ch <- Go2Result[Result]{r, err}
	}()
	return ch
}
