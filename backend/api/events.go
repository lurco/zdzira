package api

import (
	"fmt"
	"net/http"
	"time"
)

func eventsHandler(b *Broadcaster) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch, unsubscribe := b.Subscribe()
		defer unsubscribe()

		fmt.Fprintf(w, "data: connected\n\n")
		flusher.Flush()

		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ch:
				fmt.Fprintf(w, "data: refresh\n\n")
				flusher.Flush()
			case <-ticker.C:
				fmt.Fprintf(w, ": keepalive\n\n")
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	}
}

// notifyWriter captures the response status so notifyMiddleware can decide
// whether to broadcast after a successful mutation.
type notifyWriter struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func (nw *notifyWriter) WriteHeader(code int) {
	if !nw.wrote {
		nw.status = code
		nw.wrote = true
	}
	nw.ResponseWriter.WriteHeader(code)
}

func (nw *notifyWriter) Write(b []byte) (int, error) {
	if !nw.wrote {
		nw.status = http.StatusOK
		nw.wrote = true
	}
	return nw.ResponseWriter.Write(b)
}

func (nw *notifyWriter) Unwrap() http.ResponseWriter { return nw.ResponseWriter }

// NotifyMiddleware fires NotifyAll after any successful (2xx) write request.
// Apply it to any http.Handler whose mutations should trigger SSE board updates.
func NotifyMiddleware(b *Broadcaster) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}
			nw := &notifyWriter{ResponseWriter: w}
			next.ServeHTTP(nw, r)
			if nw.status >= 200 && nw.status < 300 {
				b.NotifyAll()
			}
		})
	}
}
