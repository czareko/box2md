package box2md

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func (t *TypeDefinition) ToMDString(level int) string {
	result := ""

	result += t.printDescription(level)
	result += t.printOneOf(level)
	result += t.printProperties(level)

	result += t.printItems(level)
	result += t.printAdditionalProperties(level)

	result += t.printExamples(level)

	return result
}

func printIndent(level int) string {
	return strings.Repeat(">", level)
}

func (t *TypeDefinition) printAdditionalProperties(level int) string {
	additionalProperties := t.GetAdditionalProperties()
	if additionalProperties != nil {
		result := printIndent(level) + " Map value structure: <br/>\n"
		result += additionalProperties.ToMDString(level)

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printItems(level int) string {
	items := t.GetItems()
	if items != nil {
		result := printIndent(level) + " Array items type: <br/>\n"
		result += items.ToMDString(level)

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printExamples(level int) string {
	examples := t.GetExamples()
	if len(examples) > 0 {
		result := printIndent(level) + " Examples: \n"
		result += printIndent(level) + "\n"
		for _, example := range examples {
			result += printIndent(level) + " ```yaml\n"
			result += prependEveryLine(toYaml(example), printIndent(level)+" ")
			result += "```\n"
		}

		return result
	} else {
		return ""
	}
}

func prependEveryLine(str string, prefix string) string {
	exp := regexp.MustCompile(`(?m)^`)
	return exp.ReplaceAllString(str, prefix)
}

func (t *TypeDefinition) printProperties(level int) string {
	requiredProperties := t.GetRequiredSet()

	result := ""
	for key, value := range t.GetAllProperties() {
		result += printIndent(level) + " <details>\n"
		result += printIndent(level) + " <summary>" + key + "</summary>\n"
		result += printIndent(level) + " \n"

		required := isPropertyRequired(requiredProperties, key)

		propertyType := value.GetType()
		switch {
		case isSimpleType(propertyType):
			result += t.printSimpleProperty(level, value, required, *value.GetType(), value.GetEnum())
		case isArrayType(propertyType):
			items := value.GetItems()
			if isSimpleType(items.GetType()) {
				typeDescription := "Array of " + *items.GetType() + " values"
				result += t.printSimpleProperty(level, value, required, typeDescription, items.GetEnum())
			} else {
				result += value.ToMDString(level + 1)
			}
		default:
			result += value.ToMDString(level + 1)
		}
		result += printIndent(level) + " </details>\n"
	}
	return result
}

func (t *TypeDefinition) printSimpleProperty(level int, propertyTypeDefinition *TypeDefinition, required bool, typeDescription string, enumValues []string) string {
	result := ""

	result += printIndent(level) + " **Description:** " + *propertyTypeDefinition.GetDescription() + "<br/>\n"

	result += printIndent(level) + " **Type:** " + typeDescription + "<br/>\n"

	if len(enumValues) > 0 {
		result += printIndent(level) + " **Enum:** [ "
		for i, enumValue := range enumValues {
			if i != 0 {
				result += ", "
			}
			result += enumValue
		}
		result += " ]<br/>\n"
	}

	readOnly := propertyTypeDefinition.GetReadOnly()
	if readOnly != nil {
		result += printIndent(level) + " **Read only:** " + strconv.FormatBool(*readOnly) + "<br/>\n"
	}

	result += printIndent(level) + " **Reguired:** " + strconv.FormatBool(required) + "<br/>\n"

	defaultValue := propertyTypeDefinition.GetDefault()
	if defaultValue != nil {
		result += printIndent(level) + " **Default value:** " + defaultValueToString(defaultValue) + "<br/>\n"
	}
	return result
}

func defaultValueToString(defaultValue interface{}) string {
	switch val := defaultValue.(type) {
	case string:
		return "\"" + val + "\""
	default:
		return fmt.Sprint(val)
	}
}

func (t *TypeDefinition) printOneOf(level int) string {
	oneOf := t.GetOneOf()
	if len(oneOf) > 0 {
		result := printIndent(level) + " One of: \n"
		result += printIndent(level) + "\n"
		for _, one := range oneOf {
			result += printIndent(level) + "<details open>\n"
			result += printIndent(level) + "<summary>" + *one.GetDescription() + "</summary>\n"
			result += printIndent(level) + "\n"
			result += one.ToMDString(level + 1)
			result += printIndent(level) + "</details>\n"
		}

		return result
	} else {
		return ""
	}
}

func (t *TypeDefinition) printDescription(level int) string {
	description := t.GetDescription()
	if description != nil {
		return printIndent(level) + " " + *t.GetDescription() + "<br/>\n"
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
