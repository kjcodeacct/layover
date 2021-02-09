package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {

	param := req.URL.Query().Get("msg")
	url := fmt.Sprintf("%v %v %v", req.Method, req.URL, req.Proto)
	log.Println(url)
	fmt.Fprintf(w, "%s\n", param)
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	log.Println("starting up on port 8090")
	http.HandleFunc("/test", hello)
	http.HandleFunc("/header-test", headers)

	http.ListenAndServe(":8090", nil)
}
