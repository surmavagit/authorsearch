package main

import (
	"testing"
)

func TestCheckInput(t *testing.T) {
	emptyInput := []string{}
	_, err := checkInput(emptyInput)
	if err == nil {
		t.Fatal("Expected error, got no error")
	}

	excessiveInput := []string{"a", "b", "c", "d"}
	_, err = checkInput(excessiveInput)
	if err == nil {
		t.Fatal("Expected error, got no error")
	}

	tooManyNumerics := []string{"1", "2", "a"}
	_, err = checkInput(tooManyNumerics)
	if err == nil {
		t.Fatal("Expected error, got no error")
	}

	tooManyNames := []string{"a", "b", "c"}
	_, err = checkInput(tooManyNames)
	if err == nil {
		t.Fatal("Expected error, got no error")
	}
}
