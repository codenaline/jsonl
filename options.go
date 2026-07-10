package jsonl

// ReaderOption configures a Reader.
type ReaderOption func(*readerConfig)

type unmarshalFunc func([]byte, any) error

type marshalFunc func(any) ([]byte, error)

type readerConfig struct {
	bufferSize  int
	maxLineSize int
	unmarshal   unmarshalFunc
}

type writerConfig struct {
	bufferSize int
	marshal    marshalFunc
}

// WithReaderBufferSize sets the internal buffered reader size.
func WithReaderBufferSize(n int) ReaderOption {
	return func(cfg *readerConfig) {
		if n > 0 {
			cfg.bufferSize = n
		}
	}
}

// WithMaxLineSize sets the maximum allowed raw line size in bytes.
func WithMaxLineSize(n int) ReaderOption {
	return func(cfg *readerConfig) {
		if n > 0 {
			cfg.maxLineSize = n
		}
	}
}

// WithUnmarshal replaces the JSON unmarshal used by Reader.Value and Reader.DecodeInto.
func WithUnmarshal(fn func([]byte, any) error) ReaderOption {
	return func(cfg *readerConfig) {
		if fn != nil {
			cfg.unmarshal = fn
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
