package cond

import (
	"html/template"
)

type typeCond struct {
	jsonPoint string
	neg       bool
	condType  Type
}

func newTypeCond(jsonPoint string, neg bool, condType Type) typeCond {
	return typeCond{jsonPoint: jsonPoint, neg: neg, condType: condType}
}

func (tc typeCond) JsonPoint() string {
	return tc.jsonPoint
}

func (tc typeCond) string() string {
	return "IS" + tc.CondType().String() + "('" + tc.jsonPoint + "')"
}

func (tc typeCond) IsNeg() bool {
	return tc.neg
}

func (tc typeCond) CondType() Type {
	return tc.condType
}

func (tc typeCond) Label() string {
	if tc.neg {
		return "Is Not " + tc.CondType().String()
	} else {
		return "Is " + tc.CondType().String()
	}
}

func (tc typeCond) Desc() template.HTML {
	return template.HTML(tc.jsonPoint)
}

func (tc typeCond) Negate() Condition {
	tc.neg = !tc.neg
	return tc
}
