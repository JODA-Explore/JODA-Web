package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/explore"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/query"
)

func handler(fn func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			log.Println(err)
			data := createData("Explorer", "compass")
			addError(data, err)
			templateFile := "explore.html"

			executeTemplate(w, templateFile, data)
		}
	}
}

var ins query.Ins

func exploreContent(w http.ResponseWriter, r *http.Request) (err error) {
	data := createData("Explorer", "compass")
	templateFile := "explore_content.html"

	searchOpt := explore.SearchOpt{}
	err = r.ParseForm()
	if err != nil {
		return fmt.Errorf("parse form search option:%w", err)
	}
	for i := range searchOpt.EnabledBackends {
		id := backends.ID(i)
		name := id.Name()
		searchOpt.EnabledBackends[id] = r.Form.Get(name) == "on"
	}
	maxNumber := r.Form.Get("dehMaxItems")
	if maxNumber == "" {
		return errors.New("Please input max item in Filter!")
	}
	searchOpt.Filter.MaxNumber, err = strconv.Atoi(maxNumber)
	if err != nil {
		return fmt.Errorf("The max number is invalid:%w", err)
	}
	minRating := r.Form.Get("dehMinRating")
	if minRating == "" {
		return errors.New("Please input min rating in Filter!")
	}
	searchOpt.Filter.MinRating, err = strconv.ParseFloat(minRating, 64)
	if err != nil {
		return fmt.Errorf("The min rating is invalid:%w", err)
	}

	maxDiff := r.Form.Get("dehMaxDiff")
	if maxDiff == "" {
		return errors.New("Please input max diff in Filter!")
	}
	diff, err := strconv.ParseFloat(maxDiff, 64)
	if err != nil {
		return fmt.Errorf("The max diff is invalid:%w", err)
	}
	searchOpt.Filter.MaxDiff = diff / 100

	searchOpt.Dataset = r.Form.Get("name")
	if searchOpt.Dataset == "" {
		return errors.New("please input data source")
	}

	err = explore.HandleNewDataset(
		ins,
		r.Form.Get("sourceType"),
		r.Form.Get("parent"),
		r.Form.Get("query"),
		searchOpt,
	)
	if err != nil {
		return
	}

	resultsInfo, err := explore.Results(ins, searchOpt)
	if err != nil {
		return
	}

	data["dataset"] = searchOpt.Dataset
	data["results"] = resultsInfo.Results
	executeTemplate(w, templateFile, data)
	// fmt.Printf("total time took: %v\n", time.Since(startTime).Round(time.Second))
	return nil
}

func exploreSearch(w http.ResponseWriter, r *http.Request) (err error) {
	sourceNames, err := ins.FilteredDatasets()
	if err != nil {
		return
	}
	data := createData("Doc Exploration Helper", "compass")
	templateFile := "explore.html"
	data["analysisNames"] = backends.Names()
	data["sourceNames"] = sourceNames
	executeTemplate(w, templateFile, data)
	return nil
}

func distinctValues(w http.ResponseWriter, r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		return
	}
	number, err := strconv.Atoi(r.Form.Get("number"))
	if err != nil {
		panic(err)
	}
	topk, err := ins.DistinctValuesFreq(r.Form.Get("dataset"), r.Form.Get("path"), number)
	if err != nil {
		return
	}
	topk.SortReverse()
	l := topk.List()
	values := make([]string, len(l))
	counts := make([]uint64, len(l))
	for i, x := range topk.List() {
		values[i] = x.Val
		counts[i] = x.N
	}
	data := struct {
		Values []string
		Counts []uint64
	}{Values: values, Counts: counts}
	res, err := json.Marshal(data)
	if err != nil {
		return
	}
	w.Write(res)
	return
}

func memberFreq(w http.ResponseWriter, r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		return
	}
	number, err := strconv.Atoi(r.Form.Get("number"))
	if err != nil {
		panic(err)
	}
	f, err := ins.AtomicMemberFreq(r.Form.Get("dataset"), r.Form.Get("path"), number)
	if err != nil {
		return
	}
	f.SortReverse()
	list := f.List()
	values := make([]string, len(list))
	counts := make([]uint64, len(list))
	for i, x := range list {
		values[i] = x.Val
		counts[i] = x.N
	}
	data := struct {
		Values []string
		Counts []uint64
	}{Values: values, Counts: counts}
	res, err := json.Marshal(data)
	if err != nil {
		return
	}
	w.Write(res)
	return
}

func allDistinctMembers(w http.ResponseWriter, r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		return
	}
	members, err := ins.AllDistinctMembers(r.Form.Get("dataset"), r.Form.Get("path"))
	if err != nil {
		return
	}
	res, err := json.Marshal(members)
	if err != nil {
		return
	}
	w.Write(res)
	return
}

func allDistinctMembersCount(w http.ResponseWriter, r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		return fmt.Errorf("can not parse form for allDistinctMembersCount:%w", err)
	}
	members, err := ins.AllDistinctMembers(r.Form.Get("dataset"), r.Form.Get("path"))
	if err != nil {
		return
	}
	w.Write([]byte(strconv.Itoa(len(members))))
	return
}

func sources(w http.ResponseWriter, r *http.Request) (err error) {
	sources, err := ins.AllDatasets()
	if err != nil {
		return
	}
	j, err := json.Marshal(sources)
	if err != nil {
		return
	}
	w.Write(j)
	return
}

func filteredSources(w http.ResponseWriter, r *http.Request) (err error) {
	sources, err := ins.FilteredDatasets()
	if err != nil {
		return
	}
	j, err := json.Marshal(sources)
	if err != nil {
		return
	}
	w.Write(j)
	return
}

func trackRating(w http.ResponseWriter, r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		panic(err)
	}
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		panic(err)
	}
	resultsInfo := explore.CurrentResults()
	q := resultsInfo.Results[id]
	q.Rater.Desc()
	w.Write([]byte(q.Rater.Desc()))
	return
}

func trackGuessRating(w http.ResponseWriter, r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		panic(err)
	}
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		panic(err)
	}
	resultsInfo := explore.CurrentResults()
	q := resultsInfo.Results[id]
	guessRating := q.Info.(feature.EstimaterInfo).EstimatedRating()
	w.Write([]byte(guessRating.Desc()))
	return
}

func queryMaker(w http.ResponseWriter, r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		return
	}
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		return
	}
	qry, err := explore.CurrentQueryMaker(id)
	if err != nil {
		return
	}
	w.Write([]byte(qry))
	return
}

func history(w http.ResponseWriter, r *http.Request) (err error) {
	w.Write([]byte(explore.DescRoot()))
	return
}

func datasetDesc(w http.ResponseWriter, r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		return
	}
	w.Write([]byte(explore.DescDataset(r.Form.Get("dataset"))))
	return
}

func searchOpt(w http.ResponseWriter, r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		return
	}
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		return
	}
	res, err := explore.GetOpt(r.Form.Get("dataset"), id).ToJson()
	if err != nil {
		return
	}
	w.Write(res)
	return
}

func newChild(w http.ResponseWriter, r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		return
	}
	res, err := explore.NewChildName(r.Form.Get("dataset"))
	if err != nil {
		return
	}
	w.Write([]byte(res))
	return
}
