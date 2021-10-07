package log

import "net/http"

type Writer struct {
	http.ResponseWriter
	StatusCode int
}

func NewLogWriter(w http.ResponseWriter) *Writer {
	return &Writer{w, http.StatusOK}
}

func (w *Writer) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}
