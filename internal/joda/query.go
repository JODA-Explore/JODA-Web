package joda

import (
	"errors"
	"net/url"
)

type runtime struct {
	Query        float64 `json:"Query"`
	Optimization float64 `json:"Optimization"`
	Planning     float64 `json:"Planning"`
}

// type pipelineTask struct {
// 	Num              uint64   `json:"Num"`
// 	Name             string   `json:"Name"`
// 	Async            string   `json:"Async"`
// 	HumanRuntimeTime *string  `json:"Runtime"`
// 	Runtime          *float64 `json:"Runtime_s"`
// 	TaskCount        *uint64  `json:"TaskCount"`
// }

// type pipelineConnection struct {
// 	From       []uint64 `json:"From"`
// 	To         []uint64 `json:"To"`
// 	Throughput uint64   `json:"Throughput"`
// 	Finished   bool     `json:"Finished"`
// }

// type pipeline struct {
// 	MaxThreads  uint64               `json:"MaxThreads"`
// 	Tasks       []pipelineTask       `json:"Tasks"`
// 	Connections []pipelineConnection `json:"Connections"`
// }

type benchmark struct {
	Query               string       `json:"Query"`
	Time                uint64       `json:"Time"`
	HumanTime           string       `json:"Pretty Time"`
	Threads             uint64       `json:"Threads"`
	ResultSize          int64        `json:"Result Size"`
	Container           uint64       `json:"#Container"`
	Runtime             *runtime     `json:"Runtime"`
	Pipeline            interface{}  `json:"Pipeline"`
	UnoptimizedPipeline *interface{} `json:"UnoptimizedPipeline"`
}

type queryResponse struct {
	ErrorMessage string     `json:"error"`
	Success      uint64     `json:"success"`
	Size         int64      `json:"size"`
	Message      string     `json:"message"`
	Benchmark    *benchmark `json:"benchmark"`
}

var EmptyResError = errors.New("The result is empty")

func (j Joda) ExecuteQuery(query string) (uint64, *benchmark, error) {
	// log.Println("Executing query: \"" + query + "\"")
	params := url.Values{}
	params.Add("query", query)
	response := new(queryResponse)
	err := getJSON(j.host+api+"/query?"+params.Encode(), &response)
	if err != nil {
		return 0, nil, err
	}

	if response.ErrorMessage != "" {
		return 0, nil, errors.New(response.ErrorMessage)
	}

	if response.Size <= 0 {
		return 0, nil, EmptyResError
	}

	// fmt.Printf("%#v", response)
	// fmt.Printf("%#v", response.Benchmark)

	return response.Success, response.Benchmark, nil
}
