package box2md

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type JsonSchema struct {
	Title          *string `json:"title,omitempty"`
	TypeDefinition `json:",inline"`
	Definitions    map[string]*TypeDefinition `json:"definitions,omitempty"`
}

func Read(path string) JsonSchema {
	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicf("Error %v reading json schema file %v", err, path)
	}

	result := JsonSchema{}
	err = json.Unmarshal(jsonFile, &result)
	if err != nil {
		log.Panicf("Error %v unmarshalling json file %v", err, path)
	}

	for _, definition := range result.Definitions {
		definition.propagateRefTypeDefinitions(result.Definitions)
	}
	result.propagateRefTypeDefinitions(result.Definitions)

	return result
}

func (s *JsonSchema) Write(path string) {
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		log.Panicf("Error %v marshalling schema file %v", err, path)
	}

	err = ioutil.WriteFile(path, jsonBytes, 0644)
	if err != nil {
		log.Panicf("Error %v writing schema file %v", err, path)
	}
}
