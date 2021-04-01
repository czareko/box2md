package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"
)

type DictionaryElement struct {
	name string
	path string
	//description string
	//dataType    string
	propertyMap map[string]json.RawMessage
	//properties        map[string]Property
	propertyMapString map[string]string
	jsonValue         string
	jsonElement       JsonElement
}

type JsonElement struct {
	Description string              `json:"description"`
	Datatype    string              `json:"datatype"`
	Properties  map[string]Property `json:"properties"`
}
type Property struct {
	Description          string
	DataType             string
	Default              string
	ReadOnly             string
	items                []Item
	additionalProperties []AdditionalProperty
	examples             string
}
type AdditionalProperty struct {
	DataType   string
	Properties map[string]Property
	Elements   map[string]JsonElement
}
type Item struct {
	DataType string
}

func main() {

	var files = readFiles()

	dict := make(map[string]string)

	for _, file := range files {

		jsonFile, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
		}

		jsonStringFile := string(jsonFile)

		fillDictionary(jsonStringFile, dict)

		var dictionary = make(map[string]DictionaryElement)

		for k := range dict {
			//working version
			variableName := strings.ReplaceAll(k, "#/definitions/", "")
			referenceLine := "\"$ref\": \"" + k + "\""
			newLine := "\"" + variableName + "\": " + dict[k]
			jsonStringFile = strings.ReplaceAll(jsonStringFile, referenceLine, newLine)
		}

		for k := range dict {
			//fmt.Println(k)
			//fmt.Println(dict[k])
			//var propertyMap map[string]json.RawMessage
			var definitionMap map[string]json.RawMessage
			json.Unmarshal([]byte(dict[k]), &definitionMap)
			elem := DictionaryElement{}
			elem.path = k
			variableName := strings.ReplaceAll(k, "#/definitions/", "")
			elem.name = variableName
			elem.jsonElement.Properties = make(map[string]Property)
			makeDictItemFlat(k, definitionMap, &elem, 0, dictionary)

			//elem.prepareStringProperties()
			//elem.convertToJsonElement()

			dictionary[k] = elem
			//fmt.Println("Final Element: " + k)
			//fmt.Println(elem.jsonValue)
		}

		//w aditionalProperties czasami występują elementy złożone, typu volume i descriptor, jak jechałem po słowniku to mogłem ich jeszcze nie mieć, bo kolejność z mapy jest losowa.

		for k := range dictionary {
			for ip := range dictionary[k].jsonElement.Properties {
				//fmt.Println(ip)
				//fmt.Println(dictionary[k].jsonElement.Properties[ip])
				for iap := range dictionary[k].jsonElement.Properties[ip].additionalProperties {
					if dictionary[k].jsonElement.Properties[ip].additionalProperties[iap].DataType == "volume" {
						dictionary[k].jsonElement.Properties[ip].additionalProperties[iap].Elements["volume"] = dictionary["#/definitions/volume"].jsonElement
					} else if dictionary[k].jsonElement.Properties[ip].additionalProperties[iap].DataType == "descriptor" {
						dictionary[k].jsonElement.Properties[ip].additionalProperties[iap].Elements["descriptor"] = dictionary["#/definitions/descriptor"].jsonElement
					}
				}
			}
		}

		//Elementy słownika do MD - test
		for k := range dictionary {
			fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
			//jsonElement := dictionary[k].jsonElement
			//fmt.Println(k)
			//jsonElPointer := dictionary[k].jsonElement
			fmt.Println(dictionary[k].toMD())
		}

	}
}

func makeDictItemFlat(item string, definitionMap map[string]json.RawMessage, element *DictionaryElement, level int, dictionary map[string]DictionaryElement) {

	if string(definitionMap["allOf"]) != "" {
		if level == 0 {
			//element.description = string(definitionMap["description"]) //TO REMOVE
			element.jsonElement.Description = string(definitionMap["description"])
			//element.dataType = string(definitionMap["type"]) // TO REMOVE
			element.jsonElement.Datatype = string(definitionMap["type"])
		}
		//czy to jest obiekt z typem allOf
		jsonTable := make([]json.RawMessage, 0)
		json.Unmarshal(definitionMap["allOf"], &jsonTable)

		for i := 0; i < len(jsonTable); i++ {
			var propertyTempMap map[string]json.RawMessage
			json.Unmarshal(jsonTable[i], &propertyTempMap)
			if string(propertyTempMap["properties"]) != "" {
				json.Unmarshal(propertyTempMap["properties"], &element.propertyMap)
				for pm := range element.propertyMap {
					element.jsonElement.Properties[pm] = upackProperty(element.propertyMap[pm])
					//fmt.Println("PM: " + pm + " ++ " + string(element.propertyMap[pm]))
				}
			} else {
				//po odczytaniu będzie nazwa pola
				var secondLevelMap map[string]json.RawMessage
				json.Unmarshal(jsonTable[i], &secondLevelMap)
				//na kolejnym poziomie powinny być properties
				for k := range secondLevelMap {
					var thirdLevelMap map[string]json.RawMessage
					json.Unmarshal(secondLevelMap[k], &thirdLevelMap)
					if string(thirdLevelMap["allOf"]) != "" {
						level++
						makeDictItemFlat(k, thirdLevelMap, element, level, dictionary)
					} else if string(thirdLevelMap["properties"]) != "" {
						json.Unmarshal(thirdLevelMap["properties"], &element.propertyMap)
						for pm := range element.propertyMap {
							element.jsonElement.Properties[pm] = upackProperty(element.propertyMap[pm])

							//fmt.Println(element.properties[pm].additionalProperties)
							//fmt.Println("PM: " + pm + " ++ " + string(element.propertyMap[pm]))
						}
					}
				}
			}
		}
		res, _ := json.Marshal(definitionMap)
		element.jsonValue = string(res)
	} else {
		if level == 0 {
			element.jsonElement.Description = string(definitionMap["description"])
			element.jsonElement.Datatype = string(definitionMap["type"])
			json.Unmarshal(definitionMap["properties"], &element.propertyMap)
		}

		res, _ := json.Marshal(definitionMap)
		element.jsonValue = string(res)
	}

}

func fillDictionary(jsonString string, dict map[string]string) {

	definitionsMap := gjson.Get(jsonString, "definitions").Map()

	for k := range definitionsMap {
		nv := definitionsMap[k].String()

		if !strings.Contains(nv, "$ref") {
			if dict["#/definitions/"+k] == "" {
				dict["#/definitions/"+k] = nv
			}
		}
	}

	for k := range dict {

		variableName := strings.ReplaceAll(k, "#/definitions/", "")
		referenceLine := "\"$ref\": \"" + k + "\""
		newLine := "\"" + variableName + "\": " + dict[k]
		jsonString = strings.ReplaceAll(jsonString, referenceLine, newLine)
	}

	if strings.Contains(jsonString, "$ref") {
		fillDictionary(jsonString, dict)
	}
}

func upackProperty(jsonValue json.RawMessage) Property {
	//fmt.Println(string(jsonValue))
	property := Property{}
	var definitionMap map[string]json.RawMessage
	json.Unmarshal(jsonValue, &definitionMap)
	for range definitionMap {
		//fmt.Println(k)
		property.Description = string(definitionMap["description"])
		property.DataType = string(definitionMap["type"])
		property.Default = string(definitionMap["default"])
		property.ReadOnly = string(definitionMap["readonly"])
		property.examples = string(definitionMap["examples"])
		if string(definitionMap["additionalProperties"]) != "" {
			//fmt.Println("additionalProperties !!!: " + string(definitionMap["additionalProperties"]))
			defMap := make(map[string]json.RawMessage)
			json.Unmarshal(definitionMap["additionalProperties"], &defMap)
			for ik := range defMap {
				additionalProperty := AdditionalProperty{}
				//realnie to jest błąd w przypadku jeżeli wystąpią inne właściwości niż type
				//fmt.Println("IK: " + ik)
				if ik == "type" {
					additionalProperty.DataType = string(defMap[ik])
				} else if ik == "properties" {
					additionalProperty.Properties = make(map[string]Property)
					//fmt.Println("--------------Properties")
					propMap := make(map[string]json.RawMessage)
					json.Unmarshal(defMap["properties"], &propMap)
					for ip := range propMap {
						property := Property{}
						property.Description = string(propMap["description"])
						property.DataType = string(propMap["type"])
						additionalProperty.Properties[ip] = property
					}
				} else if ik == "volume" {
					additionalProperty.Elements = make(map[string]JsonElement)
					additionalProperty.DataType = "volume"
					//fmt.Println("--------------Volume" + string(defMap[ik]))
					volumeMap := make(map[string]json.RawMessage)
					json.Unmarshal(defMap[ik], &volumeMap)
					element := JsonElement{}
					element.Description = string(volumeMap["description"])
					element.Datatype = string(volumeMap["type"])
					element.Properties = make(map[string]Property)
					propMap := make(map[string]json.RawMessage)
					json.Unmarshal(volumeMap["properties"], &propMap)
					for ip := range propMap {
						property := Property{}
						property.Description = string(propMap["description"])
						property.DataType = string(propMap["type"])
						element.Properties[ip] = property
					}
					additionalProperty.Elements[ik] = element
				} else if ik == "descriptor" {
					additionalProperty.DataType = "descriptor"
					additionalProperty.Elements = make(map[string]JsonElement)
				} else {
					fmt.Println("Unknown additional type: " + ik)
				}
				property.additionalProperties = append(property.additionalProperties, additionalProperty)
			}
			//additionalProperty := AdditionalProperty{}
			//additionalProperty.DataType = "string"
			//property.additionalProperties = append(property.additionalProperties, additionalProperty)
		}
	}
	//fmt.Println(definitionMap)
	return property
}

//-------------------------------------------------------------------------------------------------------------Helpers

func readFiles() []string {

	var files []string

	root := "in"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}

// ------------------------------------------------------------------------------------------------------------MD from structs
func (element DictionaryElement) toMD() string {
	formatedElement := "Name: " + element.name + "\n"
	formatedElement += element.jsonElement.toMD()
	return formatedElement
}
func (item Item) toMD() string {
	return item.DataType
}
func (element JsonElement) toMD() string {
	//fmt.Println(element) //return "txt"
	formatedElement := ""
	if element.Description != "" {
		formatedElement += "Description: " + element.Description + "\n"
	}
	if element.Datatype != "" {
		formatedElement += "Data type: " + element.Datatype + "\n"
	}
	if len(element.Properties) > 0 {
		formatedElement += "Properties: \n"
		for k := range element.Properties {
			formatedElement += "	Name: " + k + "\n"
			formatedElement += element.Properties[k].toMD()
		}
	}
	return formatedElement
}
func (property Property) toMD() string {

	formatedElement := ""
	formatedElement += "	Data type: " + property.DataType + "\n"
	if property.Default != "" {
		formatedElement += "	Default: " + property.Default + "\n"
	}
	if property.ReadOnly != "" {
		formatedElement += "	ReadOnly: " + property.ReadOnly + "\n"
	}
	if len(property.additionalProperties) > 0 {

		formatedElement += "		Additional properties: \n"
		for k := range property.additionalProperties {
			formatedElement += property.additionalProperties[k].toMD()
		}
	}
	if property.examples != "" {
		formatedElement += "	Example: " + property.examples + "\n"
	}

	formatedElement += "\n"

	return formatedElement
}
func (additionalProp AdditionalProperty) toMD() string {
	formatedElement := ""

	for k := range additionalProp.Properties {
		formatedElement += "			Name: " + k + "\n"
		formatedElement += additionalProp.Properties[k].toMD()
	}

	return formatedElement
}
