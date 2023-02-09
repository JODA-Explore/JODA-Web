package server

import (
	"net/http"
)

func systemHandler(w http.ResponseWriter, r *http.Request) {
	data := createData("System", "server")
	templateFile := "system.html"

	system, err := JodaInstance.GetSystem()
	if err != nil {
		addError(data, err)
		executeTemplate(w, templateFile, data)
		return
	}

	data["system"] = system.Summary()

	executeTemplate(w, templateFile, data)
}
