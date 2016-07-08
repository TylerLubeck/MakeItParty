package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
	"os"
)

var (
	build_time string
	version    string
)

type Configuration struct {
	APIKey     string
	TCUsername string
	TCPassword string
}

func get_config(file_name string) (*Configuration, error) {
	file, e := os.Open(file_name)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	configuration := Configuration{}
	err := json.NewDecoder(file).Decode(&configuration)

	return &configuration, err
}

func main() {

	logFile, err := os.OpenFile("./makeitparty.log", os.O_WRONLY, 0666)

	if err != nil {
		panic(err)
	}

	defer logFile.Close()

	log.SetLevel(log.DebugLevel)
	log.SetOutput(logFile)
	log.SetOutput(os.Stderr)

	log.WithFields(log.Fields{
		"version":    version,
		"build_time": build_time,
	}).Info("MakeItParty is starting up...")
	config, err := get_config("api_keys.json")

	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		panic(err)
	}

	s := slack.New(config.APIKey)

	HandleSlackEvents(s)
}
