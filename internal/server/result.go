package server

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/JODA-Explore/JODA-Web/internal/source"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func resultHandler(w http.ResponseWriter, r *http.Request) {
	params := getAllStringParameters(w, r)
	if results, ok := params["result"]; ok {
		result, err := strconv.ParseUint(results[0], 10, 64)
		if err != nil {
			log.Println("Could not parse int: " + results[0] + ". ")
		}

		if doc, ok := params["doc"]; ok {
			docnum, err := strconv.ParseUint(doc[0], 10, 64)
			if err != nil {
				log.Println("Could not parse int: " + doc[0] + ". ")
			}
			handleResult(w, r, result, docnum)
			return
		}
		handleResult(w, r, result, 0)
	} else {
		log.Println("Missing parameter result")
	}
}

func maybeGEOJSON(json interface{}) interface{} {
	v, ok := json.(map[string]interface{})
	if !ok {
		return nil
	}
	if _, found := v["type"]; !found {
		return nil
	}
	_, g := v["geometry"]
	_, c := v["coordinates"]
	featuresI, f := v["features"]
	if f {
		if features, ok := featuresI.([]interface{}); ok {
			v["features"] = features[:int(math.Min(float64(len(features)), 1000.0))]
		}
	}

	if g || c || f {
		return v
	}

	return nil
}

func syntaxHighlightResult(jsondoc interface{}) (string, error) {

	textjson, err := json.MarshalIndent(jsondoc, "", "  ")
	if err != nil {
		return "", err
	}
	maxChars := 100000
	if len(textjson) > maxChars { //Truncate long strings
		textjson = textjson[:maxChars]
	}

	lexer := lexers.Get("json")
	if lexer == nil {
		lexer = lexers.Fallback
	}
	style := styles.Get("tango")
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.Standalone(false))

	iterator, err := lexer.Tokenise(nil, string(textjson))
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = formatter.Format(buf, style, iterator)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func handleResult(w http.ResponseWriter, r *http.Request, result uint64, docnum uint64) {
	data := createData("Query Result", "file-alt")
	templateFile := "result.html"
	resultsets, err := JodaInstance.GetResults()
	if err != nil {
		errMessage := "Could not get resultset: " + err.Error()
		log.Println(errMessage)
		addErrorMessage(data, errMessage)
		executeTemplate(w, templateFile, data)
		return
	}
	resultset := source.Source{}
	for _, v := range resultsets {
		if v.ID == result {
			resultset = v
		}
	}

	array, err := JodaInstance.GetResultDocuments(result, docnum, 1)
	if err != nil {
		errMessage := "Could not get results: " + err.Error()
		log.Println(errMessage)
		addErrorMessage(data, errMessage)
		executeTemplate(w, templateFile, data)
		return
	}
	if len(array) == 0 {
		addWarnMessage(data, "Result is empty")
		executeTemplate(w, templateFile, data)
		return
	}
	prevDoc := docnum - 1
	if docnum == 0 {
		prevDoc = 0
	}

	nextDoc := docnum + 1
	if docnum == resultset.Documents-1 {
		nextDoc = resultset.Documents - 1
	}

	//Check geojson
	geojson := maybeGEOJSON(array[0])

	syntax, err := syntaxHighlightResult(array[0])

	if err != nil {
		errMessage := "Could not syntaxhighlight json " + err.Error()
		log.Println(errMessage)
		addErrorMessage(data, errMessage)
		executeTemplate(w, templateFile, data)
		return
	}

	data["interface"] = map[string]interface{}{
		"query":           resultset.Query,
		"result":          result,
		"doc_id":          docnum,
		"doc_id_next":     nextDoc,
		"doc_id_prev":     prevDoc,
		"doc_last":        resultset.Documents - 1,
		"doc_num":         docnum + 1,
		"doc_count":       resultset.Documents,
		"doc_count_human": resultset.HumanDocuments,
		"geojsonEnabled":  geojson != nil,
		"geojsondata":     geojson,
		"syntax":          template.HTML(syntax),
	}

	executeTemplate(w, templateFile, data)
}
