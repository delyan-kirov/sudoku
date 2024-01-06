package main

import (
	"delyan-kirov/sudoku/sudoku"
	"delyan-kirov/sudoku/data"
	"fmt"
)
func main() {
	data.Init_db()

	fmt.Println("Generating sudoku. This will take about a minute.")
	sudoku.CreateSudoku()
}

