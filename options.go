package jsonl

// Option configures a Reader.
type Option func(*readerConfig)

type decoderFunc func([]byte, any) error

type readerConfig struct {
	bufferSize  int
	maxLineSize int
	decoder     decoderFunc
}

// WithBufferSize sets the internal buffered reader size.
func WithBufferSize(n int) Option {
	return func(cfg *readerConfig) {
		if n > 0 {
			cfg.bufferSize = n
		}
	}
}

// WithMaxLineSize sets the maximum allowed raw line size in bytes.
func WithMaxLineSize(n int) Option {
	return func(cfg *readerConfig) {
		if n > 0 {
			cfg.maxLineSize = n
		}
	}
}

// WithDecoder replaces the JSON decoder used by Reader.Value and Reader.DecodeInto.
func WithDecoder(fn func([]byte, any) error) Option {
	return func(cfg *readerConfig) {
		if fn != nil {
			cfg.decoder = fn
		}
	}
}
