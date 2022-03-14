package util

import (
	"encoding/json"
	"os"
)

type Config struct {
	Email struct {
		From string `json:"from"`
		To string `json:"to"`
		SmtpPassword string `json:"smtpPassword"`
		Mailserver string `json:"mailserver"`
		Mailport string `json:"mailport"`
	} `json:"email"`

	Api struct {
		Address string `json:"address"`
		AuthToken string `json:"authtoken"`
	}
}

func LoadConfiguration(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)

	if err != nil {
		return config, err
	}
	
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}