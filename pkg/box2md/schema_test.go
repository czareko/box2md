package box2md

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestReadSchema(t *testing.T) {
	t.Parallel()

	tempFilePath := createTempFile(t)
	defer os.Remove(tempFilePath)
	schemaJsonPath := filepath.Join("testdata", "reader", "box-yaml.schema.json")

	schema := Read(schemaJsonPath)

	schema.Write(tempFilePath)
	diff := cmp.Diff(readJsonToMap(t, schemaJsonPath), readJsonToMap(t, tempFilePath))
	if diff != "" {
		t.Fatalf(diff)
	}
}

func createTempFile(t *testing.T) string {
	tempFile, err := ioutil.TempFile(os.TempDir(), "box-yaml.schema-*.json")
	if err != nil {
		t.Fatalf("Error while creating temp file")
	}
	tempFilePath := tempFile.Name()
	return tempFilePath
}

func readJsonToMap(t *testing.T, path string) map[string]interface{} {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Cannot read file %v", path)
	}

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(file, &jsonMap)
	if err != nil {
		t.Fatalf("Cannot unmarshal file %v", path)
	}

	return jsonMap
}
