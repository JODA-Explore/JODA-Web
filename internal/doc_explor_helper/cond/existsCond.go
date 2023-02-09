package cond

import (
	"html/template"
)

type existsCond struct {
	jsonPoint string
	neg       bool
}

func newExistsCond(jsonPoint string, neg bool) existsCond {
	return existsCond{jsonPoint: jsonPoint, neg: neg}
}

func (ec existsCond) JsonPoint() string {
	return ec.jsonPoint
}

func (ec existsCond) string() string {
	return "EXISTS('" + ec.jsonPoint + "')"
}

func (ec existsCond) IsNeg() bool {
	return ec.neg
}

func (ec existsCond) CondType() Type {
	return Exists
}

func (ec existsCond) Label() string {
	if ec.neg {
		return "Non Exists"
	} else {
		return "Exists"
	}
}

func (ec existsCond) Desc() template.HTML {
	return template.HTML(ec.jsonPoint)
}

func (ec existsCond) Negate() Condition {
	ec.neg = !ec.neg
	return ec
}
