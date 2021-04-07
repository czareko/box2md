package box2md

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
)

func (t *TypeDefinition) ToMDString() string {
	result := ""

	result += t.printDescription()
	result += t.printType()
	result += t.printAnyOf()
	result += t.printOneOf()
	result += t.printProperties()
	result += t.printExamples()
	result += t.printEnum()
	result += t.printReadOnly()
	result += t.printItems()
	result += t.printAdditionalProperties()
	result += t.printDefault()

	return result
}

func (t *TypeDefinition) printDefault() string {
	defaultValue := t.GetDefault()
	if defaultValue != nil {
		return fmt.Sprintf("Default value: %v", defaultValue)
	} else {
		return ""
	}
}

func (t *TypeDefinition) printAdditionalProperties() string {
	additionalProperties := t.GetAdditionalProperties()
	if additionalProperties != nil {
		result := "Additional properties type: {\n"
		result += additionalProperties.ToMDString()
		result += "}"

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printItems() string {
	items := t.GetItems()
	if items != nil {
		result := "Array items type: {\n"
		result += items.ToMDString()
		result += "}"

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printReadOnly() string {
	readOnly := t.GetReadOnly()
	if readOnly != nil && *readOnly {
		return "Readonly!"
	} else {
		return ""
	}
}

func (t *TypeDefinition) printEnum() string {
	enum := t.GetEnum()
	if len(enum) > 0 {
		result := "Enum: {\n"
		for _, enumValue := range enum {
			result += enumValue + " "
		}
		result += "}\n"

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printExamples() string {
	examples := t.GetExamples()
	if len(examples) > 0 {
		result := "Examples: {\n"
		for _, example := range examples {
			result += toYaml(example)
		}
		result += "}\n"

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printProperties() string {
	requiredProperties := t.GetRequiredSet()

	result := ""
	for key, value := range t.GetAllProperties() {
		result += "Property: " + key + "{\n"
		result += value.ToMDString()

		if _, ok := requiredProperties[key]; ok {
			result += "Required! \n"
		}

		result += "}\n"
	}
	return result
}

func (t *TypeDefinition) printOneOf() string {
	oneOf := t.GetOneOf()
	if len(oneOf) > 0 {
		result := "One of: {\n"
		for _, one := range oneOf {
			result += one.ToMDString()
		}
		result += "}\n"

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printAnyOf() string {
	anyOf := t.GetAnyOf()
	if len(anyOf) > 0 {
		result := "Any: {\n"
		for _, any := range anyOf {
			result += any.ToMDString()
		}
		result += "}\n"

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printType() string {
	getType := t.GetType()
	if isSimpleType(getType) {
		return "Type: " + *getType + "\n"
	} else {
		return ""
	}
}

func (t *TypeDefinition) printDescription() string {
	description := t.GetDescription()
	if description != nil {
		return *description + "\n"
	} else {
		return ""
	}
}

func toYaml(obj interface{}) string {
	yamlBytes, err := yaml.Marshal(obj)
	if err != nil {
		log.Panicf("Error marshaling yaml example!")
	}
	return string(yamlBytes)
}
