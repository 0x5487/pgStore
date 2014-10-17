package main

import (
	"encoding/json"
	"strconv"
)

type Money struct {
	Number int64
	Points int64
}

/*
func (source Money) New(value string) *Money{

}
*/

func (source *Money) MarshalJSON() ([]byte, error) {
	msg := strconv.FormatInt(source.Number, 10) + "." + PadRight(strconv.FormatInt(source.Points, 10), "0", 4)
	return []byte(msg), nil
}

func (source *Money) UnmarshalJSON(b []byte) error {
	logDebug("unmarshall money")
	var msg string
	json.Unmarshal(b, &msg)
	source.Number = 123
	source.Points = 456
	return nil
}
