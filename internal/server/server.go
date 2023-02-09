package server

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/query"
	"github.com/JODA-Explore/JODA-Web/internal/joda"
)

var (
	JodaInstance joda.Joda
	templates    *template.Template
)

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	data := createData("Page not Found", "exclamation-triangle")
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		log.Println("404: Unknown Page requested: \"" + r.URL.EscapedPath() + "\"")
		templates.ExecuteTemplate(w, "404.html", data)
	}
}

func createData(title string, icon string) map[string]interface{} {
	return map[string]interface{}{
		"title":    title,
		"icon":     icon,
		"messages": messages{},
	}
}

func executeTemplate(w http.ResponseWriter, temp string, data map[string]interface{}) {
	// get a buffer
	buf := new(bytes.Buffer)
	// generate the template into a buffer
	err := templates.ExecuteTemplate(buf, temp, data)
	// when there is an error, discard the buffer and display the error page
	if err != nil {
		log.Println("Can’t load template: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = templates.ExecuteTemplate(w, "error", nil)
		return
	}
	// if all is good then write buffer to the response writer
	w.Header().Set("Content-Type", " text/html; charset=UTF-8")
	buf.WriteTo(w)
}

func getAllStringParameters(w http.ResponseWriter, r *http.Request) map[string][]string {
	params := make(map[string][]string)
	switch r.Method {
	case "GET":
		for k, v := range r.URL.Query() {
			params[k] = v
		}
	case "POST":
		err := r.ParseForm()
		if err != nil {
			log.Println("Post request, could not parse form!")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		} else {
			for k, v := range r.Form {
				params[k] = v
			}
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
	return params
}

func jodaHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/joda/" {
		indexHandler(w, r)
	} else if r.URL.Path == "/joda/delete" {
		deleteHandler(w, r)
	} else if r.URL.Path == "/joda/query" {
		queryHandler(w, r)
	} else if r.URL.Path == "/joda/queryExec" {
		executeQueryHandler(w, r)
	} else if r.URL.Path == "/joda/result" {
		resultHandler(w, r)
	} else if r.URL.Path == "/joda/system" {
		systemHandler(w, r)
	} else if r.URL.Path == "/joda/analyze" {
		analyzeHandler(w, r)
	} else if r.URL.Path == "/joda/download" {
		downloadHandler(w, r)
	} else if r.URL.Path == "/joda/betze" {
		explorerHandler(w, r)
	} else if r.URL.Path == "/joda/demo" {
		icdeDemoHandler(w, r)
	} else if r.URL.Path == "/joda/module" {
		uploadModuleHandler(w, r)
	} else if r.URL.Path == "/config" {
		configHandler(w, r)
	} else if r.URL.Path == "" {
		if JodaInstance.TestConnect() != nil {
			http.Redirect(w, r, "/config", http.StatusFound)
		} else {
			http.Redirect(w, r, "/joda", http.StatusFound)
		}
	} else {
		errorHandler(w, r, http.StatusNotFound)
	}
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/config" {
		configHandler(w, r)
	} else if r.URL.Path == "/favicon.ico" {
		faviconHandler(w, r)
	} else if r.URL.Path == "/" {
		if JodaInstance.TestConnect() != nil {
			http.Redirect(w, r, "/config", http.StatusFound)
		} else {
			http.Redirect(w, r, "/joda", http.StatusFound)
		}
	} else {
		errorHandler(w, r, http.StatusNotFound)
	}
}

func nthChild(w http.ResponseWriter, r *http.Request, idx int) (string, bool) {
	splited := strings.Split(r.URL.Path, "/")
	if len(splited) < idx+1 {
		errorHandler(w, r, http.StatusNotFound)
		return "", false
	}
	return splited[idx], true
}

func api(w http.ResponseWriter, r *http.Request) error {
	point, ok := nthChild(w, r, 2)
	if !ok {
		errorHandler(w, r, http.StatusNotFound)
		return nil
	}
	switch point {
	case "track_rating":
		return trackRating(w, r)
	case "track_guess_rating":
		return trackGuessRating(w, r)
	case "distinctValues":
		return distinctValues(w, r)
	case "memberFreq":
		return memberFreq(w, r)
	case "allDistinctMembers":
		return allDistinctMembers(w, r)
	case "allDistinctMembersCount":
		return allDistinctMembersCount(w, r)
	case "queries":
		return exploreContent(w, r)
	case "QueryGenerator":
		return queryMaker(w, r)
	case "sources":
		return sources(w, r)
	case "filteredSources":
		return filteredSources(w, r)
	case "history":
		return history(w, r)
	case "datasetDesc":
		return datasetDesc(w, r)
	case "searchOpt":
		return searchOpt(w, r)
	case "newChild":
		return newChild(w, r)
	default:
		errorHandler(w, r, http.StatusNotFound)
	}
	return nil
}

func Start(host string, port uint) {
	// Register static file server
	http.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("../../web/static/"))),
	)

	// Register start page
	http.HandleFunc("/joda/explore", handler(exploreSearch))
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/joda/", jodaHandler)
	http.HandleFunc("/api/", handler(api))
	var err error
	templates = template.New("template")
	err = filepath.Walk(
		"../../web/template/",
		func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".html") {
				templates.ParseFiles(path)
			}
			return nil
		},
	)
	if err != nil {
		log.Println("Cannot parse templates:", err)
		os.Exit(-1)
	}

	if host != "" {
		log.Println("Using host: " + host)
		JodaInstance = joda.New(host)
	} else {
		JodaInstance = joda.New("http://localhost:5632")
	}
	ins = query.Ins{Joda: &JodaInstance}
	portString := ":" + strconv.FormatUint(uint64(port), 10)
	log.Println("Starting JODA-Web on Port " + portString)
	log.Fatal(http.ListenAndServe(portString, nil))
}
