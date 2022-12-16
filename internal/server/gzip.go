package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/viking311/monitoring/internal/logger"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(nextHandler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			nextHandler.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			logger.Error(err)
			_, ioErr := io.WriteString(w, err.Error())
			if ioErr != nil {
				logger.Error(ioErr)
			}
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		nextHandler.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	}

	return http.HandlerFunc(fn)
}

func UnGzip(nextHandler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			r.Body = gz
			nextHandler.ServeHTTP(w, r)
		} else {
			nextHandler.ServeHTTP(w, r)
		}
	}

	return http.HandlerFunc(fn)
}
