package utils

import (
	"log"
	"strconv"
)

func StringToUint(s string) uint {
	valUint64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		log.Println("ERROR -> Couldn't convert string to uint")
	}

	return uint(valUint64)
}

func StringToUint64(s string) uint64 {
	safeString := s
	if s == "" {
		safeString = "0"
	}

	valUint64, err := strconv.ParseUint(safeString, 10, 64)
	if err != nil {
		log.Println("ERROR -> Couldn't convert string to uint")
	}

	return valUint64

}
