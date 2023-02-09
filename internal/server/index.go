package server

import (
	"log"
	"net/http"
	"sort"

	"github.com/JODA-Explore/JODA-Web/internal/cache"
	"github.com/JODA-Explore/JODA-Web/internal/joda"
	"github.com/JODA-Explore/JODA-Web/internal/source"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexHandlerWithMessages(w, r, messages{})
}

func indexHandlerWithMessages(w http.ResponseWriter, r *http.Request, messages messages) {
	data := createData("Overview", "list")
	data["messages"] = messages
	templateFile := "index.html"
	var sources, results []source.Source
	var modules []joda.ModuleSummary

	sources, err := JodaInstance.GetSources()
	if err != nil {
		log.Println(err.Error())
		addError(data, err)
	} else {
		results, err = JodaInstance.GetResults()
		if err != nil {
			log.Println(err.Error())
			addError(data, err)
		}

		modules, err = JodaInstance.GetModules()
		if err != nil {
			log.Println(err.Error())
			addError(data, err)
		}
	}

	data["sources"] = sources
	data["results"] = results
	data["modules"] = modules
	data["datasets"] = getFrequentDatasets(results, data)

	system, err := JodaInstance.GetSystem()
	if err != nil {
		executeTemplate(w, templateFile, data)
		return
	}

	data["system_summary"] = system.Summary()

	executeTemplate(w, templateFile, data)
}

type InterfaceDataset struct {
	Name  string
	Path  string
	Query string
}

func getFrequentDatasets(results []source.Source, data map[string]interface{}) (interfaceSet []InterfaceDataset) {
	datasets, err := cache.GetFrequentDatasets()
	if err != nil {
		log.Println(err.Error())
		addWarnMessage(data, err)
	}

	for _, source := range results {
		dataset := cache.ExtractUsedDataset(source.Query)
		if dataset != nil {
			datasets = removeDataset(datasets, *dataset)
		}
	}

	sort.Slice(datasets, func(i, j int) bool {
		return datasets[i].Frequency > datasets[j].Frequency
	})

	for _, dataset := range datasets {
		interfaceSet = append(interfaceSet, InterfaceDataset{
			Name:  dataset.Name,
			Path:  dataset.Path,
			Query: dataset.ToQuery(),
		})
		if len(interfaceSet) >= 9 {
			return
		}
	}

	return
}

func removeDataset(s []cache.UsedDataset, ds cache.UsedDataset) []cache.UsedDataset {
	i := -1
	for j, d := range s {
		if d.IsSame(ds) {
			i = j
			break
		}
	}
	if i < 0 {
		return s
	}
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
