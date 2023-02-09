package joda

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/JODA-Explore/JODA-Web/internal/source"
)

func (j Joda) RemoveSource(s source.Source) error {
	params := url.Values{}
	if s.Name != "" {
		params.Add("name", s.Name)
	} else if s.ID != 0 {
		params.Add("result", strconv.FormatUint(s.ID, 10))
	} else {
		log.Fatal("Source has neither name nor ID")
	}
	resp, err := http.Get(j.host + api + "/delete?" + params.Encode())
	if err != nil {
		log.Println("Could not get delete result. HTTP Error: " + err.Error())
		return err
	}
	defer resp.Body.Close()

	return nil
}
