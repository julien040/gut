package profile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/julien040/gut/src/print"
	nanoid "github.com/matoous/go-nanoid/v2"

	keyring "github.com/99designs/keyring"
	config "github.com/gookit/config/v2"
	toml "github.com/gookit/config/v2/toml"
)

type Profile struct {
	Id       string
	Alias    string
	Username string
	Password string
	Website  string
}

type DiskProfile struct {
	Alias    string
	Website  string
	Username string
}

var configPath string

var profiles []Profile

var ring keyring.Keyring

func exit(err error, message string) {
	print.Message(message, "error")
	fmt.Println(err)
	os.Exit(1)
}

// Init a config file for the profiles and load it into the config package
func init() {
	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		exit(err, "We can't find your home directory 😓")
	}
	// Path to the config file
	configPath = filepath.Join(home, "/.gut/", "profiles.toml")

	// Check if .gut directory exists
	if _, err := os.Stat(filepath.Join(home, "/.gut/")); os.IsNotExist(err) {
		// Create .gut directory
		err = os.Mkdir(filepath.Join(home, "/.gut/"), 0755)
		if err != nil {
			exit(err, "We can't create the .gut directory 😓 at "+filepath.Join(home, "/.gut/"))
		}
	}

	// Create config file if it doesn't exist
	f, err := os.Open(configPath)
	if os.IsNotExist(err) {
		f, err = os.Create(configPath)
		if err != nil {
			exit(err, "We can't create the config file located at "+configPath+" 😓")
		}
	} else if err != nil {
		exit(err, "We can't open the config file located at "+configPath+" 😓")
	}
	f.Close()
	config.BindStruct("profile", &Profile{})

	config.AddDriver(toml.Driver)

	// Load config file
	err = config.LoadFiles(configPath)
	if err != nil {
		exit(err, "We can't load the config file located at "+configPath+" 😓")
	}

	// Load keyring
	ring, err = keyring.Open(keyring.Config{
		ServiceName: "gut",
	})
	if err != nil {
		exit(err, "We can open the keyring 😓")

	}

	// Load profiles
	data := config.Data()
	for key, val := range data {
		// Get password from keyring
		password, err := ring.Get(key)
		if err != nil {
			print.Message("The profile "+key+" doesn't have a password, we will skip it", print.Warning)
			continue
		}
		val := val.(map[string]interface{})
		alias, ok := val["Alias"].(string)
		if !ok {
			print.Message("The profile "+key+" doesn't have an alias, we will skip it", print.Warning)
			continue
		}
		website, ok := val["Website"].(string)
		if !ok {
			print.Message("The profile "+key+" doesn't have a website, we will skip it", print.Warning)
			continue
		}
		username, ok := val["Username"].(string)
		if !ok {
			print.Message("The profile "+key+" doesn't have a username, we will skip it", print.Warning)
			continue
		}

		// Add profile to the profiles array
		profiles = append(profiles, Profile{
			Id:       key,
			Alias:    alias,
			Username: username,
			Password: string(password.Data),
			Website:  website,
		})
	}

}

// Save the config file
func saveFile() {
	// Open config file in write mode
	f, err := os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		exit(err, "We can't open the config file located at "+configPath+" 😓")
	}
	// Save config file
	_, err = config.DumpTo(f, "toml")
	if err != nil {
		exit(err, "We can't save the config file located at "+configPath+" 😓")
	}
}

// Add a profile to the config file
func AddProfile(profile Profile) {
	id, err := nanoid.New()
	if err != nil {
		exit(err, "We can't generate an id for the profile 😓")
	}

	toSave := DiskProfile{
		Alias:    profile.Alias,
		Website:  profile.Website,
		Username: profile.Username,
	}
	err = ring.Set(keyring.Item{
		Key:  id,
		Data: []byte(profile.Password),
	})
	if err != nil {
		exit(err, "We can't save the password 😓")
	}

	// Add profile to the database
	err = config.Set(id, toSave)
	if err != nil {
		exit(err, "We can't save the profile in profiles.toml 😓")
	}
	saveFile()
}

// Return the profiles array
func GetProfiles() *[]Profile {
	return &profiles
}

func RemoveProfile(id string) {
	// Remove profile from the database
	err := config.Set(id, nil)
	if err != nil {
		exit(err, "We can't remove the profile from profiles.toml 😓")
	}
	// Remove password from the keyring
	err = ring.Remove(id)
	if err != nil {
		exit(err, "We can't remove the password from the keyring 😓")
	}
	saveFile()
}
