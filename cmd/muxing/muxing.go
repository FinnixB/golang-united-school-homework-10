package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func helloName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["PARAM"]
	fmt.Fprintf(w, "Hello, %s!", param)
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func dataProcess(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = w.Write([]byte(fmt.Sprintf("I got message:\n%v", string(b))))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func headersSum(w http.ResponseWriter, r *http.Request) {
	h := r.Header

	a, err := strconv.Atoi(h.Get("a"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	b, err := strconv.Atoi(h.Get("b"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("a+b", strconv.Itoa(a+b))
}

func defaultsHttp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Start(host string, port int) {
	router := mux.NewRouter()

	router.HandleFunc("/name/{PARAM}", helloName).Methods(http.MethodGet)
	router.HandleFunc("/bad", badRequest).Methods(http.MethodGet)
	router.HandleFunc("/data", dataProcess).Methods(http.MethodPost)
	router.HandleFunc("/headers", headersSum).Methods(http.MethodPost)
	router.PathPrefix("/").HandlerFunc(defaultsHttp)

	log.Println(fmt.Printf("Starting API server on %s:%d\n", host, port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router); err != nil {
		log.Fatal(err)
	}
}

//main /** starts program, gets HOST:PORT param and calls Start func.
func main() {
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8081
	}
	Start(host, port)
}
