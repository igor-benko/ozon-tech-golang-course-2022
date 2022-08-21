package main

type Package struct {
	Name    string
	Structs []Struct
}

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}
