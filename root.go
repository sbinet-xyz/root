package root

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/bibtek", bibtekHandler)

	http.Handle("/", r)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world! (sbinet.xyz)\n")
}

func bibtekHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "sbinet-xyz-bibtek.appspot.com",
	})
	proxy.Transport = client.Transport
	proxy.ServeHTTP(w, r)
}
