// The gzipmiddleware package is designed middleware to perform request body decompression.
package gzipmiddleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const (
	contentEncoding = "Content-Encoding"
)

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}
// Read - replacement of the basic similar method, is reading gzip.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close - replacement of the basic similar method, is reading gzip, also close gzip - reader.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// Decompress - middleware who can read gzip body.
func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if isContentGzip(req) {
			cr, err := newCompressReader(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			req.Body = cr
			defer cr.Close()
		}
		next.ServeHTTP(res, req)
	})
}

func isContentGzip(r *http.Request) bool {
	for _, s := range strings.Split(r.Header.Get(contentEncoding), ",") {
		if s == "gzip" {
			return true
		}
	}
	return false
}
