package server

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/JODA-Explore/JODA-Web/internal/cache"
)

func executeQueryHandler(w http.ResponseWriter, r *http.Request) {
	data := createData("Query Execution Statistics", "bar-chart")
	templateFile := "querystats.html"

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		addError(data, err)
		executeTemplate(w, templateFile, data)
		return
	}
	query := r.Form.Get("query")
	result, benchmark, err := JodaInstance.ExecuteQuery(query)
	if err != nil {
		log.Println("Error in query: " + err.Error())
		handleQuery(w, r, query, err.Error())
	} else {
		dataset := cache.ExtractUsedDataset(query)
		if dataset != nil {
			cache.AddUsedDataset(*dataset)
		}
		data["result"] = result
		if benchmark != nil { //If Statistics exist, show them
			data["benchmark"] = benchmark
			executeTemplate(w, templateFile, data)
		} else { //Else, just show result
			params := url.Values{}
			params.Add("result", strconv.FormatUint(result, 10))
			http.Redirect(w, r, "result?"+params.Encode(), http.StatusFound)
		}
	}
}
