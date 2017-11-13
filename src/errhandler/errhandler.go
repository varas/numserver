package errhandler

import (
	"log"
)

// ErrHandler handles an error
// In case errors should be handled in different ways, a proper error type could be more suitable
type ErrHandler func(error)

// Logger creates a new error handler that uses stdlib log with given prefix
func Logger(preffix string) ErrHandler {
	return func(err error) {
		if err == nil {
			return
		}

		prev := log.Prefix()
		log.SetPrefix(preffix)

		log.Println(err)

		log.SetPrefix(prev)
	}
}

// Noop discards all errors
var Noop ErrHandler = func(err error) {}
