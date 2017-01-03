package root

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauthsvc "google.golang.org/api/oauth2/v2"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
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

func bibtekHandlerREF(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: r.URL.Scheme,
		Host:   "sbinet-xyz-bibtek.appspot.com",
	})
	proxy.Transport = client.Transport
	if u := user.Current(ctx); u != nil && false {
		proxy.Director = func(r *http.Request) {
			h := &r.Header
			h.Set("X-AppEngine-User-Email", u.Email)
			h.Set("X-AppEngine-Auth-Domain", u.AuthDomain)
			h.Set("X-AppEngine-User-Id", u.ID)
			if u.Admin {
				h.Set("X-AppEngine-User-Is-Admin", "1")
			} else {
				h.Set("X-AppEngine-User-Is-Admin", "0")
			}
			h.Set("X-AppEngine-Federated-Identity", u.FederatedIdentity)
			h.Set("X-AppEngine-Federated-Provider", u.FederatedProvider)
		}
	}
	proxy.ServeHTTP(w, r)
}

func bibtekHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	src, err := google.DefaultTokenSource(ctx, oauthsvc.UserinfoEmailScope)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: src,
			Base:   &urlfetch.Transport{Context: ctx},
		},
	}
	client = oauth2.NewClient(ctx, src)

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: r.URL.Scheme,
		Host:   "bibtek.sbinet.xyz",
	})
	proxy.Transport = client.Transport

	svc, err := oauthsvc.New(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ui, err := svc.Userinfo.Get().Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("--> userinfo: %v\n", ui.Email)

	proxy.ServeHTTP(w, r)
}
