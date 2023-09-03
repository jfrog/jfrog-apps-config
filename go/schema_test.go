package jfrogappsconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

var loadConfigIfExistCases = []struct {
	fileName       string
	expectedConfig JFrogAppsConfig
}{
	{fileName: "no-scanners.yml", expectedConfig: JFrogAppsConfig{
		Version: "1.0",
		Modules: []Module{{
			Name:            "NoScanners",
			SourceRoot:      "src",
			ExcludePatterns: []string{"docs/"},
			ExcludeScanners: []string{"secrets"}},
		}},
	},
	{fileName: "sast.yml", expectedConfig: JFrogAppsConfig{
		Version: "1.0",
		Modules: []Module{{
			Name: "Sast",
			Scanners: Scanners{Sast: &SastScanner{
				Language: "java",
				Scanner: Scanner{
					WorkingDirs:     []string{"src/module1", "src/module2"},
					ExcludePatterns: []string{"src/module1/test"},
				},
				ExcludedRules: []string{"xss-injection"},
			}},
		}}},
	},
	{fileName: "all-scanners.yml", expectedConfig: JFrogAppsConfig{
		Version: "1.0",
		Modules: []Module{{
			Name: "AllScanners",
			Scanners: Scanners{Sast: &SastScanner{
				Language: "java",
				Scanner:  Scanner{WorkingDirs: []string{"src/module1"}}},
				Iac:     &Scanner{WorkingDirs: []string{"src/module2"}},
				Secrets: &Scanner{WorkingDirs: []string{"src/module3"}},
			},
		}}},
	},
}

func TestLoadConfigIfExist(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "load-config-test")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tempDir))
	}()
	assert.NoError(t, os.Mkdir(filepath.Join(tempDir, dotJFrogDir), 0750))

	cwd, err := os.Getwd()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, os.Chdir(cwd))
	}()

	for _, testCase := range loadConfigIfExistCases {
		t.Run(testCase.fileName, func(t *testing.T) {
			assert.NoError(t, os.Chdir(cwd))

			content, err := os.ReadFile(filepath.Join("testdata", "goodschemas", testCase.fileName))
			assert.NoError(t, err)

			assert.NoError(t, os.Chdir(tempDir))
			assert.NoError(t, os.WriteFile(filepath.Join(dotJFrogDir, configFileName), content, 0600))

			jfrogAppsConfig, err := LoadConfigIfExist()
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedConfig, *jfrogAppsConfig)
		})
	}
}

func TestJFrogYamlSchema(t *testing.T) {
	// Load JFrog Apps Config schema
	schema, err := os.ReadFile(filepath.Join("..", "schema.json"))
	assert.NoError(t, err)
	schemaLoader := gojsonschema.NewBytesLoader(schema)

	// Validate all JFrog Apps Configs in testdata
	validateYamlsInDirectory(t, filepath.Join("testdata", "goodschemas"), schemaLoader)

	// Validate bad schema
	validateYamlSchema(t, schemaLoader, filepath.Join("testdata", "badschemas", "bad-schema.yml"), "Invalid type. Expected: string, given: integer")

}

// Validate all yml files in the given directory against the input schema
// t            - Testing object
// schemaLoader - JFrog Apps Config schema
// path	         - Yaml directory path
func validateYamlsInDirectory(t *testing.T, path string, schemaLoader gojsonschema.JSONLoader) {
	err := filepath.Walk(path, func(jfrogAppConfigFilePath string, info os.FileInfo, err error) error {
		assert.NoError(t, err)
		if strings.HasSuffix(info.Name(), "yml") {
			validateYamlSchema(t, schemaLoader, jfrogAppConfigFilePath, "")
		}
		return nil
	})
	assert.NoError(t, err)
}

// Validate a Yaml file against the input Yaml schema
// t            - Testing object
// schemaLoader - JFrog Apps Config schema
// yamlFilePath - Yaml file path
// expectError  - Expected error or an empty string if error is not expected
func validateYamlSchema(t *testing.T, schemaLoader gojsonschema.JSONLoader, yamlFilePath, expectError string) {
	t.Run(filepath.Base(yamlFilePath), func(t *testing.T) {
		// Read JFrog Apps Config
		// #nosec G304 - false positive
		yamlFile, err := os.ReadFile(yamlFilePath)
		assert.NoError(t, err)

		// Unmarshal JFrog Apps Config
		var jfrogAppsConfigYaml interface{}
		err = yaml.Unmarshal(yamlFile, &jfrogAppsConfigYaml)
		assert.NoError(t, err)

		// Convert the Yaml config to JSON config to help the json parser validate it.
		// The reason we don't do the convert by as follows:
		// YAML -> Unmarshall -> Go Struct -> Marshal -> JSON
		// is because the config's struct includes only YAML annotations.
		jfrogAppsConfigJson := convertYamlToJson(jfrogAppsConfigYaml)

		// Load and validate JFrog Apps Config
		documentLoader := gojsonschema.NewGoLoader(jfrogAppsConfigJson)
		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		assert.NoError(t, err)
		if expectError != "" {
			assert.False(t, result.Valid())
			assert.Contains(t, result.Errors()[0].String(), expectError)
		} else {
			assert.True(t, result.Valid(), result.Errors())
		}
	})
}

// Recursively convert yaml interface to JSON interface
func convertYamlToJson(yamlValue interface{}) interface{} {
	switch yamlMapping := yamlValue.(type) {
	case map[interface{}]interface{}:
		jsonMapping := map[string]interface{}{}
		for key, value := range yamlMapping {
			if key == true {
				// "on" is considered a true value for the Yaml Unmarshaler. To work around it, we set the true to be "on".
				key = "on"
			}
			jsonMapping[fmt.Sprint(key)] = convertYamlToJson(value)
		}
		return jsonMapping
	case []interface{}:
		for i, value := range yamlMapping {
			yamlMapping[i] = convertYamlToJson(value)
		}
	}
	return yamlValue
}
