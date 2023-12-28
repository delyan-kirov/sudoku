package main

import (
	"delyan-kirov/sudoku/sudoku"
	"fmt"
)
func main() {
	fmt.Println("Generating sudoku. This will take about a minute.")
	sudoku.CreateSudoku()
}

