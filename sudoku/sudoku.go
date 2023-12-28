package sudoku

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Sudoku [9][9]int

// Functions for pretty printing

func printBlue(printStr string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 34, printStr)
	fmt.Printf(colored)
}

func printBlueLn(printStr string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 34, printStr)
	fmt.Printf(colored)
	fmt.Println("")
}

func clearConsole() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func PrintSudoku(sudoku Sudoku) {
	const space = "    "
	printBlueLn(space + "|-----------+-----------+-----------|")
	for i := 0; i < 9; i++ {
		fmt.Print(space)
		for j := 0; j < 9; j++ {
			if j%3 != 0 {
				printBlue(" ")
			} else {
				printBlue("|")
			}
			fmt.Printf("%2d ", sudoku[i][j])
		}
		printBlue("|")
		fmt.Println("")
		if (i+1)%3 == 0 {
			printBlueLn(space + "|-----------+-----------+-----------|")
		} else {
			printBlueLn(space + "|           |           |           |")
		}
	}
}

// Solving the sudoku

func initSudoku() Sudoku {
	var sudoku Sudoku
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			sudoku[i][j] = 0
		}
	}
	return (sudoku)
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
	return "language Essence 1.3\n\n" + param + "\n   " + intRange + " ]\n"
}

func readParam(filePath string) (Sudoku, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return initSudoku(), errors.New("Error opening file")
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return initSudoku(), errors.New("Error reading file")
	}
	param := string(content)
	param = strings.ReplaceAll(param, "language Essence 1.3", "")
	param = strings.ReplaceAll(param, "letting initial be", "")
	param = strings.ReplaceAll(param, "; int(1..9)", "")
	param = strings.ReplaceAll(param, "int(1..9)", "")
	param = strings.ReplaceAll(param, ";", "")
	param = strings.ReplaceAll(param, " ", "")
	param = strings.ReplaceAll(param, "\n", "")

	var paramArray Sudoku
	err = json.Unmarshal([]byte(param), &paramArray)
	if err != nil {
		return initSudoku(), errors.New("Conversion to go array fail")
	}
	return paramArray, nil
}

func writeParam(sudoku Sudoku) (string, error) {
	content := genSudokuParam(sudoku)
	const paramPath = "./solutions/"
	paramFiles, err := filepath.Glob(paramPath + "*")
	keyIndex := len(paramFiles)
	if err != nil {
		return "", err
	}
	newParamPath := paramPath + strconv.Itoa(keyIndex) + ".param"
	newParamFile, err := os.Create(newParamPath)
	if err != nil {
		return "", errors.New("Could not create file")
	}
	_, err = newParamFile.Write([]byte(content))
	if err != nil {
		return "", err
	}
	defer newParamFile.Close()
	return newParamPath, nil
}

func IsValidSudoku(sudoku Sudoku) bool {
	checkRepeats := func(arr []int) bool {
		seen := make(map[int]bool)
		for _, num := range arr {
			if num == 0 {
				continue
			}
			if seen[num] {
				return false
			}
			seen[num] = true
		}
		return true
	}

	extractBlocks := func(sudoku Sudoku) [][]int {
		blocks := make([][]int, 0, 9)
		for rowStart := 0; rowStart < 9; rowStart += 3 {
			for colStart := 0; colStart < 9; colStart += 3 {
				block := make([]int, 0, 9)
				for i := rowStart; i < rowStart+3; i++ {
					for j := colStart; j < colStart+3; j++ {
						block = append(block, sudoku[i][j])
					}
				}
				blocks = append(blocks, block)
			}
		}
		return blocks
	}

	blocks := extractBlocks(sudoku)
	cols := make([]int, 9)
	for i := 0; i < 9; i++ {
		rows := sudoku[i][:]
		if !checkRepeats(rows) ||
			!checkRepeats(cols) ||
			!checkRepeats(blocks[i]) {
			return false
		}
	}
	return true
}

func solve_sudoku(sudoku Sudoku) (int, error) {
	sudoku_param := genSudokuParam(sudoku)
	// Write param to file
	file, err := os.Create("./.solve/sudoku.param")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return 0, err
	}
	_, err = io.WriteString(file, sudoku_param)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return 0, err
	}
	// solve param
	cmd := exec.Command("bash", "./solve.sh")
	cmd.Dir = "./.solve"
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing script:", err)
		return 0, err
	}
	// count solutions
	count_solutions := 0
	err = filepath.Walk("./.solve/conjure-output/",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".solution") {
				count_solutions++
			}
			return nil
		})
	if err != nil {
		return 0, err
	}
	// clear cached files
	cmd = exec.Command("bash", "./clear.sh")
	cmd.Dir = "./.solve"
	cmd.Stderr = os.Stderr
	cmd.Run()
	return count_solutions, nil
}

// Generating the sudoku puzzle

func gen_rand_sudoku(curr_sudoku Sudoku) Sudoku {
	// Make random assignment of an empty sudoku cell
	sudoku_row := rand.Intn(9)
	sudoku_col := rand.Intn(9)
	value_to_add := rand.Intn(9) + 1
	new_sudoku := curr_sudoku
	new_sudoku[sudoku_row][sudoku_col] = value_to_add
	is_free_pos := curr_sudoku[sudoku_row][sudoku_col] == 0
	if !IsValidSudoku(new_sudoku) {
		return gen_rand_sudoku(curr_sudoku)
	}
	if is_free_pos {
		num_sudoku_sols, err := solve_sudoku(new_sudoku)
		if err != nil {
			fmt.Println("Error in solution process")
			return (initSudoku())
		}
		if num_sudoku_sols == 0 {
			return gen_rand_sudoku(curr_sudoku)
		}
		if num_sudoku_sols == 1 {
			return new_sudoku
		}
		return gen_rand_sudoku(new_sudoku)
	}

	return gen_rand_sudoku(curr_sudoku)
}

func CreateSudoku() {
	sudoku := gen_rand_sudoku(initSudoku())
	PrintSudoku(sudoku)
	writeParam(sudoku)
}
