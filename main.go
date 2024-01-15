package main

import (
	"delyan-kirov/sudoku/data"
	"delyan-kirov/sudoku/sudoku"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// Serve static files from the build directory
	r.Static("/static", "./build/static")

	r.GET("/", func(c *gin.Context) {
		c.File("./build/index.html")
	})

	r.POST("/check_solution", func(c *gin.Context) {
		var jsonData [9][9]int
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON format"})
			fmt.Printf("Error parsing JSON: %v\n", err)
			return
		}

		sudoku_board := sudoku.Sudoku(jsonData)
		sudoku.PrintSudoku(sudoku_board)

		responseData := gin.H{"message": "Solution checked successfully", "checkedSolution": "example"}
		c.JSON(200, responseData)
	})

	r.GET("/initial_board", func(c *gin.Context) {
		initialBoard, err := data.Read(32)
		if err != nil {
			fmt.Println("Could not read the sudoku")
			fmt.Printf("ERROR: %s\n", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(200, gin.H{"initialBoard": initialBoard})
	})

	port := ":8080"
	r.Run(port)
}
