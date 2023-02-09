package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/JODA-Explore/JODA-Web/internal/source"
)

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	data := createData("Delete", "trash")

	params := getAllStringParameters(w, r)
	if names, ok := params["name"]; ok {
		for _, name := range names {
			err := JodaInstance.RemoveSource(source.Source{Name: name})
			if err != nil {
				addErrorMessage(data, err)
			} else {
				addSuccessMessage(data, "Removed Collection '"+name+"'")
			}

		}
	}
	if results, ok := params["result"]; ok {
		for _, result := range results {
			numResult, err := strconv.ParseUint(result, 10, 64)
			if err == nil {
				err := JodaInstance.RemoveSource(source.Source{ID: numResult})
				if err != nil {
					addErrorMessage(data, err)
				} else {
					addSuccessMessage(data, "Removed Result '"+result+"'")
				}
			} else {
				log.Println("Could not parse int: " + result + ". ")
				addErrorMessage(data, "Could not parse int: "+result+". ")
			}
		}
	}

	indexHandlerWithMessages(w, r, data["messages"].(messages))

}
