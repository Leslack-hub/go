package main

import "github.com/gin-gonic/gin"

func main()  {
	engine := gin.Default()
	engine.GET("/hello", func(context *gin.Context){
		context.Writer.Write([]byte("hello gin\n"))
	})

	engine.Run()
}