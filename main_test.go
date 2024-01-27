package main

import (
	"testing"
)

func TestCheckInput(t *testing.T) {
	inputOne := []string{"test1"}
	str, err := checkInput(inputOne)

	if str != "" || err == nil {
		t.Fatalf("Expected: '', error ; Got: '%s', '%v'", str, err)
	}
}
