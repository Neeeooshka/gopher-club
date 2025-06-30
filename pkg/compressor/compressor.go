package compressor

import (
	"io"
	"log"
	"net/http"
	"strings"
)

type (
	Compressor interface {
		NewWriter(http.ResponseWriter) io.WriteCloser
		NewReader(io.ReadCloser) (io.ReadCloser, error)
		GetEncoding() string
	}

	compressorWriter struct {
		w        http.ResponseWriter
		cw       io.WriteCloser
		encoding string
	}
	compressorReader struct {
		r  io.ReadCloser
		cr io.ReadCloser
	}
)

func newCompressorWriter(w http.ResponseWriter, c Compressor) *compressorWriter {
	return &compressorWriter{
		w:        w,
		cw:       c.NewWriter(w),
		encoding: c.GetEncoding(),
	}
}

func (c *compressorWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressorWriter) Write(p []byte) (int, error) {
	return c.cw.Write(p)
}

func (c *compressorWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", c.encoding)
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressorWriter) Close() error {
	return c.cw.Close()
}

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

type compressor struct {
	c Compressor
}

func NewCompressor(c Compressor) *compressor {
	return &compressor{c: c}
}

func (c *compressor) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ow := w

		contentEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, c.c.GetEncoding()) {
			cr, err := newCompressorReader(r.Body, c.c)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = cr
			defer func() {
				if err := cr.Close(); err != nil {
					log.Printf("failed to close compressor reader: %v", err)
				}
			}()
		}

		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, c.c.GetEncoding()) {
			w.Header().Set("Content-Type", "application/x-"+c.c.GetEncoding())
			w.Header().Set("Content-Encoding", c.c.GetEncoding())
			cw := newCompressorWriter(w, c.c)
			ow = cw
			defer func() {
				if err := cw.Close(); err != nil {
					log.Printf("failed to close compressor writer: %v", err)
				}
			}()
		}

		next.ServeHTTP(ow, r)
	}

	return http.HandlerFunc(fn)
}
