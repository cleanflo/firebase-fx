package register

import (
	"net/http"

	"github.com/gorilla/mux"
)

// HttpEntrypoint is the entrypoint for http functions
// it will route the request to the correct function
func (f *FunctionRegistrar) HttpEntrypoint(w http.ResponseWriter, r *http.Request) {
	f.http.ServeHTTP(w, r)
}

// HTTP registers the given path to an underlying mux.Router
// add advanced routing options (Headers, Methods) to the returned HttpFunction
func (f *FunctionRegistrar) HTTP(path string, handler http.HandlerFunc) *HttpFunction {
	r := f.http.HandleFunc(path, handler)
	fn := &HttpFunction{
		reg:  f,
		r:    r,
		fn:   handler,
		path: path,
	}
	return fn
}

// Middleware registers the given mux.MiddlewareFunc to the underlying mux.Router
func (f *FunctionRegistrar) MiddleWare(wares ...mux.MiddlewareFunc) {
	f.http.Use(wares...)
}

// HttpFunction is a wrapper for mux.Route and the parent FunctionRegistrar
type HttpFunction struct {
	reg             *FunctionRegistrar
	r               *mux.Route
	unauthenticated bool
	path            string
	fn              http.HandlerFunc
}

// Unauthenticated marks the function as --allow-unauthenticated
func (h *HttpFunction) Unauthenticated(t bool) *HttpFunction {
	h.unauthenticated = t
	return h
}

// Methods registers the given methods to the underlying mux.Route
func (h *HttpFunction) Methods(methods ...string) *HttpFunction {
	h.r.Methods(methods...)
	return h
}

// Headers registers the given headers to the underlying mux.Route
func (h *HttpFunction) Headers(pairs ...string) *HttpFunction {
	h.r.Headers(pairs...)
	return h
}

// Host registers the given host to the underlying mux.Route
func (h *HttpFunction) Host(host string) *HttpFunction {
	h.r.Host(host)
	return h
}

// Queries registers the given queries to the underlying mux.Route
func (h *HttpFunction) Queries(queries ...string) *HttpFunction {
	h.r.Queries(queries...)
	return h
}
