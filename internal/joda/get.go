package joda

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/JODA-Explore/JODA-Web/internal/source"
)

func (j Joda) GetSources() ([]source.Source, error) {
	endpoint := j.host + api + "/sources"
	log.Println("Getting Sources (" + endpoint + ")")
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Println("Could not get sources. HTTP Error: " + err.Error())
		return nil, errors.New("Could not get sources. HTTP Error: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not read sources. Parse Error: " + err.Error())
		return nil, errors.New("Could not read sources. Parse Error: " + err.Error())
	}

	sources := make([]source.Source, 0)
	err = json.Unmarshal(body, &sources)
	if err != nil {
		log.Println("Could not parse sources. Parse Error: " + err.Error())
		return nil, errors.New("Could not parse sources. Parse Error: " + err.Error())
	}
	return sources, nil
}

func (j Joda) GetResults() ([]source.Source, error) {
	endpoint := j.host + api + "/results"
	log.Println("Getting Results (" + endpoint + ")")
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Println("Could not get sources. HTTP Error: " + err.Error())
		return nil, errors.New("Could not get sources. HTTP Error: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not read sources. Parse Error: " + err.Error())
		return nil, errors.New("Could not read sources. Parse Error: " + err.Error())
	}

	sources := make([]source.Source, 0)
	err = json.Unmarshal(body, &sources)
	if err != nil {
		log.Println("Could not parse sources. Parse Error: " + err.Error())
		return nil, errors.New("Could not parse sources. Parse Error: " + err.Error())
	}

	return sources, nil
}
