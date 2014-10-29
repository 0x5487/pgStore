package main

import (
	//"encoding/json"
	//"fmt"
	"errors"
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

		switch len(money) {
		case 1:
			number, err := ToInt64(msg)
			if err != nil {
				return errors.New("invalid format for money")
			}
			source.Number = number
			source.Points = 0
			break
		case 2:
			number, err := ToInt64(money[0])
			if err != nil {
				return errors.New("invalid format for money")
			}
			digital, err := ToInt64(money[1])
			if err != nil {
				return errors.New("invalid format for money")
			}

			source.Number = number
			source.Points = digital
			break
		}

	}

	return nil
}
