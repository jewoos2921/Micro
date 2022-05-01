package main

import (
	"flag"
	"log"
	"net/http"
)

// 모든 투표를 보여주는 index.html
// 특정 투표 결과를 보여주는 view.html
// 사용자가 새로운 투표를 만들 수 있는 new.html

func main() {
	var addr = flag.String("addr", ":8081", "website address")
	flag.Parse()
	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/",
		http.FileServer(http.Dir("public"))))
	log.Println("Serving website at:", *addr)
	http.ListenAndServe(*addr, mux)
}
