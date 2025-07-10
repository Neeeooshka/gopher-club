package compressor

import (
	"io"
	"log"
	"net/http"
	"strings"
)

type (
	Compressor interface {
		NewReader(io.ReadCloser) (io.ReadCloser, error)
		GetEncoding() string
	}
	compressorReader struct {
		r  io.ReadCloser
		cr io.ReadCloser
	}
)

func newCompressorReader(r io.ReadCloser, c Compressor) (*compressorReader, error) {
	cr, err := c.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressorReader{
		r:  r,
		cr: cr,
	}, nil
}

func (c *compressorReader) Read(p []byte) (n int, err error) {
	return c.cr.Read(p)
}

func (c *compressorReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.cr.Close()
}

type CompressorWrap struct {
	c Compressor
}

func NewCompressor(c Compressor) *CompressorWrap {
	return &CompressorWrap{c: c}
}

func (c *CompressorWrap) GetEncoding() string {
	return c.c.GetEncoding()
}

func (c *CompressorWrap) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		contentEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, c.GetEncoding()) {
			cr, err := newCompressorReader(r.Body, c.c)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer func() {
				if err := cr.Close(); err != nil {
					log.Printf("failed to close compressor reader: %v", err)
				}
			}()
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
