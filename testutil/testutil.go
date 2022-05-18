package testutil

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func CreateFile(t *testing.T, path string, data map[string]interface{}) {
	t.Helper()

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		t.Fatalf("failed to encode testdata")
	}
	if err := os.WriteFile(path, buf.Bytes(), os.ModePerm); err != nil {
		t.Fatal("failed to os.WriteFile")
	}
}

func CaptureOutput(t *testing.T, fn func()) string {
	t.Helper()

	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read buf: %v", err)
	}
	return strings.TrimRight(buf.String(), "\n")
}