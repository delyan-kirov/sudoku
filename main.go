package main

import (
	"delyan-kirov/sudoku/data"
	"delyan-kirov/sudoku/sudoku"
	"fmt"
)

func main() {
	fmt.Println("Initializing sudoku inside main")
	data.Init_db()

	test_sudoku, err := sudoku.ReadParam("./solutions/0.param")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Getting sudoku 0 from db")
	err = data.Write(test_sudoku)

	if err != nil {
		println(err)
	}

	sudoku_from_db, err := data.Read(1)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Sudoku from db is:")
	sudoku.PrintSudoku(sudoku_from_db)

	fmt.Println("Generating sudoku. This will take about a minute.")
	sudoku.CreateSudoku()
}
