package mainrun

import (
	"fmt"
	"os"
	"sync"

	"github.com/win-t/go-errors"
	"github.com/win-t/go-errors/trace"
)

var onError struct {
	sync.Mutex
	fn func(error) int
}

// When function that passed to [Func] is returned error or panic,
// run f, the returned int will be used to os.Exit function.
func OnError(f func(error) int) {
	onError.Lock()
	defer onError.Unlock()

	onError.fn = f
}

type ExitCodeError struct {
	code int
}

func NewExitCodeError(code int) error {
	return &ExitCodeError{code}
}

func (e *ExitCodeError) Error() string {
	return fmt.Sprintf("exit code %d", e.code)
}

func defaultOnError(err error) int {
	if realErr := (*ExitCodeError)(nil); errors.As(err, &realErr) {
		return realErr.code
	}

	fmt.Fprintln(os.Stderr, errors.FormatWithFilter(err,
		func(l trace.Location) bool { return !l.InPkg("github.com/win-t/go-mainrun") },
	))

	return 1
}
