document.addEventListener("DOMContentLoaded", function () {
  let currentBoard = initializeBoard();

  fetchInitialBoardData()
    .then((serverInitialBoard) => {
      updateCurrentBoard(serverInitialBoard);
      renderBoard(currentBoard);
    })
    .catch((error) => {
      console.error("Error fetching initial board data:", error);
    });

  const boardContainer = document.getElementById("sudoku-board");
  const checkResultButton = document.getElementById("check-result");

  function initializeBoard() {
    return Array.from(
      { length: 9 },
      () => Array(9).fill({ value: 0, editable: true }),
    );
  }

  async function fetchInitialBoardData() {
    const response = await fetch("/initial_board");
      const data = await response.json();
      return data.initialBoard;
  }

  function updateCurrentBoard(serverInitialBoard) {
    currentBoard = serverInitialBoard.map((row) =>
      row.map((cell) => ({ value: cell, editable: cell === 0 }))
    );
  }

  function renderBoard(board) {
    const table = document.createElement("table");

    board.forEach((row, rowIndex) => {
      const tr = document.createElement("tr");

      row.forEach((cell, colIndex) => {
        const td = document.createElement("td");

        const input = document.createElement("input");
        input.type = "number";
        input.min = "1";
        input.maxLength = 1;
        input.pattern = [1 - 9];
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
    if (currentBoard[row][col].editable) {
      currentBoard[row][col].value = value;
    }
  }

  function checkResult() {
    const current_board = currentBoard.map((row) =>
      row.map((cell) => cell.value)
    );
    const jsonData = JSON.stringify(current_board);

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
  }

  // Event listeners
  boardContainer.addEventListener("input", function (event) {
    const inputElement = event.target;
    const row = parseInt(inputElement.dataset.row, 10);
    const col = parseInt(inputElement.dataset.col, 10);
    const value = parseInt(inputElement.value, 10);
    handleCellChange(row, col, value || 0);
  });

  checkResultButton.addEventListener("click", checkResult);
});

// TODO:
// - [ ] Create a demo in React 
// - [ ] Add an index so that the user can chose to solve another sudoku
// - [ ] Make it so that the database has the solutions
// - [ ] Integrate the database init into the server logic
// - [ ] Separate the javascript code
