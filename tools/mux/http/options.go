package http

// Options can be used to create a customized mux.
type Options struct {
	Port int
}

// Option is a function on the options for a http mux.
type Option func(*Options)

func WithPort(p int) Option {
	return func(o *Options) {
		o.Port = p
	}
}
