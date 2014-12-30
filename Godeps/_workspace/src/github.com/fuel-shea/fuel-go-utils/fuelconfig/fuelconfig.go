package fuelconfig

import (
	"bitbucket.org/kardianos/osext"
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	DBHost string
	DBName string
}

func CreateConfig(appName string) (Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "DEVELOPMENT"
	}
	configFilename := appName + "." + env + ".config.json"

	execDir, err := osext.ExecutableFolder()
	if err != nil {
		return Config{}, err
	}
	configPath := path.Join(execDir, configFilename)

	configFile, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	decoder := json.NewDecoder(configFile)

	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
