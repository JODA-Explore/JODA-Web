package server

import (
	"net/http"

	"github.com/JODA-Explore/JODA-Web/internal/source"
)

func foundID(foundMap map[string]bool, id string, source source.Source) {
	if source.Name == id {
		foundMap[id] = true
	}
}

func icdeDemoHandler(w http.ResponseWriter, r *http.Request) {
	data := createData("Example Queries", "laptop")
	templateFile := "demo-queries.html"

	disabledQueries := map[string]bool{}

	sources, err := JodaInstance.GetSources()
	if err != nil {
		addError(data, err)
		disabledQueries["all"] = true
	} else {
		found := map[string]bool{}
		for _, v := range sources {
			//Movies
			foundID(found, "movies", v)
			foundID(found, "actors", v)
		}

		if _, ok := found["movies"]; !ok {
			disabledQueries["qM_2"] = true
			disabledQueries["qM_3"] = true
			disabledQueries["qM_4"] = true
		} else {
			disabledQueries["qM_1"] = true
		}

		if _, ok := found["actors"]; !ok {
			disabledQueries["qM_5"] = true
		} else {
			disabledQueries["qM_4"] = true
		}
	}

	data["disabledQueries"] = disabledQueries
	executeTemplate(w, templateFile, data)
}
