package source

import (
	"html/template"
	"log"
	"net/url"
	"strconv"

	"github.com/dustin/go-humanize"
)

//Source represents a source in the JODA system
type Source struct {
	Name      string `json:"name"`
	ID        uint64 `json:"id"`
	Documents uint64 `json:"documents"`
	Container uint64 `json:"container"`
	Memory    uint64 `json:"memory"`
	Query     string `json:"query"`
}

//New creates a new Source from its fields
func New(name string, id uint64, documents uint64, container uint64, memory uint64) Source {
	s := Source{name, id, documents, container, memory, ""}
	return s
}

///////////////////////////////////////////////
///////////Display functions///////////////////
///////////////////////////////////////////////

//HumanDocuments returns the number of documents in the source as human readable String
func (s Source) HumanDocuments() string {
	return humanize.Comma(int64(s.Documents))
}

//HumanContainer returns the number of container in the source as human readable String
func (s Source) HumanContainer() string {
	return humanize.Comma(int64(s.Container))
}

//HumanMemory returns the memory size of the source as human readable String
func (s Source) HumanMemory() string {
	return humanize.IBytes(uint64(s.Memory))
}

//GetRemoveParams returns the URL params required to delete this source
func (s Source) GetRemoveParams() template.URL {
	params := url.Values{}
	if s.Name != "" {
		params.Add("name", s.Name)
	} else if s.ID != 0 {
		params.Add("result", strconv.FormatUint(s.ID, 10))
	} else {
		log.Fatal("Source has neither name nor ID")
	}
	return template.URL(params.Encode())
}
