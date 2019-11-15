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

type Frontend struct {
	client *http.Client
}

func NewFrontend() *Frontend {
	return &Frontend{
		client: http.DefaultClient,
	}
}

func (f *Frontend) ping(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	var seq int
	if q != "" {
		seq, _ = strconv.Atoi(q)
	}

	reqURL := fmt.Sprintf("http://backend/ping?q=%d", seq)

	outReq, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		fmt.Fprintln(w, errors.Wrap(err, "error creating backend request"))
		return
	}

	resp, err := f.client.Do(outReq.WithContext(r.Context()))
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

func routes(f *Frontend) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(chitrace.Middleware(chitrace.WithServiceName("frontend")))

	r.Get("/ping", f.ping)

	return r
}

func main() {
	tracer.Start()
	defer tracer.Stop()

	f := NewFrontend()

	log.Print("starting HTTP server")
	if err := http.ListenAndServe(":80", routes(f)); err != nil {
		log.Fatal(err)
	}
}
