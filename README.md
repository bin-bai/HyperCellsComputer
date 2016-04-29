# HyperCellsComputer
A research project of parallel computing, which consists massive ALU as hyper cells, and parallel memory bus.

Below example uses 2000 cells, which handle about 4000 integers.

The first example is in-place addition of O(logn):
```
[1 2 3 ... 3998 3999 4000]
[8002000 0 0 ... 0 0 0]

Tickcount is 188
Run 14 loops
```

The second example is naive in-place sorting of O(n):
```
[8384 1205 7869 ... 7876 9751 3306]
[1 1 3 ... 9988 9989 9996]

Tickcount is 51106
Run 1965 loops
```

Volunteers wanted
-----------------
Parallel computing is a 'game-changing' technology, could be considered as a Million-Human-Year project, all contributions are welcomed, including but not limit to:
```
	High level programming language and compiler
	Parallel algorithm and data structure
	Artificial Intelligence
	Software Emulator
	Operating System
	Network
	Graphic & Video
	CPU, Bus, and memory design
	Mainboard
	Peripheral
	...
	Applications
```

Getting started
---------------
Retrieve source code:
```
	go get -u github.com/bin-bai/HyperCellsComputer
```
