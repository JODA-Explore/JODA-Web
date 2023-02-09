package jsontype

type Type int

const (
	Invalid Type = iota
	String
	Int
	Float
	Bool
	Object
)
