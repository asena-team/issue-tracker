package main

func Compare(n, max, min int) bool {
	return n <= max && n >= min
}

func Contains(i string, a []string) bool {
	for _, n := range a {
		if i == n {
			return true
		}
	}

	return false
}
