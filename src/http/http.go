package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func main() {
	request, err := http.NewRequest(http.MethodGet, "http://www.douyu.com", nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Mobile Safari/537.36")
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println("request", req)
            return nil
		},
	}
	resp, err := client.Do(request)
	//request, err := http.Get("http://www.baidu.com")
	if err != nil {
		panic(err)
	}
    defer resp.Body.Close()

	s, err := httputil.DumpResponse(resp,true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", s)
}
