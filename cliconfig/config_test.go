package cliconfig

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kmova/waytogo/cliconfig/configfile"
	"github.com/kmova/waytogo/pkg/homedir"
)

func TestEmptyConfigDir(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	SetConfigDir(tmpHome)

	config, err := Load("")
	if err != nil {
		t.Fatalf("Failed loading on empty config dir: %q", err)
	}

	expectedConfigFilename := filepath.Join(tmpHome, ConfigFileName)
	if config.Filename != expectedConfigFilename {
		t.Fatalf("Expected config filename %s, got %s", expectedConfigFilename, config.Filename)
	}

	// Now save it and make sure it shows up in new form
	saveConfigAndValidateNewFormat(t, config, tmpHome)
}

func TestMissingFile(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	config, err := Load(tmpHome)
	if err != nil {
		t.Fatalf("Failed loading on missing file: %q", err)
	}

	// Now save it and make sure it shows up in new form
	saveConfigAndValidateNewFormat(t, config, tmpHome)
}

func TestSaveFileToDirs(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	tmpHome += "/.waytogo"

	config, err := Load(tmpHome)
	if err != nil {
		t.Fatalf("Failed loading on missing file: %q", err)
	}

	// Now save it and make sure it shows up in new form
	saveConfigAndValidateNewFormat(t, config, tmpHome)
}

func TestEmptyFile(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	if err := ioutil.WriteFile(fn, []byte(""), 0600); err != nil {
		t.Fatal(err)
	}

	_, err = Load(tmpHome)
	if err == nil {
		t.Fatalf("Was supposed to fail")
	}
}

func TestEmptyJson(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	if err := ioutil.WriteFile(fn, []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	if err != nil {
		t.Fatalf("Failed loading on empty json file: %q", err)
	}

	// Now save it and make sure it shows up in new form
	saveConfigAndValidateNewFormat(t, config, tmpHome)
}

func TestOldInvalidsAuth(t *testing.T) {
	invalids := map[string]string{
		`username = test`: "The Auth config file is empty",
		`username
password`: "Invalid Auth config file",
		`username = test
email`: "Invalid auth configuration file",
	}

	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	homeKey := homedir.Key()
	homeVal := homedir.Get()

	defer func() { os.Setenv(homeKey, homeVal) }()
	os.Setenv(homeKey, tmpHome)

	for content, expectedError := range invalids {
		fn := filepath.Join(tmpHome, oldConfigfile)
		if err := ioutil.WriteFile(fn, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}

		config, err := Load(tmpHome)
		// Use Contains instead of == since the file name will change each time
		if err == nil || !strings.Contains(err.Error(), expectedError) {
			t.Fatalf("Should have failed\nConfig: %v\nGot: %v\nExpected: %v", config, err, expectedError)
		}

	}
}

func TestOldValidAuth(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	homeKey := homedir.Key()
	homeVal := homedir.Get()

	defer func() { os.Setenv(homeKey, homeVal) }()
	os.Setenv(homeKey, tmpHome)

	fn := filepath.Join(tmpHome, oldConfigfile)
	js := `username = am9lam9lOmhlbGxv
	email = user@example.com`
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	if err != nil {
		t.Fatal(err)
	}

	// defaultIndexserver is https://index.waytogo.io/v1/
	ac := config.AuthConfigs["https://index.waytogo.io/v1/"]
	if ac.Username != "joejoe" || ac.Password != "waytogo" {
		t.Fatalf("Missing data from parsing:\n%q", config)
	}

	// Now save it and make sure it shows up in new form
	configStr := saveConfigAndValidateNewFormat(t, config, tmpHome)

	expConfStr := `{
	"auths": {
		"https://index.waytogo.io/v1/": {
			"auth": "am9lam9lOmhlbGxv"
		}
	}
}`

	if configStr != expConfStr {
		t.Fatalf("Should have save in new form: \n%s\n not \n%s", configStr, expConfStr)
	}
}

func TestOldJsonInvalid(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	homeKey := homedir.Key()
	homeVal := homedir.Get()

	defer func() { os.Setenv(homeKey, homeVal) }()
	os.Setenv(homeKey, tmpHome)

	fn := filepath.Join(tmpHome, oldConfigfile)
	js := `{"https://index.waytogo.io/v1/":{"auth":"test","email":"user@example.com"}}`
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	// Use Contains instead of == since the file name will change each time
	if err == nil || !strings.Contains(err.Error(), "Invalid auth configuration file") {
		t.Fatalf("Expected an error got : %v, %v", config, err)
	}
}

func TestOldJson(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	homeKey := homedir.Key()
	homeVal := homedir.Get()

	defer func() { os.Setenv(homeKey, homeVal) }()
	os.Setenv(homeKey, tmpHome)

	fn := filepath.Join(tmpHome, oldConfigfile)
	js := `{"https://index.waytogo.io/v1/":{"auth":"am9lam9lOmhlbGxv","email":"user@example.com"}}`
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	if err != nil {
		t.Fatalf("Failed loading on empty json file: %q", err)
	}

	ac := config.AuthConfigs["https://index.waytogo.io/v1/"]
	if ac.Username != "joejoe" || ac.Password != "waytogo" {
		t.Fatalf("Missing data from parsing:\n%q", config)
	}

	// Now save it and make sure it shows up in new form
	configStr := saveConfigAndValidateNewFormat(t, config, tmpHome)

	expConfStr := `{
	"auths": {
		"https://index.waytogo.io/v1/": {
			"auth": "am9lam9lOmhlbGxv",
			"email": "user@example.com"
		}
	}
}`

	if configStr != expConfStr {
		t.Fatalf("Should have save in new form: \n'%s'\n not \n'%s'\n", configStr, expConfStr)
	}
}

func TestNewJson(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	js := ` { "auths": { "https://index.waytogo.io/v1/": { "auth": "am9lam9lOmhlbGxv" } } }`
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	if err != nil {
		t.Fatalf("Failed loading on empty json file: %q", err)
	}

	ac := config.AuthConfigs["https://index.waytogo.io/v1/"]
	if ac.Username != "joejoe" || ac.Password != "waytogo" {
		t.Fatalf("Missing data from parsing:\n%q", config)
	}

	// Now save it and make sure it shows up in new form
	configStr := saveConfigAndValidateNewFormat(t, config, tmpHome)

	expConfStr := `{
	"auths": {
		"https://index.waytogo.io/v1/": {
			"auth": "am9lam9lOmhlbGxv"
		}
	}
}`

	if configStr != expConfStr {
		t.Fatalf("Should have save in new form: \n%s\n not \n%s", configStr, expConfStr)
	}
}

func TestNewJsonNoEmail(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	js := ` { "auths": { "https://index.waytogo.io/v1/": { "auth": "am9lam9lOmhlbGxv" } } }`
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	if err != nil {
		t.Fatalf("Failed loading on empty json file: %q", err)
	}

	ac := config.AuthConfigs["https://index.waytogo.io/v1/"]
	if ac.Username != "joejoe" || ac.Password != "waytogo" {
		t.Fatalf("Missing data from parsing:\n%q", config)
	}

	// Now save it and make sure it shows up in new form
	configStr := saveConfigAndValidateNewFormat(t, config, tmpHome)

	expConfStr := `{
	"auths": {
		"https://index.waytogo.io/v1/": {
			"auth": "am9lam9lOmhlbGxv"
		}
	}
}`

	if configStr != expConfStr {
		t.Fatalf("Should have save in new form: \n%s\n not \n%s", configStr, expConfStr)
	}
}

func TestJsonWithPsFormat(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	js := `{
		"auths": { "https://index.waytogo.io/v1/": { "auth": "am9lam9lOmhlbGxv", "email": "user@example.com" } },
		"psFormat": "table {{.ID}}\\t{{.Label \"com.waytogo.label.cpu\"}}"
}`
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	if err != nil {
		t.Fatalf("Failed loading on empty json file: %q", err)
	}

	if config.PsFormat != `table {{.ID}}\t{{.Label "com.waytogo.label.cpu"}}` {
		t.Fatalf("Unknown ps format: %s\n", config.PsFormat)
	}

	// Now save it and make sure it shows up in new form
	configStr := saveConfigAndValidateNewFormat(t, config, tmpHome)
	if !strings.Contains(configStr, `"psFormat":`) ||
		!strings.Contains(configStr, "{{.ID}}") {
		t.Fatalf("Should have save in new form: %s", configStr)
	}
}

// Save it and make sure it shows up in new form
func saveConfigAndValidateNewFormat(t *testing.T, config *configfile.ConfigFile, homeFolder string) string {
	if err := config.Save(); err != nil {
		t.Fatalf("Failed to save: %q", err)
	}

	buf, err := ioutil.ReadFile(filepath.Join(homeFolder, ConfigFileName))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(buf), `"auths":`) {
		t.Fatalf("Should have save in new form: %s", string(buf))
	}
	return string(buf)
}

func TestConfigDir(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	if ConfigDir() == tmpHome {
		t.Fatalf("Expected ConfigDir to be different than %s by default, but was the same", tmpHome)
	}

	// Update configDir
	SetConfigDir(tmpHome)

	if ConfigDir() != tmpHome {
		t.Fatalf("Expected ConfigDir to %s, but was %s", tmpHome, ConfigDir())
	}
}

func TestConfigFile(t *testing.T) {
	configFilename := "configFilename"
	configFile := NewConfigFile(configFilename)

	if configFile.Filename != configFilename {
		t.Fatalf("Expected %s, got %s", configFilename, configFile.Filename)
	}
}

func TestJsonReaderNoFile(t *testing.T) {
	js := ` { "auths": { "https://index.waytogo.io/v1/": { "auth": "am9lam9lOmhlbGxv", "email": "user@example.com" } } }`

	config, err := LoadFromReader(strings.NewReader(js))
	if err != nil {
		t.Fatalf("Failed loading on empty json file: %q", err)
	}

	ac := config.AuthConfigs["https://index.waytogo.io/v1/"]
	if ac.Username != "joejoe" || ac.Password != "waytogo" {
		t.Fatalf("Missing data from parsing:\n%q", config)
	}

}

func TestOldJsonReaderNoFile(t *testing.T) {
	js := `{"https://index.waytogo.io/v1/":{"auth":"am9lam9lOmhlbGxv","email":"user@example.com"}}`

	config, err := LegacyLoadFromReader(strings.NewReader(js))
	if err != nil {
		t.Fatalf("Failed loading on empty json file: %q", err)
	}

	ac := config.AuthConfigs["https://index.waytogo.io/v1/"]
	if ac.Username != "joejoe" || ac.Password != "waytogo" {
		t.Fatalf("Missing data from parsing:\n%q", config)
	}
}

func TestJsonWithPsFormatNoFile(t *testing.T) {
	js := `{
		"auths": { "https://index.waytogo.io/v1/": { "auth": "am9lam9lOmhlbGxv", "email": "user@example.com" } },
		"psFormat": "table {{.ID}}\\t{{.Label \"com.waytogo.label.cpu\"}}"
}`
	config, err := LoadFromReader(strings.NewReader(js))
	if err != nil {
		t.Fatalf("Failed loading on empty json file: %q", err)
	}

	if config.PsFormat != `table {{.ID}}\t{{.Label "com.waytogo.label.cpu"}}` {
		t.Fatalf("Unknown ps format: %s\n", config.PsFormat)
	}

}

func TestJsonSaveWithNoFile(t *testing.T) {
	js := `{
		"auths": { "https://index.waytogo.io/v1/": { "auth": "am9lam9lOmhlbGxv" } },
		"psFormat": "table {{.ID}}\\t{{.Label \"com.waytogo.label.cpu\"}}"
}`
	config, err := LoadFromReader(strings.NewReader(js))
	err = config.Save()
	if err == nil {
		t.Fatalf("Expected error. File should not have been able to save with no file name.")
	}

	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create a temp dir: %q", err)
	}
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	f, _ := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	err = config.SaveToWriter(f)
	if err != nil {
		t.Fatalf("Failed saving to file: %q", err)
	}
	buf, err := ioutil.ReadFile(filepath.Join(tmpHome, ConfigFileName))
	if err != nil {
		t.Fatal(err)
	}
	expConfStr := `{
	"auths": {
		"https://index.waytogo.io/v1/": {
			"auth": "am9lam9lOmhlbGxv"
		}
	},
	"psFormat": "table {{.ID}}\\t{{.Label \"com.waytogo.label.cpu\"}}"
}`
	if string(buf) != expConfStr {
		t.Fatalf("Should have save in new form: \n%s\nnot \n%s", string(buf), expConfStr)
	}
}

func TestLegacyJsonSaveWithNoFile(t *testing.T) {

	js := `{"https://index.waytogo.io/v1/":{"auth":"am9lam9lOmhlbGxv","email":"user@example.com"}}`
	config, err := LegacyLoadFromReader(strings.NewReader(js))
	err = config.Save()
	if err == nil {
		t.Fatalf("Expected error. File should not have been able to save with no file name.")
	}

	tmpHome, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create a temp dir: %q", err)
	}
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	f, _ := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err = config.SaveToWriter(f); err != nil {
		t.Fatalf("Failed saving to file: %q", err)
	}
	buf, err := ioutil.ReadFile(filepath.Join(tmpHome, ConfigFileName))
	if err != nil {
		t.Fatal(err)
	}

	expConfStr := `{
	"auths": {
		"https://index.waytogo.io/v1/": {
			"auth": "am9lam9lOmhlbGxv",
			"email": "user@example.com"
		}
	}
}`

	if string(buf) != expConfStr {
		t.Fatalf("Should have save in new form: \n%s\n not \n%s", string(buf), expConfStr)
	}
}
