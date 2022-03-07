package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)
//const VSBASEURL = "http://localhost:8080"



//func CreateDeviceConfig() {
//
//}

func CreateDeviceConfig(c *gin.Context) {
	// 接收的参数是gid+idx

}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}


func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(Logger())
	//r.GET("/handleDeviceRequest", HandleDeviceRequest)
	return r
}

func main() {
	//ldc.Init()
	r := gin.Default()
	r.Use(Logger())
	//r.POST("/handleDeviceRequest", HandleDeviceRequest)
	r.Run(":8081") // 监听并在 0.0.0.0:8081 上启动服务
}