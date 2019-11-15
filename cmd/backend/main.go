package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	chitrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Backend struct {
	client *http.Client
}

func NewBackend() *Backend {
	return &Backend{
		client: http.DefaultClient,
	}
}

func (b *Backend) ping(w http.ResponseWriter, r *http.Request) {
	var seq int
	var err error

	seqParam := r.URL.Query().Get("q")
	seq, err = strconv.Atoi(seqParam)
	if err != nil {
		fmt.Fprintln(w, errors.Wrap(err, "error parsing URL param"))
		return
	}

	if seq == 0 {
		fmt.Fprintf(w, "pong")
		return
	}

	reqURL := fmt.Sprintf("http://backend/ping?q=%d", seq-1)

	outReq, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		fmt.Fprintln(w, errors.Wrap(err, "error creating backend request"))
		return
	}

	resp, err := b.client.Do(outReq.WithContext(r.Context()))
	if err != nil {
		fmt.Fprintln(w, errors.Wrap(err, "error executing backend request"))
		return
	}

	defer resp.Body.Close()
	if _, err := io.Copy(w, resp.Body); err != nil {
		fmt.Fprintln(w, errors.Wrap(err, "error streaming backend response"))
		return
	}
}

func routes(b *Backend) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(chitrace.Middleware(chitrace.WithServiceName("backend")))

	r.Get("/ping", b.ping)

	return r
}

func main() {
	tracer.Start()
	defer tracer.Stop()

	svc := NewBackend()
	if err := http.ListenAndServe(":80", routes(svc)); err != nil {
		log.Fatal(err)
	}
}
