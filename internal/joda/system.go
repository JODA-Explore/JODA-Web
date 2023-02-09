package joda

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"github.com/dustin/go-humanize"
)

type Memory struct {
	Total   uint64 `json:"total"`
	Used    uint64 `json:"used"`
	Joda    uint64 `json:"joda"`
	Allowed uint64 `json:"allowed_memory"`
	Stored  uint64 `json:"calculated_memory"`
}

type Host struct {
	Kernel string `json:"kernel"`
	OS     string `json:"os"`
}

type Version struct {
	Version   string `json:"version"`
	API       uint   `json:"api"`
	Commit    string `json:"commit"`
	BuildTime string `json:"build-time"`
}
type System struct {
	Memory  Memory  `json:"memory"`
	Version Version `json:"version"`
	Host    Host    `json:"host"`
}

type MemorySummary struct {
	TotalHuman   string
	UsedHuman    string
	UsedFloat    float64
	UsedPerc     int
	UsedColor    string
	JodaHuman    string
	AllowedHuman string
	StoredHuman  string
	StoredPerc   int
	StoredFloat  float64
	StoredColor  string
	Min          func(int, int) int
}
type Summary struct {
	Memory MemorySummary
	Raw    System
}

func (j Joda) GetSystem() (*System, error) {
	log.Println("Getting System")
	resp, err := http.Get(j.host + api + "/system")
	if err != nil {
		log.Println("Could not get system info. HTTP Error: " + err.Error())
		return nil, errors.New("Could not get system info. HTTP Error: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not read system info. Parse Error: " + err.Error())
		return nil, errors.New("Could not read system info. Parse Error: " + err.Error())
	}

	system := new(System)
	json.Unmarshal(body, &system)
	return system, nil
}

func getStatusColor(f float64) string {
	ret := "green"
	if f > 0.7 {
		ret = "amber"
	}
	if f > 0.8 {
		ret = "deep-orange"
	}
	if f > 0.9 {
		ret = "red"
	}
	return ret
}

func (s System) Summary() Summary {
	used := float64(s.Memory.Used) / float64(s.Memory.Total)
	usedColor := getStatusColor(used)
	stored := float64(s.Memory.Stored) / float64(s.Memory.Allowed)
	storedColor := getStatusColor(stored)
	return Summary{
		Memory: MemorySummary{
			TotalHuman:   humanize.IBytes(s.Memory.Total),
			UsedHuman:    humanize.IBytes(s.Memory.Used),
			UsedFloat:    used,
			UsedPerc:     int(math.Round(used * 100)),
			UsedColor:    usedColor,
			JodaHuman:    humanize.IBytes(s.Memory.Joda),
			AllowedHuman: humanize.IBytes(s.Memory.Allowed),
			StoredHuman:  humanize.IBytes(s.Memory.Stored),
			StoredFloat:  stored,
			StoredPerc:   int(math.Round(stored * 100)),
			StoredColor:  storedColor,
			Min: func(i1 int, i2 int) int {
				return int(math.Min(float64(i1), float64(i2)))
			},
		},
		Raw: s,
	}
}
