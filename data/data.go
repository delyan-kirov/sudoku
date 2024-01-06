package data

import(
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"delyan-kirov/sudoku/sudoku"
)

const db_file = "./data/sudoku.db"

func Init_db() {

	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		fmt.Println("Could not extablish connection to the database")
		fmt.Println(err)
	} else {
		fmt.Printf("Database connection with %s was establish\n", db_file)
	}

	createTableStmt := `
	CREATE TABLE IF NOT EXISTS sudoku (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sudoku INTEGER
	);
	`

	_, err = db.Exec(createTableStmt)
	if err != nil {
		fmt.Printf("Could not create table schema %s\n", err)
	} else {
		fmt.Println("Schema created")
	}

	defer db.Close()
}

func print_data(new_sudoku sudoku.Sudoku) {
	sudoku.PrintSudoku(new_sudoku)
}
