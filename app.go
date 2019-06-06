package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Competition is and struct contains competition description
type Competition struct {
	Name      string `json:"Name"`
	Date      string `json:"Date"`
	PrizePool int    `json:"PrizePool"`
	Desc      string `json:"Desc"`
}

// Middleware is use for chaining handler function
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging function for logging every API hits
func Logging() Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var start = time.Now()
			defer func() {
				log.Println(r.URL.Path, time.Since(start))
			}()
			handler(w, r)
		}
	}
}

// Method is a function to check whether the coresponse route is having the set method
func Method(m string) Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			handler(w, r)
		}
	}
}

// Tracing is a function to trace the current url location
func Tracing() Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Tracing request for %s\n", r.RequestURI)
			handler(w, r)
		}
	}
}

// Chain is a function to chain all the middleware
func Chain(middlewares ...Middleware) Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for _, middleware := range middlewares {
				handler = middleware(handler)
			}
			handler(w, r)
		}
	}
}

// Index is a HandlerFunction for root url about all competition
func Index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func main() {
	mw := Chain(Logging(), Method("GET"), Tracing())
	http.Handle("/", mw(Index))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	port := 8080
	fmt.Println("Connected to port " + strconv.Itoa(port) + ", Have a nice day!")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
