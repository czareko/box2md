package main

import (
	"encoding/json"
	"strings"
)

type Key string

const (
	DESCRIPTION           Key = "description"
	TYPE                      = "type"
	EXAMPLES                  = "examples"
	ADDITIONAL_PROPERTIES     = "additionalProperties"
)

func isSimpleValue(json json.RawMessage) bool {

	value := string(json)
	if string(value[0]) == "\"" {
		//	 && value != "description" && value != "type" && value != "examples" {
		return true
	}
	return false

}
func isTable(json json.RawMessage) bool {
	value := string(json)
	if string(value[0]) == "[" {
		return true
	}
	return false
}

func detailsOpened(value string) bool {

	var expandedDetails = make(map[string]bool)
	expandedDetails["descriptorWithEnvironments"] = true
	//expandedDetails["descriptor"] = true
	expandedDetails["properties"] = true
	//var expandedDetails []string :={"descriptorWithEnvironments"}
	if expandedDetails[value] {
		return true
	}
	return false
}
func quoteRemover(value string) string {
	return strings.ReplaceAll(value, "\"", "")
}
