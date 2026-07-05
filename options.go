package jsonl

// Option configures a Reader.
type Option func(*readerConfig)

type readerConfig struct {
	bufferSize int
}

// WithBufferSize sets the internal buffered reader size.
func WithBufferSize(n int) Option {
	return func(cfg *readerConfig) {
		if n > 0 {
			cfg.bufferSize = n
		}
	}
}
