package main

import (
	"delyan-kirov/sudoku/data"
	"delyan-kirov/sudoku/sudoku"
	"fmt"
)

func main2() {
	fmt.Println("Migrating database")
	err := data.Migrate()
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}
	fmt.Println("Generating sudoku. This will take about a minute.")
	sudoku.CreateSudoku()
}
