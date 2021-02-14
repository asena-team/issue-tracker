package main

import (
	"net/url"
	"reflect"
	"strings"
)

type InputValue struct {
	name     string
	pattern  string
	min, max int
}

type SelectorValue struct {
	name  string
	array []string
}

var values = []interface{}{
	InputValue{"title", "", 1, 150},
	InputValue{"description", "", 30, 3000},
	InputValue{"reporter", ReporterRegEx, 1, 50},
	InputValue{"mail", EMailRegEx, 0, 150},
	SelectorValue{"type", IssueTypes},
	SelectorValue{"priority", Priorities},
}

func ValidateForm(form *url.Values) (*Issue, bool) {
	issue := &Issue{}
	v := reflect.ValueOf(issue).Elem()
	for _, item := range values {
		switch reflect.TypeOf(item) {
		case reflect.TypeOf(InputValue{}):
			item, _ := item.(InputValue)
			str := strings.Trim(form.Get(item.name), " ")
			if str == "" || !Compare(len(str), item.min, item.max) {
				return issue, false
			}

			if item.pattern != "" {
				if !Match(item.pattern, str) {
					return issue, false
				}
			}

			field := v.FieldByName(strings.Title(item.name))
			field.SetString(str)
			break

		case reflect.TypeOf(SelectorValue{}):
			item, _ := item.(SelectorValue)
			str := strings.Trim(form.Get(item.name), " ")
			if !Contains(str, item.array) {
				return issue, false
			}

			field := v.FieldByName(strings.Title(item.name))
			field.SetString(str)
			break
		}
	}

	return issue, true
}
