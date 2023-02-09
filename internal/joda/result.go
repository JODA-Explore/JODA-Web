package joda

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/valyala/fastjson"
)

type resultResponse struct {
	ErrorMessage string        `json:"error"`
	Result       []interface{} `json:"result"`
}

func (j Joda) GetResultDocuments(result uint64, offset uint64, count uint64) ([]interface{}, error) {
	params := url.Values{}
	params.Add("id", strconv.FormatUint(result, 10))
	params.Add("offset", strconv.FormatUint(offset, 10))
	params.Add("count", strconv.FormatUint(count, 10))
	response := new(resultResponse)
	err := getJSON(j.host+api+"/result?"+params.Encode(), &response)
	if err != nil {
		return nil, err
	}

	if response.ErrorMessage != "" {
		return nil, errors.New(response.ErrorMessage)
	}
	// fmt.Printf("Result:%v", response.Result)
	return response.Result, nil
}

func (j Joda) GetResultDocumentsFastJson(result uint64, offset uint64, count uint64) (v *fastjson.Value, err error) {
	params := url.Values{}
	params.Add("id", strconv.FormatUint(result, 10))
	params.Add("offset", strconv.FormatUint(offset, 10))
	params.Add("count", strconv.FormatUint(count, 10))
	str := j.host + api + "/result?" + params.Encode()
	_ = str
	v, err = getFastJSON(j.host + api + "/result?" + params.Encode())
	if err != nil {
		return
	}
	if v.Exists("error") {
		err = errors.New(string(v.GetStringBytes("error")))
		return
	}
	v = v.Get("result")
	return
}
