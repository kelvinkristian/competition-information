package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

const (
	port = 8080
)

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// NullBool is an alias for sql.NullBool data type
type NullBool struct {
	sql.NullBool
}

// MarshalJSON for NullBool
func (nb *NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON for NullFloat64
func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// NullTime is an alias for mysql.NullTime data type
type NullTime struct {
	mysql.NullTime
}

// MarshalJSON for NullTime
func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

// Competition is and struct contains competition description
type Competition struct {
	cmpName                 string     `json:"Name"`
	cmpLastRegistrationDate NullTime   `json:"LastRegistrationDate"`
	cmpStartDate            NullTime   `json:"StartDate"`
	cmpPrizePool            NullInt64  `json:"PrizePool"`
	cmpDesc                 NullString `json:"Desc"`
	cmpImageSrc             NullString `json:"ImageSrc"`
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	mw := Chain(Logging(), Method("GET"), Tracing())
	http.Handle("/", mw(Index))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	fmt.Println("Connected to port " + strconv.Itoa(port) + ", Have a nice day!")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
