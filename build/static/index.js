document.addEventListener("DOMContentLoaded", function () {
  let currentBoard = [
    [0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0],
  ];

  fetch("/initial_board")
    .then((response) => response.json())
    .then((data) => {
      const serverInitialBoard = data.initialBoard;

      // Update the currentBoard with the fetched data
      currentBoard = serverInitialBoard.map((row) =>
        row.map((cell) => ({
          value: cell,
          editable: cell === 0,
        }))
      );
      console.log(currentBoard);

      // Update the UI with the fetched data
      renderBoard(currentBoard);
    })
    .catch((error) => {
      console.error("Error fetching initial board data:", error);
    });

  const boardContainer = document.getElementById("sudoku-board");
  const checkResultButton = document.getElementById("check-result");

  function renderBoard(board) {
    const table = document.createElement("table");

    board.forEach((row, rowIndex) => {
      const tr = document.createElement("tr");

      row.forEach((cell, colIndex) => {
        const td = document.createElement("td");

        const input = document.createElement("input");
        input.type = "number";
        input.min = "1";
        input.max = "9";
        input.value = cell.value === 0 ? "" : String(cell.value);
        input.dataset.row = String(rowIndex);
        input.dataset.col = String(colIndex);

        input.addEventListener("input", function () {
          const row = parseInt(this.dataset.row, 10);
          const col = parseInt(this.dataset.col, 10);
          const value = parseInt(this.value, 10);

          if ((value >= 1 && value <= 9) || this.value === "") {
            if (board[row][col].editable) {
              handleCellChange(row, col, value || 0);
            } else {
              this.value = String(board[row][col].value);
            }
          } else {
            this.value = String(board[row][col].value);
          }
        });

        td.appendChild(input);
        tr.appendChild(td);

        if ((colIndex + 1) % 3 === 0 && colIndex !== 8) {
          const separatorCell = document.createElement("td");
          separatorCell.classList.add("separator");
          tr.appendChild(separatorCell);
        }
      });

      table.appendChild(tr);

      if ((rowIndex + 1) % 3 === 0 && rowIndex !== 8) {
        const separatorRow = document.createElement("tr");
        const separatorCell = document.createElement("td");
        separatorCell.colSpan = 9;
        separatorCell.classList.add("separator");
        separatorRow.appendChild(separatorCell);
        table.appendChild(separatorRow);
      }
    });

    boardContainer.innerHTML = "";
    boardContainer.appendChild(table);
  }

  function handleCellChange(row, col, value) {
    currentBoard[row][col].value = value;
  }

  checkResultButton.addEventListener("click", function () {
    const current_board = currentBoard.map((row) =>
      row.map((cell) => cell.value)
    );
    const jsonData = JSON.stringify(current_board);
    console.log("JSON Data:", jsonData);

    fetch("/check_solution", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: jsonData,
    })
      .then((response) => response.json())
      .then((data) => {
        console.log("Server response:", data);
        // Handle the data returned by the server here
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  });
});

// TODO:
// - [ ] Add an index so that the user can chose to solve another sudoku 
// - [ ] Make it so that the database has the solutions
// - [ ] Integrate the database init into the server logic
// - [ ] Separate the javascript code
