package main

import (
	"log"

	limiter "github.com/davidleitw/gin-limiter"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func NewServer() *gin.Engine {
	server := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	limitControl, _ := limiter.DefaultController(rdb, "24-H", 21000, "debug")
	_ = limitControl.Add("/post1", "20-M", "post", 120)
	_ = limitControl.Add("/api/post2", "15-H", "post", 200)
	_ = limitControl.Add("/post3", "24-H", "post", 120)

	server.Use(limiter.LimitMiddle(limitControl))

	server.POST("/post1", post1) // /post1

	server.POST("api/post2", post2) // /api/post2

	server.POST("/post3", post3) // /post3

	return server
}

func post1(ctx *gin.Context) {
	ctx.String(200, ctx.FullPath())
}

func post2(ctx *gin.Context) {
	ctx.String(200, ctx.FullPath())
}

func post3(ctx *gin.Context) {
	ctx.String(200, ctx.ClientIP())
}

func main() {
	server := NewServer()

	err := server.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}
