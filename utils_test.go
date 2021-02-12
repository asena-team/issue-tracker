package main

import "testing"

type Range struct {
	min, max int
}

func TestCompare(t *testing.T) {
	cases := map[int]Range{
		1:     {1, 9},
		100:   {100, 200},
		1234:  {1234, 5678},
		12345: {12345, 12345},
	}

	for k, v := range cases {
		if res := Compare(k, v.min, v.max); !res {
			t.Errorf("Range error: %d (Min: %d, Max: %d)", k, v.min, v.max)
		}
	}
}

func TestContains(t *testing.T) {
	arr := []string{"a", "b", "c", "d", "e", "f"}
	if res := Contains("a", arr); !res {
		t.Errorf("Item not found: %s", "a")
	}

	if res := Contains("g", arr); res {
		t.Errorf("Item found: %s", "g")
	}
}

func TestMatch(t *testing.T) {
	if res := Match(ReporterRegEx, "test_foo_işö$#2001"); !res {
		t.Errorf("Expected: true, Actual: %t", res)
	}

	if res := Match(EMailRegEx, "test@test.com"); !res {
		t.Errorf("Expected: true, Actual: %t", res)
	}
}
