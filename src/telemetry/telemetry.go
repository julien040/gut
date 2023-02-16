package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/julien040/gut/src/print"

	"github.com/BurntSushi/toml"
)

var telemetryEnabled = false

var consentStateKnown = false

const eventServerURL = "https://api-events.gut-cli.dev/v1"

const gutVersion = "0.1.0"

var wg sync.WaitGroup

func exit(err error, message string) {
	print.Message(message, "error")
	fmt.Println(err)
	os.Exit(1)
}

func getConsentStateFromFile() {
	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		exit(err, "Sorry, I can't find your home directory 😓")
	}
	// Path to the config file
	consentPath := filepath.Join(home, "/.gut/", "consent.toml")

	// Check if .gut directory exists
	if _, err := os.Stat(filepath.Join(home, "/.gut/")); os.IsNotExist(err) {
		// Create .gut directory
		err = os.Mkdir(filepath.Join(home, "/.gut/"), 0755)
		if err != nil {
			exit(err, "I can't create the .gut directory 😓 at "+filepath.Join(home, "/.gut/"))
		}
	}

	// Check if consent file exists
	if _, err := os.Stat(consentPath); os.IsNotExist(err) {
		consentStateKnown = false
		return

	}

	// Read the consent file
	consentFile, err := os.Open(consentPath)
	if err != nil {
		exit(err, "I can't open the consent file 😓 at "+consentPath)
	}
	defer consentFile.Close()

	var consent struct {
		Telemetry bool
	}
	_, err = toml.NewDecoder(consentFile).Decode(&consent)
	if err != nil {
		print.Message("I can't read the consent file at %s. Please, can you check if I have the right permissions?", "error", consentPath)
		consentStateKnown = false
		return
	}

	telemetryEnabled = consent.Telemetry
	consentStateKnown = true
}

func SetConsentState(state bool) {
	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		exit(err, "Sorry, I can't find your home directory 😓")
	}
	// Path to the config file
	consentPath := filepath.Join(home, "/.gut/", "consent.toml")

	// Check if .gut directory exists
	if _, err := os.Stat(filepath.Join(home, "/.gut/")); os.IsNotExist(err) {
		// Create .gut directory
		err = os.Mkdir(filepath.Join(home, "/.gut/"), 0755)
		if err != nil {
			exit(err, "I can't create the .gut directory 😓 at "+filepath.Join(home, "/.gut/"))
		}
	}

	// Check if consent file exists
	if _, err := os.Stat(consentPath); os.IsNotExist(err) {
		// Create the consent file
		consentFile, err := os.Create(consentPath)
		if err != nil {
			exit(err, "I can't create the consent file 😓 at "+consentPath)
		}
		defer consentFile.Close()
	}

	// Open the consent file
	consentFile, err := os.OpenFile(consentPath, os.O_RDWR, 0755)
	if err != nil {
		exit(err, "I can't open the consent file 😓 at "+consentPath)
	}
	defer consentFile.Close()

	/*
		I don't know why, but I can't use writeString() to write to the file.
		It would write Telemetry = truee with two e at the end.
		I can't find why, so I use a struct to write to the file.
	*/
	type ToWrite struct {
		Telemetry bool
	}

	err = toml.NewEncoder(consentFile).Encode(ToWrite{Telemetry: state})
	if err != nil {
		exit(err, "I can't write to the consent file 😓 at "+consentPath)
	}

	// Update the telemetryEnabled variable
	telemetryEnabled = state
	consentStateKnown = true
}

type Event struct {
	EventName  string `json:"event_name"`
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	CPUCores   int    `json:"cpu_cores"`
	GutVersion string `json:"gut_version"`
}

func GetConsentState() bool {
	if !consentStateKnown {
		getConsentStateFromFile()
	}
	return telemetryEnabled
}

func IsConsentStateKnown() bool {
	return consentStateKnown
}

func LogCommand(command string) {
	if !consentStateKnown {
		return
	}
	if !telemetryEnabled {
		return
	}

	event := Event{
		EventName:  command,
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		CPUCores:   runtime.NumCPU(),
		GutVersion: gutVersion,
	}
	wg.Add(1)

	go sendEvent(event)

	wg.Wait()
}

func sendEvent(event Event) {
	defer wg.Done()
	client := &http.Client{}
	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", eventServerURL, io.NopCloser(bytes.NewReader(jsonEvent)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		return
	}

}

func init() {
	getConsentStateFromFile()

}
