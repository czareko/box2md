package box2md

import (
	"strings"
)

type TypeDefinition struct {
	RefTypeDefinitions   map[string]*TypeDefinition `json:"-"`
	Schema               *string                    `json:"$schema,omitempty"`
	Id                   *string                    `json:"$id,omitempty"`
	Type                 *string                    `json:"type,omitempty"`
	Description          *string                    `json:"description,omitempty"`
	Properties           map[string]*TypeDefinition `json:"properties,omitempty"`
	AllOf                []*TypeDefinition          `json:"allOf,omitempty"`
	AnyOf                []*TypeDefinition          `json:"anyOf,omitempty"`
	OneOf                []*TypeDefinition          `json:"oneOf,omitempty"`
	Ref                  *string                    `json:"$ref,omitempty"`
	Examples             []interface{}              `json:"examples,omitempty"`
	ReadOnly             *bool                      `json:"readOnly,omitempty"`
	Items                *TypeDefinition            `json:"items,omitempty"`
	AdditionalProperties *TypeDefinition            `json:"additionalProperties,omitempty"`
	Enum                 []string                   `json:"enum,omitempty"`
	Default              interface{}                `json:"default,omitempty"`
	Required             []string                   `json:"required,omitempty"`
}

func (t *TypeDefinition) RefKey() *string {
	if t.Ref == nil {
		return nil
	} else {
		key := strings.TrimPrefix(*t.Ref, "#/definitions/")
		return &key
	}
}

func (t *TypeDefinition) propagateRefTypeDefinitions(definitions map[string]*TypeDefinition) {
	t.RefTypeDefinitions = definitions
	for _, referencedTypeDefinition := range t.Properties {
		referencedTypeDefinition.propagateRefTypeDefinitions(definitions)
	}
	for _, referencedTypeDefinition := range t.AllOf {
		referencedTypeDefinition.propagateRefTypeDefinitions(definitions)
	}
	for _, referencedTypeDefinition := range t.AnyOf {
		referencedTypeDefinition.propagateRefTypeDefinitions(definitions)
	}
	for _, referencedTypeDefinition := range t.OneOf {
		referencedTypeDefinition.propagateRefTypeDefinitions(definitions)
	}
	if t.Items != nil {
		t.Items.propagateRefTypeDefinitions(definitions)
	}
	if t.AdditionalProperties != nil {
		t.AdditionalProperties.propagateRefTypeDefinitions(definitions)
	}
}

func (t *TypeDefinition) GetSchema() *string {
	if t.Schema != nil {
		return t.Schema
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetSchema()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetId() *string {
	if t.Id != nil {
		return t.Id
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetId()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetType() *string {
	if t.Type != nil {
		return t.Type
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetType()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetDescription() *string {
	if t.Description != nil {
		return t.Description
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetDescription()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetProperties() map[string]*TypeDefinition {
	if t.Properties != nil {
		return t.Properties
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetProperties()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetAllProperties() map[string]*TypeDefinition {
	allProperties := t.GetProperties()
	for _, child := range t.GetAllOf() {
		for key, value := range child.GetAllProperties() {
			if allProperties == nil {
				allProperties = make(map[string]*TypeDefinition)
			}
			allProperties[key] = value
		}
	}
	return allProperties
}

func (t *TypeDefinition) GetAllOf() []*TypeDefinition {
	if t.AllOf != nil {
		return t.AllOf
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetAllOf()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetAnyOf() []*TypeDefinition {
	if t.AnyOf != nil {
		return t.AnyOf
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetAnyOf()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetOneOf() []*TypeDefinition {
	if t.OneOf != nil {
		return t.OneOf
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetOneOf()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetExamples() []interface{} {
	if t.Examples != nil {
		return t.Examples
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetExamples()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetReadOnly() *bool {
	if t.ReadOnly != nil {
		return t.ReadOnly
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetReadOnly()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetItems() *TypeDefinition {
	if t.Items != nil {
		return t.Items
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetItems()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetAdditionalProperties() *TypeDefinition {
	if t.AdditionalProperties != nil {
		return t.AdditionalProperties
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetAdditionalProperties()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetEnum() []string {
	if t.Enum != nil {
		return t.Enum
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetEnum()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetDefault() interface{} {
	if t.Default != nil {
		return t.Default
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetDefault()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetRequired() []string {
	if t.Required != nil {
		return t.Required
	} else {
		refKey := t.RefKey()
		if refKey != nil {
			return t.RefTypeDefinitions[*refKey].GetRequired()
		} else {
			return nil
		}
	}
}

func (t *TypeDefinition) GetRequiredSet() map[string]struct{} {
	required := t.GetRequired()
	result := make(map[string]struct{}, len(required))
	for _, element := range required {
		result[element] = struct{}{}
	}
	return result
}

func isSimpleType(t *string) bool {
	if t == nil {
		return false
	} else {
		switch *t {
		case "object":
			return false
		case "array":
			return false
		default:
			return true
		}
	}
}
