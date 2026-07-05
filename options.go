package jsonl

// Option configures a Reader.
type Option func(*readerConfig)

type decoderFunc func([]byte, any) error

type readerConfig struct {
	bufferSize int
	decoder    decoderFunc
}

// WithBufferSize sets the internal buffered reader size.
func WithBufferSize(n int) Option {
	return func(cfg *readerConfig) {
		if n > 0 {
			cfg.bufferSize = n
		}
	}
}

// WithDecoder replaces the JSON decoder used by Reader.Value.
func WithDecoder(fn func([]byte, any) error) Option {
	return func(cfg *readerConfig) {
		if fn != nil {
			cfg.decoder = fn
		}
	}
}
