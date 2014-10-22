package main

import (
	//"encoding/json"
	"strconv"
	"strings"
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
	//logDebug("MarshalJSON:" + msg)
	return []byte(msg), nil
}

func (source *Money) UnmarshalJSON(b []byte) error {

	msg := string(b[:])

	if len(msg) > 1 {
		money := strings.Split(msg, ".")

		if len(money) == 2 {
			source.Number, _ = ToInt64(money[0])
			source.Points, _ = ToInt64(money[1])
		}

	}

	return nil
}
