package parser

//go:generate echo Generating query language parser...

//go:generate go tool peg -strict -inline -switch ./query.peg
