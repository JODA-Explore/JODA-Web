package query

import (
	"fmt"
	"strings"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/cmd"
)

const StoreVarPrefix = "__AUTO_GENERATED__"

type storeVar struct {
	name, parent, query string
	ins                 *Ins
}

func (sv storeVar) init() error {
	return sv.ins.Store(sv.query, sv.name)
}

func (sv storeVar) clear() {
	sv.ins.deleteStoreVar(sv.name)
}

func (ins Ins) newStoreVar(newName, parent, query string) (storeVar, error) {
	var name string
	if !strings.HasPrefix(parent, StoreVarPrefix) {
		name = StoreVarPrefix
	}
	name += parent + newName
	sv := storeVar{
		name:   name,
		parent: parent,
		query:  cmd.Load(parent) + query,
		ins:    &ins,
	}
	err := sv.init()
	if err != nil {
		err = fmt.Errorf("try to store var %v: %w", sv.name, err)
	}
	return sv, err
}

func (ins Ins) deleteStoreVar(storeVar string) {
	ins.Execute(cmd.Load(storeVar) + cmd.Delete(storeVar))
}

func (ins Ins) Store(query, storeVar string) error {
	ins.deleteStoreVar(storeVar)
	return ins.Execute(query + cmd.Store(storeVar))
}
