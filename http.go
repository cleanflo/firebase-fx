package register

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
)

// HttpEntrypoint is the entrypoint for http functions
// it will route the request to the correct function
func (f *FunctionRegistrar) HttpEntrypoint(w http.ResponseWriter, r *http.Request) {
	f.http.ServeHTTP(w, r)
}

// RegisterHTTP allows a user to register a HttpFunction without using mux.Router
// Using this method, functions cannot be inlined and must be exported (ie. capitalized) by the package.
// The deployment script for functions registered via this method will be in the below format:
//   ~$ gcloud functions deploy  <FunctionName> --trigger-http --allow-unauthenticated
// and can be executed like:
//   ~$ curl "https://REGION-PROJECT_ID.cloudfunctions.net/FunctionName"
func (f *FunctionRegistrar) RegisterHTTP(handler http.HandlerFunc) *HttpFunction {
	fpath := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

	fps := strings.Split(fpath, ".")
	// // fpkg := fps[0]
	fname := fps[len(fps)-1]

	// r := f.http.HandleFunc(fname, handler)
	fn := &HttpFunction{
		reg: f,
		// r:    r,
		fn:   handler,
		path: fname,
	}

	f.handlers[fname] = fn

	return fn
}

// HTTP registers the given path/name to an underlying mux.Router
// add advanced routing options (Headers, Methods) to the returned HttpFunction
// The deployment script for functions registered via this method will be in the below format:
//   ~$ gcloud functions deploy Register.HttpEntrypoint --trigger-http --allow-unauthenticated
// and can be executed like:
//   ~$ curl "https://REGION-PROJECT_ID.cloudfunctions.net/Register.HttpEntrypoint/{path}"
func (f *FunctionRegistrar) HTTP(path string, handler http.HandlerFunc) *HttpFunction {
	r := f.http.HandleFunc(path, handler)

	r = r.Name(path)
	fn := &HttpFunction{
		reg:  f,
		r:    r,
		fn:   handler,
		path: path,
	}

	f.handlers[path] = fn

	return fn
}

// Middleware registers the given mux.MiddlewareFunc to the underlying mux.Router
func (f *FunctionRegistrar) MiddleWare(wares ...mux.MiddlewareFunc) {
	f.http.Use(wares...)
}

func (f *FunctionRegistrar) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.http.ServeHTTP(w, r)
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
