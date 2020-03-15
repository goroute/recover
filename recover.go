package recover

import (
	"fmt"
	"runtime"

	"github.com/goroute/route"
)

// Options defines the config for Recover middleware.
type Options struct {
	// Skipper defines a function to skip middleware.
	Skipper route.Skipper

	// Size of the stack to be printed.
	// Optional. Default value 4KB.
	StackSize int `yaml:"stack_size"`

	// DisableStackAll disables formatting stack traces of all other goroutines
	// into buffer after the trace for the current goroutine.
	// Optional. Default value false.
	DisableStackAll bool `yaml:"disable_stack_all"`

	// OnError allows to set callback func to handle panic.
	OnError func(err error, stack []byte) `yaml:"disable_print_stack"`
}

// Option defines option func.
type Option func(*Options)

// GetDefaultOptions returns default options.
func GetDefaultOptions() Options {
	return Options{
		Skipper:         route.DefaultSkipper,
		StackSize:       4 << 10, // 4 KB
		DisableStackAll: false,
		OnError:         func(err error, stack []byte) {},
	}
}

// Skipper sets skipper option.
func Skipper(skipper route.Skipper) Option {
	return func(o *Options) {
		o.Skipper = skipper
	}
}

// StackSize sets stack size option.
func StackSize(stackSize int) Option {
	return func(o *Options) {
		o.StackSize = stackSize
	}
}

// DisableStackAll sets disable stack option.
func DisableStackAll(disableStackAll bool) Option {
	return func(o *Options) {
		o.DisableStackAll = disableStackAll
	}
}

// OnError sets error callback option.
func OnError(fn func(err error, stack []byte)) Option {
	return func(o *Options) {
		o.OnError = fn
	}
}

// New returns a Recover middleware with config.
func New(options ...Option) route.MiddlewareFunc {
	// Apply options.
	opts := GetDefaultOptions()
	for _, opt := range options {
		opt(&opts)
	}

	return func(c route.Context, next route.HandlerFunc) (err error) {
		if opts.Skipper(c) {
			return next(c)
		}

		defer func() {
			if r := recover(); r != nil {
				var ok bool
				err, ok = r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				stack := make([]byte, opts.StackSize)
				length := runtime.Stack(stack, !opts.DisableStackAll)
				opts.OnError(err, stack[:length])
				c.Error(err)
			}
		}()
		err = next(c)
		return
	}
}
