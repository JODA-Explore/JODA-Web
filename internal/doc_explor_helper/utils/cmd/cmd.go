package cmd

import (
	"strings"
)

func Load(name string) string {
	return "LOAD " + name + " "
}

func Choose(cond string) string {
	return " CHOOSE " + cond
}

func leftRight(keyword string, length int, left, right func(i int) string) string {
	var sb strings.Builder
	sb.WriteString(keyword)
	sb.WriteString(" ")
	for i := 0; i < length; i++ {
		sb.WriteString("(")
		sb.WriteString(left(i))
		sb.WriteString(":")
		sb.WriteString(right(i))
		sb.WriteString(")")
		if i != length-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

func AggFunc(length int, left, right func(i int) string) string {
	return leftRight(" AGG", length, left, right)
}

func Agg(left, right string) string {
	return " AGG " + Paren(Quote(left)+":"+right)
}

func As(left, right string) string {
	return " AS " + Paren(Quote(left)+":"+right)
}

func AsFunc(length int, left, right func(i int) string) string {
	return leftRight(" AS", length, left, right)
}

func Fun(fun string, arg ...string) string {
	return fun + "(" + strings.Join(arg, ",") + ")"
}

func Quote(path string) string {
	return "'" + path + "'"
}

func DoubleQuote(s string) string {
	return `"` + s + `"`
}

func Paren(s string) string {
	return "(" + s + ")"
}

func Store(s string) string {
	return " STORE " + s
}

func Delete(s string) string {
	return " DELETE " + s
}

func GroupedBy(dataset, jsonPoint, storeVar string) string {
	return Load(dataset) +
		Store(" GROUPED BY "+Quote(jsonPoint)) +
		"\n" +
		Load(storeVar) + " FROM GROUPED " + Quote(jsonPoint)
}

func Group(left, as, right string) string {
	return " GROUP " + left + " AS " + as + " BY " + right
}
