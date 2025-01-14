package service

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func Run() {
	var addr = flag.String("addr", ":5000", "The addr of the application.") // :5000
	flag.Parse()

	gomniauth.SetSecurityKey("SECURITY KEY")
	gomniauth.WithProviders(
		google.New(
			"624312092698-n5djffqsu9cocq008fr729osg7657l8m.apps.googleusercontent.com",
			"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",
			"http://localhost:5000/auth/callback/google",
		)) // :5000

	lmHost := os.Getenv("LINK_MANAGER_SERVICE_HOST")
	if lmHost == "" {
		lmHost = "localhost"
	}
	smHost := os.Getenv("SOCIAL_GRAPH_MANAGER_SERVICE_HOST")
	if smHost == "" {
		smHost = "localhost"
	}
	lmURL, err := url.Parse("http://" + lmHost + ":8080") // :8080
	if err != nil {
		log.Fatal(err)
	}
	smURL, err := url.Parse("http://" + smHost + ":9090")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", MustAuth(&templateHandler{filename: "mainpage.html"}))

	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/auth/login/", loginHandler)
	http.HandleFunc("/auth/callback/", callbackHandler)

	http.Handle("/links", jwtMiddleware(httputil.NewSingleHostReverseProxy(lmURL)))
	http.Handle("/followers", jwtMiddleware(httputil.NewSingleHostReverseProxy(smURL)))
	http.Handle("/following", jwtMiddleware(httputil.NewSingleHostReverseProxy(smURL)))
	http.Handle("/follow", jwtMiddleware(httputil.NewSingleHostReverseProxy(smURL)))
	http.Handle("/unfollow", jwtMiddleware(httputil.NewSingleHostReverseProxy(smURL)))

	log.Println("Starting api gateway service", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
