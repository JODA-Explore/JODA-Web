package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/JODA-Explore/JODA-Web/internal/joda"
)

func configHandler(w http.ResponseWriter, r *http.Request) {
	data := createData("Configure", "cog")
	templateFile := "config.html"

	data["url"] = "http://localhost:5632"

	params := getAllStringParameters(w, r)

	if url, ok := params["url"]; ok && len(url) == 1 {
		data["url"] = url[0]
		tmpInstance := joda.New(url[0])
		err := tmpInstance.TestConnect()
		if err != nil {
			addError(data, err)
		} else {
			JodaInstance = tmpInstance
			addSuccessMessage(data, template.HTML("You can now use the webinterface <a href=\"/joda\" title=\"Webinterface\">here</a>."))
			log.Println("Using host: " + url[0])
		}
	}

	executeTemplate(w, templateFile, data)
}
