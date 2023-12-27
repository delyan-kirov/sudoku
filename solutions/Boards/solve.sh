#!/usr/bin/bash

conjure solve -ac --number-of-solutions=$1 --solver=minion ./sudoku.essence ./initial.param
