package main

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
	_ "strconv"
	"strings"
	"time"
)

type Sudoku [9][9]int

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

func printBlue(printStr string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 34, printStr)
	fmt.Printf(colored)
}

func printBlueLn(printStr string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 34, printStr)
	fmt.Printf(colored)
	fmt.Println("")
}

func randomPermutation(numbers []int) []int {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Use Fisher-Yates shuffle algorithm to permute the slice
	for i := len(numbers) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}

	return numbers
}

// PrintMatrix prints the matrix to the console.
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

func countSolutions(param string) (int, error) {
	const solutionPath = "./solutions/"
	_, err := os.Stat(solutionPath + "Params/" + param)
	if err != nil {
		return 0, os.ErrNotExist
	}
	cmd := exec.Command("bash", "solve.sh", param)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return 0, err
	}

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Command finished with error:", err)
		return 0, err
	}

	solutions, err := filepath.Glob(solutionPath + "*.solution")
	if err != nil {
		return 0, err
	}

	return len(solutions), nil
}

func writeParam(sudoku Sudoku) (string, error) {
	content := genSudokuParam(sudoku)
	const paramPath = "./solutions/Params/"
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

func readSolution(solutionPath string) (string, error) {
	file, err := os.Open(solutionPath)
	if err != nil {
		return "", errors.New("Error opening file")
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", errors.New("Error reading file")
	}
	return string(content), nil
}

func solve_sudoku(sudoku Sudoku) (int, error) {
	sudoku_param := genSudokuParam(sudoku)
	// Write param to file
	file, err := os.Create("./solve/sudoku.param")
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
	cmd.Dir = "./solve"
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing script:", err)
		return 0, err
	}
	// count solutions
	count_solutions := 0
	err = filepath.Walk("./solve/conjure-output/", func(path string, info os.FileInfo, err error) error {
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
	return count_solutions, nil
}

func main() {
	// clearConsole()
	var mysudoku Sudoku = initSudoku()
	// PrintSudoku(mysudoku)
	// fmt.Println("")
	// fmt.Println("Random Permutation:", randomPermutation([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}))
	// time.Sleep(50 * time.Millisecond)
	// fmt.Println(genSudokuParam(mysudoku))
	// fmt.Println(readParam("./solutions/Params/example1.param"))
	var err error
	mysudoku, err = readParam("./solutions/Params/example1.param")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// fmt.Println(genSudokuParam(mysudoku))
	// fmt.Println(writeParam(initSudoku()))
	// fmt.Println(countSolutions("example1.param"))
	// fmt.Println(readSolution("./solutions/Boards/sudoku-initial-000288.solution"))
	fmt.Println(solve_sudoku(mysudoku))
}

// algorithm
// 1. Start with a filled board
// 2. Initialize a permutation index matrix
// 3. Initialize an array with assigned null spaces
// 2. Assign zero from the random index matrix
// 3. Check for unique solutions
// 4. If unique - add to assined array
//    else create a new permutation index matrix by removing the indices assigned
// 5. If assigned indices is full, stop

// TODO
// - [X] Create a parser to generate and read param files
// - [ ] Make a function to solve a sudoku board from go
// - [X] Create a function that counts the number of solutions
// - [ ] Make it so that the initial block is randomly filled with numbers that work
// // // - Algorithm
// // // // - For each row - row a dice and decide to fill or not
// // // // - If decided to fill, generate a random digit
