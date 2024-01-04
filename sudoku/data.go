package sudoku

import(
	"database/sql"
	"fmt"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

const db_file = "./sudoku.db"

func main() {
	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("Database connection with %s was establish\n", db_file)
	}
	fmt.Printf("Creating database analytics\n")
	defer db.Close()
}
