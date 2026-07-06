package jsonl

// Option configures a Reader.
type Option func(*readerConfig)

type decoderFunc func([]byte, any) error

type marshalFunc func(any) ([]byte, error)

type readerConfig struct {
	bufferSize  int
	maxLineSize int
	decoder     decoderFunc
}

type writerConfig struct {
	bufferSize int
	marshal    marshalFunc
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

// WriterOption configures a Writer.
type WriterOption func(*writerConfig)

// WithWriterBufferSize sets the internal buffered writer size.
func WithWriterBufferSize(n int) WriterOption {
	return func(cfg *writerConfig) {
		if n > 0 {
			cfg.bufferSize = n
		}
	}
}

// WithMarshal replaces the JSON marshal function used by Writer.Write.
func WithMarshal(fn func(any) ([]byte, error)) WriterOption {
	return func(cfg *writerConfig) {
		if fn != nil {
			cfg.marshal = fn
		}
	}
}
