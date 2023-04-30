package controller

import (
	"log"
	"strconv"
	"net/http"
	"html/template"
)


func AtoiSlice(seq []string) []int {
	result := make([]int, len(seq))
	for index, el := range seq {
		number, err := strconv.Atoi(el)
		if err != nil {
			log.Printf("[WARNING] - Could not convert string \"%v\" to int.", el)
			return result
		}
		result[index] = number
	}
	return result
}


func ServeFile(webpage http.ResponseWriter, pageName string, object interface{}) {
	content, err := template.ParseFiles("./view/html/" + pageName)
	if err != nil {
		log.Printf("[ERROR] - File \"%v\" does not exist or is not accessible.\n%v",
					pageName, err.Error())
	}
	err = content.ExecuteTemplate(webpage, pageName, object)
	if err != nil {
		log.Printf("[ERROR] - Template execution.\n" + err.Error() + "\n\n")
	}
}


func GetClosestTo(seq []int, target int) int {
	min, min_diff := -1, 99999999999999
	for _, value := range seq {
		if value < target {
			continue
		}
		if value - target < min_diff {
			min, min_diff = value, value - target
		}
	}
	return min
}