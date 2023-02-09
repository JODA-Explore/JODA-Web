package cond

type Type int

const (
	Exists Type = iota
	IsNull
	IsBool
	IsNumber
	IsString
	IsObject
	IsArray
	Equal
)

var condTypeSlice = []string{
	"EXISTS",
	"NULL",
	"BOOL",
	"NUMBER",
	"STRING",
	"OBJECT",
	"ARRAY",
	"Equal",
}

func (ct Type) String() string {
	return condTypeSlice[ct]
}
