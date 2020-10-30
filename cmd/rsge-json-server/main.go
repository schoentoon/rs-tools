package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

var port = flag.Int("port", 8000, "http port")

type server struct {
	ItemCache map[int64]string // TODO this needs a mutex
	Client    *http.Client
}

func main() {
	flag.Parse()

	s := server{
		ItemCache: make(map[int64]string),
	}

	// initialize routes, and start http server
	http.HandleFunc("/search", cors(s.search))
	http.HandleFunc("/query", cors(s.query))
	http.HandleFunc("/", cors(func(w http.ResponseWriter, r *http.Request) {
		data, err := httputil.DumpRequest(r, true)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", data)
	}))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}
