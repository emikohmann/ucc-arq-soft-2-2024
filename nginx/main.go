package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type Body struct {
	// json tag to serialize json body
	Name string `json:"name"`
}

func main() {
	engine := gin.New()
	engine.GET("/test", func(context *gin.Context) {
		body := Body{}
		body.Name = os.Getenv("HOSTNAME")
		// using BindJson method to serialize body with struct body.Name=os.Getenv("HOSTNAME")
		context.JSON(http.StatusAccepted, &body)
	})
	engine.Run(":3000")
}
