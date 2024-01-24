package xp

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/neuro-team-femto/revcor/helpers"
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

var introHeaders = []string{"subj", "trial", "block", "date", "stim", "stim_order", "response", "rt"}

func genRecordHeaders(es ExperimentSettings, p Participant, r Result) (headers []string, err error) {
	paramHeaders, err := getAssetDefHeaders(es, r.Stimulus)
	if err != nil {
		return
	}
	headers = append(headers, introHeaders...)
	headers = append(headers, p.getInfoKeys()...)
	headers = append(headers, "param_index")
	headers = append(headers, paramHeaders...)
	return
}

func newRecord(p Participant, r Result, index int, values []string) []string {
	record := []string{
		p.Id,
		r.Trial,
		r.Block,
		r.Date,
		r.Stimulus,
		r.Order,
		r.Response,
		r.Rt,
	}
	record = append(record, p.getInfoValues()...)
	record = append(record, fmt.Sprint(index))
	record = append(record, values...)
	return record
}

// API

func WriteToCSV(es ExperimentSettings, p Participant, r Result) (err error) {
	def, err := getAssetDefAllValues(es, r.Stimulus)
	if err != nil {
		return
	}

	var records [][]string
	path := "data/" + p.ExperimentId + "/results/" + p.Id + ".csv"
	if !helpers.PathExists(path) {
		headers, e := genRecordHeaders(es, p, r)
		if e != nil {
			return e
		}
		records = append(records, headers)
	}

	for index, values := range def {
		records = append(records, newRecord(p, r, index, values))
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	w := csv.NewWriter(file)
	w.WriteAll(records)
	return
}
