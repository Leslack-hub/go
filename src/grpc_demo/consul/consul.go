package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Println("hello web3! this is n3 or n2")
	fmt.Fprintf(w, "hello web3! this is n3 or n2")
}

func healthHandler(w http.ResponseWriter, r * http.Request){
	fmt.Println("health check! n3 or n2")
}
func main()  {
	http.HandleFunc("/",handler)
	http.HandleFunc("/health", healthHandler)
	http.ListenAndServe(":10000", nil)
}
