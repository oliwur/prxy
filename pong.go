package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func pongHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	log.Println(string(b))
	w.Write(b)
}

func main() {
	addr := flag.String("a", ":8765", "the address to listen to")
	flag.Parse()

	log.Println("using listening addr ", *addr)

	http.HandleFunc("/", pongHandler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
