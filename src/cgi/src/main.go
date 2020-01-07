package main // 代码包声明语句。
//系统包用来输出的
import (
	"fmt"
	"mock"
	real2 "real"
	"time"
)

type Retriever interface {
	Get(Url string) string
}

func download(r Retriever) string {
	return r.Get("http://www.douyu.com")
}
func main() {
	// 打印函数调用语句。用于打印输出信息。
	var r Retriever
	r = mock.Retriever{"xxxx"}
	inspect(r)
	r = &real2.Retriever{
		UserAgent: "http://www.huya.com",
		TimeOut:   time.Minute,
	}
	inspect(r)
	if realRetriever, ok := r.(mock.Retriever); ok {
		fmt.Println(realRetriever.Content)
	} else {
		fmt.Println("not a mock assertion")
	}
}

func inspect(r Retriever) {
	fmt.Printf("%T %v \n", r, r)
	switch v := r.(type) {
	case mock.Retriever:
		fmt.Println("contents :", v.Content)
	case *real2.Retriever:
		fmt.Println("contents", v)
	}
}
