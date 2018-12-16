package logger

type Options struct {
	Module string
}

type Option func(*Options)

func WithModule(m string) Option {
	return func(o *Options) {
		o.Module = m
	}
}
