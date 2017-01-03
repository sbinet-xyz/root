package root

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/x/", goGetRepoHandler)
	http.HandleFunc("/login", loginHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world! (sbinet.xyz)\n")
	ctx := appengine.NewContext(r)
	if u := user.Current(ctx); u != nil {
		fmt.Fprintf(w, "user: %q\n", u.String())
	}
	fmt.Fprintf(w, "\n\nscheme: %q\n", r.URL.Scheme)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if u == nil {
		url, _ := user.LoginURL(ctx, "/")
		fmt.Fprintf(w, `<a href="%s">Sign in</a>`, url)
		return
	}
	url, _ := user.LogoutURL(ctx, "/")
	fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
}

func goGetRepoHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	url = url[len("/x/"):]
	repo := url
	if strings.Contains(url, "/") {
		repo = filepath.Dir(url)
	}

	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="content-type" content="text/html; charset=utf-8">
  <meta name="go-import" content="sbinet.xyz/x/%[1]s git https://github.com/sbinet/%[1]s">
  <meta name="go-source" content="sbinet.xyz/x/%[1]s https://github.com/sbinet https://github.com/sbinet/%[1]s/tree/master{/dir} https://github.com/sbinet/%[1]s/blob/master{/dir}/{file}#L{line}">
  <meta http-equiv="refresh" content="0; url=https://godoc.org/sbinet.xyz/x/%[1]s">
</head>
<body>
Nothing to see here; <a href="https://godoc.org/sbinet.xyz/x/%[1]s">move along</a>.
</body>
</html>
`,
		repo,
	)
}
