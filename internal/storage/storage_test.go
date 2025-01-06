package storage

import (
	"os"
	"testing"
)

func TestMoveFile(t *testing.T) {

	// Setup directories
	os.MkdirAll("./upload", os.ModePerm)
	os.MkdirAll("./temp", os.ModePerm)
	defer os.RemoveAll("./upload")
	defer os.RemoveAll("./temp")

	// Create a test file
	file, err := os.Create("./temp/test.txt")
	if err != nil {
		t.Fatalf("cannot create test file: %v", err)
	}
	file.Close()

	str, err := New("./upload", "./temp")
	if err != nil {
		t.Fatalf("cannot initialize storage: %v", err)
	}

	// Move the file
	err = str.Move("./temp/test.txt", "./upload/test.txt")
	if err != nil {
		t.Fatalf("cannot move test file: %v", err)
	}

	// Check that the file exists in the new location
	if _, err := os.Stat("./upload/test.txt"); os.IsNotExist(err) {
		t.Fatalf("file not found in destination: %v", err)
	}
}
