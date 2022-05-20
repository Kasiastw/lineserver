package helpers

import (
	"log"
	"regexp"
)

func GetRegex(text string) bool {
	match, err := regexp.MatchString("(?i)"+"(the)", text)
	if err != nil {
		log.Println("error matching string: err:", err)
	}
	return match
}


