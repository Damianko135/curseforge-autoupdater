package env

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// ConfigTemplate holds the default config structure
type ConfigTemplate struct {
	APIKey          string   `json:"api_key" yaml:"api_key" toml:"api_key"`
	ModID           int      `json:"mod_id" yaml:"mod_id" toml:"mod_id"`
	ModName         string   `json:"mod_name" yaml:"mod_name" toml:"mod_name"`
	ModVersion      string   `json:"mod_version" yaml:"mod_version" toml:"mod_version"`
	ModDescription  string   `json:"mod_description" yaml:"mod_description" toml:"mod_description"`
	ModAuthor       string   `json:"mod_author" yaml:"mod_author" toml:"mod_author"`
	ModLicense      string   `json:"mod_license" yaml:"mod_license" toml:"mod_license"`
	ModURL          string   `json:"mod_url" yaml:"mod_url" toml:"mod_url"`
	ModLogo         string   `json:"mod_logo" yaml:"mod_logo" toml:"mod_logo"`
	ModCategories   []string `json:"mod_categories" yaml:"mod_categories" toml:"mod_categories"`
	ModGameVersions []string `json:"mod_game_versions" yaml:"mod_game_versions" toml:"mod_game_versions"`
	ModDependencies []string `json:"mod_dependencies" yaml:"mod_dependencies" toml:"mod_dependencies"`
	ModChangelog    string   `json:"mod_changelog" yaml:"mod_changelog" toml:"mod_changelog"`
}

// defaultConfigTemplate returns a ConfigTemplate with default values
func defaultConfigTemplate() ConfigTemplate {
	return ConfigTemplate{
		APIKey:          "",
		ModID:           0,
		ModName:         "",
		ModVersion:      "1.0.0",
		ModDescription:  "",
		ModAuthor:       "",
		ModLicense:      "",
		ModURL:          "",
		ModLogo:         "",
		ModCategories:   []string{},
		ModGameVersions: []string{},
		ModDependencies: []string{},
		ModChangelog:    "",
	}
}

// WriteTOMLTemplate writes a TOML config template
func WriteTOMLTemplate(filename string) error {
	data, err := toml.Marshal(defaultConfigTemplate())
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// WriteJSONTemplate writes a JSON config template
func WriteJSONTemplate(filename string) error {
	data, err := json.MarshalIndent(defaultConfigTemplate(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// WriteYAMLTemplate writes a YAML config template
func WriteYAMLTemplate(filename string) error {
	data, err := yaml.Marshal(defaultConfigTemplate())
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// WriteDotenvTemplate writes a .env template file based on the ConfigTemplate struct
func WriteDotenvTemplate(filename string) error {
	t := defaultConfigTemplate()
	v := reflect.ValueOf(t)
	typeOfT := v.Type()
	var sb strings.Builder

	for i := 0; i < v.NumField(); i++ {
		field := typeOfT.Field(i)
		key := strings.ToUpper(field.Tag.Get("json"))
		if key == "" {
			key = strings.ToUpper(field.Name)
		}
		val := v.Field(i)

		switch val.Kind() {
		case reflect.Slice:
			sb.WriteString(fmt.Sprintf("%s=\n", key))
		default:
			sb.WriteString(fmt.Sprintf("%s=%s\n", key, toString(val.Interface())))
		}
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

// toString converts supported value types to string
func toString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int, int64, int32:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// WriteTemplate writes a config template in the given format
// Supported formats: "toml", "json", "yaml", "env"
func WriteTemplate(format, filename string) error {
	switch strings.ToLower(format) {
	case "toml":
		return WriteTOMLTemplate(filename)
	case "json":
		return WriteJSONTemplate(filename)
	case "yaml", "yml":
		return WriteYAMLTemplate(filename)
	case "env", "dotenv":
		return WriteDotenvTemplate(filename)
	default:
		return fmt.Errorf("unsupported config format: %s", format)
	}
}
