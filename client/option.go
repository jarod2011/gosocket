package client

// Options
// Client Options
type Options struct {
	ServerAddr string // Server Address
}

// Option
// callback to set Options value
type Option func(options *Options)

// WithOptions
// Set new Options overwrite current Options
func WithOptions(opt Options) Option {
	return func(options *Options) {
		options = &opt
	}
}

// WithServerAddr
// Set new address overwrite current Options.ServerAddr
func WithServerAddr(addr string) Option {
	return func(options *Options) {
		options.ServerAddr = addr
	}
}
