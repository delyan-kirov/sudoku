package main

import (
    "github.com/gin-gonic/gin"
  // "net/http"
)

func main() {
    r := gin.Default()

    // Serve static files from the build directory
    r.Static("/static", "./build/static")

    // Serve the React app's index.html
    r.GET("/", func(c *gin.Context) {
        c.File("./build/index.html")
    })

    port := ":8080"
    r.Run(port)
}
