package main

import (
	"testing"
)

func TestFetchAll(t *testing.T) {
	body, err := fetchAll()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if body != "" {
		t.Errorf("unexpected body: %s", body)
	}
}
