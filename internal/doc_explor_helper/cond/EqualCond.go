package cond

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
)

type equalCond struct {
	jsonPoint string
	right     string
	neg       bool
}

func newEqualCond(jsonPoint string, right string, neg bool) equalCond {
	return equalCond{jsonPoint: jsonPoint, right: right, neg: neg}
}

func (ec equalCond) JsonPoint() string {
	return ec.jsonPoint
}

func (ec equalCond) string() string {
	return "'" + ec.jsonPoint + "'" + " == " + ec.right
}

func (ec equalCond) IsNeg() bool {
	return ec.neg
}

func (ec equalCond) CondType() Type {
	return Equal
}

func (ec equalCond) Label() string {
	if ec.neg {
		return "Not Equal"
	} else {
		return "Equal"
	}
}

func (ec equalCond) Desc() template.HTML {
	left := desc.JsonPoint("'" + ec.jsonPoint + "'")
	var mid string
	if ec.IsNeg() {
		mid = ` ≠ `
	} else {
		mid = ` = `
	}
	return left + template.HTML(mid+ec.right)
}

func (ec equalCond) Negate() Condition {
	ec.neg = !ec.neg
	return ec
}
