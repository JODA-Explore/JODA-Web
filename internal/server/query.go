package server

import (
	"net/http"
	"strings"
)

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	params := getAllStringParameters(w, r)
	query := "LOAD"
	errorMessage := ""
	if queries, ok := params["query"]; ok {
		query = queries[0]
	}
	if errors, ok := params["error"]; ok {
		errorMessage = errors[0]
	}
	handleQuery(w, r, query, errorMessage)
}

func handleQuery(w http.ResponseWriter, r *http.Request, query string, errorMessage string) {
	data := createData("Query", "search")
	templateFile := "query.html"
	if errorMessage != "" {
		addErrorMessage(data, errorMessage)
	}
	data["query"] = query
	data["editorLines"] = max(6, strings.Count(query, "\n")+1)
	executeTemplate(w, templateFile, data)
}
