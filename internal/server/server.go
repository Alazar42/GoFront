package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Options struct {
	Addr      string
	Root      string
	Dev       bool
	Logger    io.Writer
	Handler   http.Handler
	ReloadHub *ReloadHub
}

type ReloadHub struct {
	mu   sync.Mutex
	subs map[chan string]struct{}
}

func NewReloadHub() *ReloadHub { return &ReloadHub{subs: make(map[chan string]struct{})} }

func (h *ReloadHub) Subscribe() chan string {
	ch := make(chan string, 1)
	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *ReloadHub) Broadcast() {
	if h == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.subs {
		select {
		case ch <- "reload":
		default:
		}
	}
}

func Serve(opts Options) error {
	fs := http.FileServer(http.Dir(opts.Root))
	mux := http.NewServeMux()
	mux.Handle("/", rewriteIndex(fs, opts.Root, opts.Dev))
	if opts.Dev && opts.ReloadHub != nil {
		mux.HandleFunc("/__gofront_reload", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "streaming unsupported", http.StatusInternalServerError)
				return
			}
			ch := opts.ReloadHub.Subscribe()
			for {
				select {
				case <-r.Context().Done():
					return
				case <-ch:
					_, _ = fmt.Fprint(w, "data: reload\n\n")
					flusher.Flush()
				}
			}
		})
		mux.HandleFunc("/__gofront_error", func(w http.ResponseWriter, r *http.Request) {
			data, err := os.ReadFile(filepath.Join(opts.Root, ".gofront-error"))
			if err != nil && !os.IsNotExist(err) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write(data)
		})
	}
	server := &http.Server{Addr: opts.Addr, Handler: mux}
	if opts.Logger != nil {
		_, _ = fmt.Fprintf(opts.Logger, "GoFront serving %s on http://localhost%s\n", opts.Root, opts.Addr)
	}
	return server.ListenAndServe()
}

func rewriteIndex(next http.Handler, root string, dev bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".wasm") || strings.HasSuffix(r.URL.Path, ".js") {
			next.ServeHTTP(w, r)
			return
		}
		if r.URL.Path != "/" && filepath.Ext(r.URL.Path) != "" {
			next.ServeHTTP(w, r)
			return
		}
		indexPath := filepath.Join(root, "index.html")
		data, err := os.ReadFile(indexPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if dev {
			data = append([]byte("<script>window.__GOFRONT_DEV__=true;</script>\n"), data...)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	})
}
