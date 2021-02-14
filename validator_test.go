package main

import (
	"net/url"
	"testing"
)

func TestValidateForm(t *testing.T) {
	perfectForm := url.Values{
		"reporter":    []string{"Asena#2614"},
		"mail":        []string{"test@test.com"},
		"title":       []string{"Lorem Ipsum"},
		"description": []string{"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam tincidunt dolor risus, a scelerisque nisl varius in. Phasellus et tortor non enim efficitur dapibus vitae eu est. Cras venenatis metus eget massa mollis, ac blandit tortor suscipit. Vestibulum convallis turpis in gravida interdum."},
		"type":        []string{"Bug"},
		"priority":    []string{"Lowest"},
	}

	if _, ok := ValidateForm(&perfectForm); !ok {
		t.Error("Actual: false, Expected: true")
	}

	badForm := url.Values{
		"reporter":    []string{"foo#1234"},
		"mail":        []string{"badmail@com"},
		"title":       []string{"Lorem Ipsum"},
		"description": []string{"Short Description"},
		"type":        []string{"invalid-type"},
		"priority":    []string{"invalid-priority"},
	}

	if _, ok := ValidateForm(&badForm); ok {
		t.Error("Actual: true, Expected: false")
	}
}
