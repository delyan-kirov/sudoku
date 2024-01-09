package data

import (
	"database/sql"
	"delyan-kirov/sudoku/sudoku"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const db_file = "./data/sudoku.db"

func Init_db() error {
	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		fmt.Printf("Could not extablish connection to the database\n")
		return err
	} else {
		fmt.Printf("Database connection with %s was established\n", db_file)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS sudoku (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sudoku TEXT
	);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Printf("Could not create table schema for %s\n", db_file)
		return err
	} else {
		fmt.Println("Schema created")
	}

	err = db.Close()
	if err != nil {
		fmt.Printf("Database %s could not be closed\n", db_file)
		return err
	}

	return nil
}

func Write(sudoku_board sudoku.Sudoku) error {
	sudoku_str := ""
	for _, row := range sudoku_board {
		for _, row_entry := range row {
			sudoku_str += strconv.Itoa(row_entry)
		}
	}
	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		fmt.Printf("ERROR: Could not connect to %s\n", db_file)
		return err
	}
	_, err = db.Exec("INSERT INTO sudoku (sudoku) VALUES (?)", sudoku_str)
	if err != nil {
		fmt.Printf("Could not write to table %s\n", db_file)
		return err
	}
	err = db.Close()
	if err != nil {
		fmt.Printf("Coulf not close database %s\n", db_file)
		return err
	}
	return nil
}

func Read(id int) (sudoku.Sudoku, error) {
	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		fmt.Printf("ERROR: Could not connect to %s\n", db_file)
		return sudoku.InitSudoku(), err
	}

	sudoku_query := db.QueryRow("SELECT sudoku FROM sudoku WHERE id = ?", id)
	sudoku_str := ""
	err = sudoku_query.Scan(&sudoku_str)

	if err != nil {
		fmt.Printf("ERROR: No matches for the provided index inside the database.\nProbable cause: index out of bound\n")
		return sudoku.InitSudoku(), err
	}

	if len(sudoku_str) != 81 {
		return sudoku.InitSudoku(),
			fmt.Errorf(
				`ERROR: Invalid Sudoku string length
				 The string received was: %s
				 Probable cause: incorect conversion`,
				sudoku_str)
	}

	var sudoku_board sudoku.Sudoku
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			// Convert sudoku to string
			num, err := strconv.Atoi(string(sudoku_str[i*9+j]))
			if err != nil {
				return sudoku_board, fmt.Errorf("Error converting character to string: %v", err)
			}
			sudoku_board[i][j] = num
		}
	}

	return sudoku_board, nil
}
