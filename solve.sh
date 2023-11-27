#!/usr/bin/bash

cd ./solutions/

rm -f *.solution 

conjure solve -ac --number-of-solutions=all --solver=minion ./sudoku.essence ./Params/example1.param
