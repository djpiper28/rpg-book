package parser 

//go:generate echo Generating query language parser...

//go:generate go tool peg -inline -switch ./query.peg
