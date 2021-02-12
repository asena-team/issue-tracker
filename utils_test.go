package main

import "testing"

func TestMatch(t *testing.T) {
	if res := Match(ReporterRegEx, "test_foo_işö$#2001"); !res {
		t.Errorf("Expected: true, Actual: %t", res)
	}

	if res := Match(EMailRegEx, "test@test.com"); !res {
		t.Errorf("Expected: true, Actual: %t", res)
	}
}
