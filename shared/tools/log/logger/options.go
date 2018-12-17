package logger

type Options struct {
	Prefix string
	Debug bool
}

type Option func(*Options)

func WithPrefix(p string) Option {
	return func(o *Options) {
		o.Prefix = p
	}
}

func WithDebug(d bool) Option {
	return func(o *Options) {
		o.Debug = d
	}
}

