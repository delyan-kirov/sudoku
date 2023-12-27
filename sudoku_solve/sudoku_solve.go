package sudoku_solve

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type Sudoku [9][9]int

func Say_hi() {
	fmt.Println("hello")
}

func genSudokuParam(sudoku Sudoku) string {
	const intRange = "int(1..9)"
	param := "letting initial be [ \n"
	for i, row := range sudoku {
		stringRow := ""
		for j, num := range row {
			if j%3 == 0 && j != 0 {
				stringRow = stringRow + "  "
			}
			if j == len(row)-1 {
				stringRow = stringRow + strconv.Itoa(num)
			} else {
				stringRow = stringRow + strconv.Itoa(num) + ", "
			}
		}
		if i%3 == 0 {
			param = param + "\n"
		}
		if i <= 7 {
			param = param + "   [ " + stringRow + "; " + intRange + "], \n"
		} else {
			param = param + "   [ " + stringRow + "; " + intRange + "]; \n"
		}
	}
	return "language Essence 1.3\n\n" + param + "\n   " + intRange + " ]"
}

func Solve_sudoku(sudoku Sudoku) (bool, error) {
	// - [ ] write the param file
	//  - [ ] run solve.sh
	//  - [ ] read the result
	//  - [ ] return bool
	sudoku_param := genSudokuParam(sudoku)
	const param_path = "./solve/sudoku.param"
	file, err := os.Create("./sudoku_solve.go")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return false, err
	}
	_, err = io.WriteString(file, sudoku_param)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return false, err
	}
	return true, nil
}
