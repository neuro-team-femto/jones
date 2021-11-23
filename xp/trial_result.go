package xp

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/creamlab/revcor/helpers"
)

// result

type Result struct {
	Trial    string `json:"trial"`
	Block    string `json:"block"`
	Date     string `json:"date"`
	Stimulus string `json:"stimulus"`
	Order    string `json:"order"`
	Response string `json:"response"`
	Rt       string `json:"rt"`
}

// minimal implementation, we may check more
func (r Result) IsValid() bool {
	return len(r.Trial) > 0 &&
		len(r.Block) > 0 &&
		len(r.Date) > 0 &&
		len(r.Stimulus) > 0 &&
		len(r.Order) > 0 &&
		len(r.Response) > 0 &&
		len(r.Rt) > 0
}

// record formatting

var Headers = []string{"subj", "trial", "block", "sex", "age", "date", "stim", "stim_order", "filter", "filter_freq", "filter_gain", "response", "rt"}

func newRecord(p Participant, r Result, filterIndex int, f filter) []string {
	return []string{
		p.Id,
		r.Trial,
		r.Block,
		p.Sex,
		p.Age,
		r.Date,
		r.Stimulus,
		r.Order,
		fmt.Sprint(filterIndex),
		f.Freq,
		f.Gain,
		r.Response,
		r.Rt,
	}
}

// API

func WriteToCSV(p Participant, r1, r2 Result) (err error) {
	fs1, err := LoadFilters(p.XpId, r1.Stimulus)
	if err != nil {
		return
	}
	fs2, err := LoadFilters(p.XpId, r2.Stimulus)
	if err != nil {
		return
	}

	var records [][]string
	path := "data/" + p.XpId + "/results/" + p.Id + ".csv"
	if !helpers.PathExists(path) {
		records = append(records, Headers)
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	for index, f := range fs1 {
		records = append(records, newRecord(p, r1, index, f))
	}
	for index, f := range fs2 {
		records = append(records, newRecord(p, r2, index, f))
	}

	w := csv.NewWriter(file)
	w.WriteAll(records)
	return
}
