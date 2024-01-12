package main

import (
	"fmt"
	// "delyan-kirov/sudoku/data"
	"delyan-kirov/sudoku/sudoku"
	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "net/http"
)

func main() {
	r := gin.Default()

	// Serve static files from the build directory
	r.Static("/static", "./build/static")

	r.GET("/", func(c *gin.Context) {
		c.File("./build/index.html")
	})

	// r.Use(cors.Default())

	r.POST("/check_solution", func(c *gin.Context) {
		var jsonData [9][9]int
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON format"})
			fmt.Printf("Error parsing JSON: %v\n", err)
			return
		}
		// fmt.Printf("Received JSON data: %v\n", jsonData)
		sudoku.PrintSudoku(sudoku.Sudoku(jsonData))
		responseData := gin.H{"message": "Solution checked successfully", "checkedSolution": "example"}
		c.JSON(200, responseData)
	})

	port := ":8080"
	r.Run(port)
}
