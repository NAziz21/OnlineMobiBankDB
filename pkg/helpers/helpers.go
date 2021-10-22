package helpers

import (
	"regexp"
	"strconv"
	"strings"
)

const key =`aaa`

//Lower Case
func TittleName(OString string) string {
	newString := strings.Title(strings.ToLower(OString))

	return newString
}

func ToUpper (oldString string) string {
	newString := strings.ToUpper(oldString)
	return newString
}

func StrInt(OldString string) (int,error) {
	nInt, err := strconv.ParseInt(OldString,10,0)
	if err != nil {
		return 0, err
	}
	return int(nInt), nil
}

func PhoneNumber (aString string) string {
   newString := regexp.MustCompile(`\D+`).ReplaceAllString(aString , "")
   return newString
}
