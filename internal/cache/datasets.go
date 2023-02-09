package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type UsedDataset struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	SourceType int    `json:"type"` // 0 = FROM FILES, 1 = FROM FILE, 2 = FROM URL
	Frequency  int    `json:"frequency"`
}

func (thisdataset UsedDataset) IsSame(other UsedDataset) bool {
	return thisdataset.Name == other.Name && thisdataset.Path == other.Path && thisdataset.SourceType == other.SourceType
}

func (dataset UsedDataset) ToQuery() string {
	typestring := ""
	switch dataset.SourceType {
	case 0:
		typestring = "FILES"
	case 1:
		typestring = "FILE"
	case 2:
		typestring = "URL"
	default:
		log.Panicf("Unknown Source type: %d", dataset.SourceType)
	}

	return fmt.Sprintf("LOAD %s FROM %s \"%s\" ", dataset.Name, typestring, dataset.Path)
}

func AddUsedDataset(ds UsedDataset) {
	datasets, err := GetFrequentDatasets()
	if err != nil {
		log.Printf("Could not get frequent datasets from cache: %v", err)
		return
	}
	for i, dataset := range datasets {
		if ds.IsSame(dataset) {
			datasets[i].Frequency++
			err := storeFrequentDatasets(datasets)
			if err != nil {
				log.Printf("Could not store frequent datasets: %v", err)
				return
			}
			return
		}
	}
	datasets = append(datasets, ds)
	err = storeFrequentDatasets(datasets)
	if err != nil {
		log.Printf("Could not store frequent datasets: %v", err)
		return
	}
}

func getDataSetFile() string {
	cache := GetCacheDir()
	if len(cache) == 0 {
		return ""
	}
	return cache + "/frequentsets.json"
}

func storeFrequentDatasets(datasets []UsedDataset) error {
	fileName := getDataSetFile()
	if len(fileName) == 0 {
		return fmt.Errorf("could net get dataset file name (no cache/temp dir found")
	}
	content, err := json.MarshalIndent(datasets, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, content, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func GetFrequentDatasets() ([]UsedDataset, error) {
	fileName := getDataSetFile()
	if len(fileName) == 0 {
		return nil, fmt.Errorf("could net get dataset file name (no cache/temp dir found")
	}

	datasets := []UsedDataset{}

	// Ensure file exists and initialize
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			content, err := json.MarshalIndent(datasets, "", " ")
			if err != nil {
				return nil, err
			}
			err = ioutil.WriteFile(fileName, content, os.ModePerm)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	//Read file
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &datasets)
	if err != nil {
		return nil, err
	}
	return datasets, nil
}

func ExtractUsedDataset(query string) *UsedDataset {
	reg := regexp.MustCompile("(?i)LOAD\\s+(.*)\\s+FROM\\s+(FILE|FILES|URL)\\s+\"(.*)\"")

	matches := reg.FindStringSubmatch(query)
	if len(matches) == 4 {
		t := 0
		switch source := strings.ToLower(matches[2]); source {
		case "files":
			t = 0
		case "file":
			t = 1
		case "url":
			t = 2
		default:
			t = 99
		}
		if t == 99 {
			log.Panicln("Unsupported source type")
		}
		return &UsedDataset{
			Name:       matches[1],
			Path:       matches[3],
			SourceType: t,
			Frequency:  1,
		}
	}

	return nil
}
