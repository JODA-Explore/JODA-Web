package joda

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/valyala/fastjson"
)

type Joda struct {
	host string
}

const api = "/api/v2"

func New(host string) Joda {
	j := Joda{
		host: host,
	}
	return j
}

func (j Joda) TestConnect() error {
	_, err := j.GetSystem()
	if err != nil {
		return err
	}
	return nil
}

func (j Joda) GetHost() string {
	return j.host
}

func getJSON(url string, target interface{}) error {
	// log.Println("Calling endpoint: " + url)
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func getFastJSON(url string) (v *fastjson.Value, err error) {
	r, err := http.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()
	sb := new(strings.Builder)
	_, err = io.Copy(sb, r.Body)
	if err != nil {
		return
	}
	var p fastjson.Parser
	return p.Parse(sb.String())
}
