package protal_web

import (
	"context"
	"github.com/micro/go-micro/web"
	"leslack/src/GoMicroDemo/proto/sum"
	"net/http"
	"strconv"
)

var (
	srvClient sum.SumService
)

func main() {
	service := web.NewService(
		web.Name("go.micro.learning.web.portal"),
		web.Address(":202020"),
		web.StaticDir("html"),
	)
	service.Init()

	srvClient = sum.NewSumService("go.micro.learning.srv.sum", service.Options().Service.Client())
	service.HandleFunc("/learning/sum", Sum)
}

func Sum(w http.ResponseWriter, r *http.Request) {
	inputString := r.URL.Query().Get("input")
	input, _ := strconv.ParseInt(inputString, 10, 10)
	req := &sum.SumRequest{
		Input: input,
	}

	rsp, err := srvClient.GetSum(context.Background(), req)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(strconv.Itoa(int(rsp.Output))))

}
