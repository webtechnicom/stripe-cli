package profile

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var runtimeViper = viper.New()

// ConfigureProfile creates a profile when logging in
func (p *Profile) ConfigureProfile(apiKey string) error {

	err := p.setFilePath()
	if err != nil {
		return err
	}

	if p.DeviceName == "" {
		deviceName, err := os.Hostname()
		if err != nil {
			deviceName = "unknown"
		}
		p.DeviceName = deviceName
	}

	runtimeViper.Set(p.ProfileName +".device_name", strings.TrimSpace(p.DeviceName))
	runtimeViper.Set(p.ProfileName +".secret_key", strings.TrimSpace(apiKey))

	writeErr := writeConfig("secret_key")
	if writeErr != nil {
		return writeErr
	}

	fmt.Println("You're configured and all set to get started")

	return nil
}

// Temporary workaround until https://github.com/spf13/viper/pull/519 can remove a key from viper
func writeConfig(key string) error {
	configMap := viper.AllSettings()
	delete(configMap, key)
	buf := new(bytes.Buffer)
	encodeErr := toml.NewEncoder(buf).Encode(configMap)
	if encodeErr != nil {
		return encodeErr
	}
	err := runtimeViper.ReadConfig(buf)
	if err != nil {
		return err
	}

	runtimeViper.MergeInConfig()
	runtimeViper.WriteConfig()

	return nil

}

func (p *Profile) setFilePath() error {
	configPath := p.GetConfigFolder(os.Getenv("XDG_CONFIG_HOME"))
	dotfilePath := filepath.Join(configPath, "config.toml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	runtimeViper.SetConfigFile(dotfilePath)
	runtimeViper.SetConfigType("toml")

	return nil
}
