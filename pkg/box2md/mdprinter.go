package box2md

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

func (t *TypeDefinition) ToMDString(level int) string {
	result := ""

	//result += printIndent(level) + fmt.Sprintf("ON LEVEL %v \n", level)

	result += t.printDescription(level)
	result += t.printType(level)
	result += t.printAnyOf(level)
	result += t.printOneOf(level)
	result += t.printProperties(level)
	result += t.printExamples(level)
	result += t.printEnum(level)
	result += t.printReadOnly(level)
	result += t.printItems(level)
	result += t.printAdditionalProperties(level)
	result += t.printDefault(level)

	return result
}

func printIndent(level int) string {
	return strings.Repeat("  ", level)
}

func (t *TypeDefinition) printDefault(level int) string {
	defaultValue := t.GetDefault()
	if defaultValue != nil {
		return printIndent(level) + fmt.Sprintf("Default value: %v\n", defaultValue)
	} else {
		return ""
	}
}

func (t *TypeDefinition) printAdditionalProperties(level int) string {
	additionalProperties := t.GetAdditionalProperties()
	if additionalProperties != nil {
		result := printIndent(level) + "Additional properties type: {\n"
		result += additionalProperties.ToMDString(level + 1)
		result += printIndent(level) + "}\n"

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printItems(level int) string {
	items := t.GetItems()
	if items != nil {
		result := printIndent(level) + "Array items type: {\n"
		result += items.ToMDString(level + 1)
		result += printIndent(level) + "}\n"

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printReadOnly(level int) string {
	readOnly := t.GetReadOnly()
	if readOnly != nil && *readOnly {
		return printIndent(level) + "Readonly!\n"
	} else {
		return ""
	}
}

func (t *TypeDefinition) printEnum(level int) string {
	enum := t.GetEnum()
	if len(enum) > 0 {
		result := printIndent(level) + "Enum: { "
		for i, enumValue := range enum {
			if i != 0 {
				result += ", "
			}
			result += enumValue
		}
		result += " }\n"

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printExamples(level int) string {
	examples := t.GetExamples()
	if len(examples) > 0 {
		result := printIndent(level) + "Examples: {\n"
		for _, example := range examples {
			result += toYaml(example)
		}
		result += printIndent(level) + "}\n"

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printProperties(level int) string {
	requiredProperties := t.GetRequiredSet()

	result := ""
	for key, value := range t.GetAllProperties() {
		result += printIndent(level)
		if _, ok := requiredProperties[key]; ok {
			result += "Required property "
		} else {
			result += "Property "
		}
		result += key + ": {\n"

		result += value.ToMDString(level + 1)
		result += printIndent(level) + "}\n"
	}
	return result
}

func (t *TypeDefinition) printOneOf(level int) string {
	oneOf := t.GetOneOf()
	if len(oneOf) > 0 {
		result := printIndent(level) + "One of: {\n"
		for _, one := range oneOf {
			result += one.ToMDString(level + 1)
		}
		result += printIndent(level) + "}\n"

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printAnyOf(level int) string {
	anyOf := t.GetAnyOf()
	if len(anyOf) > 0 {
		result := printIndent(level) + "Any: {\n"
		for _, any := range anyOf {
			result += any.ToMDString(level + 1)
		}
		result += printIndent(level) + "}\n"

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printType(level int) string {
	getType := t.GetType()
	if isSimpleType(getType) {
		return printIndent(level) + "Type: " + *getType + "\n"
	} else {
		return ""
	}
}

func (t *TypeDefinition) printDescription(level int) string {
	description := t.GetDescription()
	if description != nil {
		return printIndent(level) + "Description: " + *description + "\n"
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
