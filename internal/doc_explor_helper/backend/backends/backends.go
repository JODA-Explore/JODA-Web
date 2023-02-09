package backends

import (
	"html/template"
	"strings"
)

type ID int

const (
	ValueType ID = iota
	StructureDiff
	Distinct
	Array
	Objects
)

type info struct {
	name string
	desc template.HTML
}

var infos = [...]info{
	{"Value Type", "Value Type"},
	{"Structure Difference", "Structure Difference"},
	{"Distinct Values", "Distinct Values"},
	{"Array Analysis", "Array Analysis"},
	{"Array of Objects", "Array of Objects"},
}

const Num = len(infos)

func Names() (names []string) {
	names = make([]string, Num)
	for i, x := range infos {
		names[i] = x.name
	}
	return
}

func (id ID) Name() string {
	return infos[id].name
}

func (id ID) Desc() template.HTML {
	return infos[id].desc
}

type IDs []ID

func (ids IDs) Main() ID {
	return ids[0]
}

func (ids IDs) Names() (names []string) {
	names = make([]string, len(ids))
	for i, id := range ids {
		names[i] = id.Name()
	}
	return
}

func (ids IDs) NamesDesc() template.HTML {
	return template.HTML(strings.Join(ids.Names(), " && "))
}

type Temp [Num][]int
