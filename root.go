package root

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/bibtek", bibtekHandler)

	http.Handle("/", r)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world! (sbinet.xyz)\n")
	ctx := appengine.NewContext(r)
	if u := user.Current(ctx); u != nil {
		fmt.Fprintf(w, "user: %q\n", u.String())
	}
}

func bibtekHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.URL.Scheme+"//bibtek.sbinet.xyz", http.StatusMovedPermanently)
}
