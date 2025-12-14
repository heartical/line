package validation

import (
	"sort"
	"strings"
)

type TemplateParameter struct {
	Key   string
	Value string
}

type TemplateParameterList []TemplateParameter

func (params TemplateParameterList) Prepend(parameters ...TemplateParameter) TemplateParameterList {
	return append(parameters, params...)
}

func renderMessage(template string, parameters []TemplateParameter) string {
	sort.SliceStable(parameters, func(i, j int) bool {
		return len(parameters[i].Key) > len(parameters[j].Key)
	})

	message := template
	for _, p := range parameters {
		message = strings.ReplaceAll(message, p.Key, p.Value)
	}

	return message
}
