package helpers

import (
	"encoding/json"
	"log"
	"math"
)

func AbsInt64(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func Round(v float64) float64 {
	return math.Round(v*1000) / 1000
}

func PrettyPrint(msg string, obj interface{}) {
	log.Println(msg, Pretty(obj))
}

func Pretty(obj interface{}) string {
	str, err := json.MarshalIndent(obj, "", "   ")
	if err != nil {
		log.Fatalln(err.Error())
	}
	return string(str)
}
