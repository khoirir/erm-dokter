package docs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type openAPISpec struct {
	Paths map[string]map[string]struct {
		Summary   string `yaml:"summary"`
		Responses map[string]struct {
			Description string `yaml:"description"`
			Content     map[string]struct {
				Schema struct {
					Ref        string                 `yaml:"$ref"`
					AllOf      []map[string]any       `yaml:"allOf"`
					Properties map[string]struct {
						Type    string `yaml:"type"`
						Example any    `yaml:"example"`
					} `yaml:"properties"`
				} `yaml:"schema"`
			} `yaml:"content"`
		} `yaml:"responses"`
	} `yaml:"paths"`
}

func TestVerifyAllModuleResponsesHaveExplicitExamples(t *testing.T) {
	entries, err := os.ReadDir("modules")
	if err != nil {
		t.Fatalf("failed to read modules dir: %v", err)
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		filePath := filepath.Join("modules", entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read %s: %v", entry.Name(), err)
		}

		var spec openAPISpec
		if err := yaml.Unmarshal(content, &spec); err != nil {
			t.Fatalf("failed to parse yaml %s: %v", entry.Name(), err)
		}

		for path, methods := range spec.Paths {
			for method, op := range methods {
				if method != "get" && method != "post" && method != "put" && method != "delete" && method != "patch" {
					continue
				}

				// Check 200 or 201 response
				resp200, ok200 := op.Responses["200"]
				resp201, ok201 := op.Responses["201"]

				if !ok200 && !ok201 {
					t.Errorf("[%s] %s %s has neither 200 nor 201 response", entry.Name(), strings.ToUpper(method), path)
					continue
				}

				resp := resp200
				respCode := "200"
				if !ok200 && ok201 {
					resp = resp201
					respCode = "201"
				}

				jsonContent, hasJSON := resp.Content["application/json"]
				if !hasJSON {
					// Some endpoints might return image/jpeg or raw stream
					if path == "/api/v1/berkas-digital/{id_berkas}" {
						continue
					}
					t.Errorf("[%s] %s %s response %s missing application/json content", entry.Name(), strings.ToUpper(method), path, respCode)
					continue
				}

				if jsonContent.Schema.Ref != "" {
					t.Errorf("[%s] %s %s response %s uses $ref '%s' instead of explicit envelope schema", entry.Name(), strings.ToUpper(method), path, respCode, jsonContent.Schema.Ref)
				}

				if len(jsonContent.Schema.AllOf) > 0 {
					t.Errorf("[%s] %s %s response %s uses allOf instead of explicit envelope schema", entry.Name(), strings.ToUpper(method), path, respCode)
				}

				msgProp, hasMsg := jsonContent.Schema.Properties["message"]
				if !hasMsg {
					t.Errorf("[%s] %s %s response %s schema missing 'message' property", entry.Name(), strings.ToUpper(method), path, respCode)
					continue
				}

				if msgProp.Example == nil || msgProp.Example == "" {
					t.Errorf("[%s] %s %s response %s message property missing 'example'", entry.Name(), strings.ToUpper(method), path, respCode)
					continue
				}

				msgExample, isStr := msgProp.Example.(string)
				if !isStr {
					t.Errorf("[%s] %s %s response %s message example is not a string: %v", entry.Name(), strings.ToUpper(method), path, respCode, msgProp.Example)
					continue
				}

				if msgExample == "Berhasil memproses data" {
					t.Errorf("[%s] %s %s response %s message example is generic: '%s'", entry.Name(), strings.ToUpper(method), path, respCode, msgExample)
				}
				t.Logf("[%s] %s %s -> %s: %s", entry.Name(), strings.ToUpper(method), path, respCode, msgExample)
			}
		}
	}
}
