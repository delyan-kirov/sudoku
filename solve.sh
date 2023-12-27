#!/usr/bin/bash

cd ./solutions/

#rm -f *.solution 

param="./Params/$1"

conjure solve -ac --number-of-solutions=20 --solver=minion ./sudoku.essence $param

# mv *.solutions ./solutions/Boards/
