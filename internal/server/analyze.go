package server

import (
	"log"
	"net/http"

	"github.com/JODA-Explore/JODA-Web/internal/source"
	"github.com/mitchellh/mapstructure"
)

// {"Count_Total":89473,"Count_Object":89473,"Count_Array":0,"Count_Null":0,"Count_Boolean":0,"Count_String":0,"Count_Number":0,"Children":[{"Key":"created_at",
type attribute struct {
	Key          string      `json:"Key"`
	CountTotal   uint64      `json:"Count_Total"`
	CountObject  uint64      `json:"Count_Object"`
	CountArray   uint64      `json:"Count_Array"`
	CountNull    uint64      `json:"Count_Null"`
	CountNumber  uint64      `json:"Count_Number"`
	CountBoolean uint64      `json:"Count_Boolean"`
	CountString  uint64      `json:"Count_String"`
	Children     []attribute `json:"Children"`
}

type d3TContent struct {
	CountTotal   uint64 `json:"Count_Total"`
	CountObject  uint64 `json:"Count_Object"`
	CountArray   uint64 `json:"Count_Array"`
	CountNull    uint64 `json:"Count_Null"`
	CountNumber  uint64 `json:"Count_Number"`
	CountBoolean uint64 `json:"Count_Boolean"`
	CountString  uint64 `json:"Count_String"`
}
type d3Tree struct {
	Name     string     `json:"name"`
	Content  d3TContent `json:"content"`
	Children []d3Tree   `json:"children"`
}

func (a attribute) ToD3() d3Tree {
	ret := d3Tree{
		Name: a.Key,
		Content: d3TContent{
			CountTotal:   a.CountTotal,
			CountObject:  a.CountObject,
			CountArray:   a.CountArray,
			CountNull:    a.CountNull,
			CountNumber:  a.CountNumber,
			CountBoolean: a.CountBoolean,
			CountString:  a.CountString,
		},
	}

	for _, v := range a.Children {
		ret.Children = append(ret.Children, v.ToD3())
	}

	return ret
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	data := createData("Analyze Source", "flask")
	template := "analyze.html"

	params := getAllStringParameters(w, r)
	if sources, ok := params["source"]; ok && len(sources) == 1 {
		sourceName := sources[0]
		data["title"] = "Analyze Source \"" + sourceName + "\""

		result, _, err := JodaInstance.ExecuteQuery("LOAD " + sourceName + " AGG ('':ATTSTAT(''))")
		if err != nil {
			message := "Could not analyze source: " + err.Error()
			log.Println(message)
			addErrorMessage(data, message)
			executeTemplate(w, template, data)
			return
		}
		//Remove result again
		defer JodaInstance.RemoveSource(source.Source{ID: result})

		resultDocuments, err := JodaInstance.GetResultDocuments(result, 0, 1)
		if err != nil {
			message := "Could not analyze source: " + err.Error()
			log.Println(message)
			addErrorMessage(data, message)
			executeTemplate(w, template, data)
			return
		}
		if len(resultDocuments) != 1 {
			message := "Expected one result document"
			log.Println(message)
			addErrorMessage(data, message)
			executeTemplate(w, template, data)
			return
		}

		analyze := attribute{}
		config := &mapstructure.DecoderConfig{TagName: "json", Result: &analyze}
		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			message := "Could not decode analyze result: " + err.Error()
			log.Println(message)
			addErrorMessage(data, message)
			executeTemplate(w, template, data)
			return
		}
		decoder.Decode(resultDocuments[0])
		tree := analyze.ToD3()

		data["tree"] = tree
		executeTemplate(w, template, data)

	} else {
		message := "Missing/Wrong parameter 'source'. Expect exactly one source"
		addErrorMessage(data, message)
		executeTemplate(w, template, data)
		return
	}

}
