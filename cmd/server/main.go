package main

import (
	"crypto/subtle"
	"encoding/json"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

func main() {
	data, err := ioutil.ReadFile("asset")
	if err !=nil {
		//panic(err)
	}
	var results shared.CloudDiscoveryResults
	if err := json.Unmarshal(data, &results); err != nil {
		//panic(err)
	}
	http.HandleFunc("/logo.png", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("web/logo.png")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}))
	http.HandleFunc("/", BasicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(shared.CloudDiscoveryPage(results))
	}), "admin", "123456", "Please enter your username and password for this site"))
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.DefaultServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
