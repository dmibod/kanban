package mongo

// Options declares repository factory options
type Options struct {
	service *DatabaseService
	db      string
}

// Option is a closure which should initialize specific Options properties
type Option func(*Options)

// WithDatabase initializes db option
func WithDatabase(db string) Option {
	return func(o *Options) {
		o.db = db
	}
}
