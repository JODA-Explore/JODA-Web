package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func downloadResultset(w http.ResponseWriter, r *http.Request, result uint64, linesep bool) {
	bulkSize := uint64(100)
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Fatalln("expected http.ResponseWriter to be an http.Flusher")
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Disposition", "attachment; filename=result.json")
	w.Header().Set("Content-Type", "application/json")
	if !linesep {
		w.Write([]byte("[\n"))
	}
	offset := uint64(0)
	var err error
	for {
		array, err := JodaInstance.GetResultDocuments(result, offset, bulkSize)
		if len(array) == 0 || err != nil {
			break
		}
		for _, elt := range array {
			if offset != 0 && !linesep {
				w.Write([]byte(",\n"))
			}
			offset++
			textjson, err := json.Marshal(elt)
			if err != nil {
				log.Fatalln(err.Error())
				return
			}
			if !linesep {
				w.Write([]byte("  "))
			}
			w.Write(textjson)
			if linesep {
				w.Write([]byte("\n"))
			}
		}
		flusher.Flush()
	}
	if err != nil {
		errMessage := "Could not get results: " + err.Error()
		log.Println(errMessage)
		return
	}
	if !linesep {
		w.Write([]byte("\n]"))
	}
}

func downloadDocument(w http.ResponseWriter, r *http.Request, result uint64, doc uint64) {
	array, err := JodaInstance.GetResultDocuments(result, doc, 1)
	if err != nil {
		errMessage := "Could not get results: " + err.Error()
		log.Println(errMessage)
		return
	}
	if len(array) == 0 {
		log.Println("Result is empty")
		return
	}
	textjson, err := json.MarshalIndent(array[0], "", "  ")
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=result-document.json")
	w.Header().Set("Content-Type", "application/json")
	w.Write(textjson)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
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
			downloadDocument(w, r, result, docnum)
			return
		} else { //Download whole result set
			linesep := false
			if _, ok := params["lineseparated"]; ok {
				linesep = true
			}
			downloadResultset(w, r, result, linesep)
		}
	} else {
		//TODO Error handling
		log.Println("Missing 'results' parameter")
	}

}
