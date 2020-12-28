package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func (w http.ResponseWriter, r* http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"msg": "hello, world"}`))
	})

	http.HandleFunc("/healthz", func (w http.ResponseWriter, r* http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"msg": "ok"}`))
	})

	log.Print("starting server on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}


}
