package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/julien040/gut/src/print"
	nanoid "github.com/matoous/go-nanoid/v2"

	"github.com/BurntSushi/toml"
	keyring "github.com/zalando/go-keyring"
)

type Profile struct {
	Id       string
	Alias    string
	Username string
	Password string
	Website  string
	Email    string
}

type DiskProfile struct {
	Alias    string
	Website  string
	Username string
	Email    string
}

var configPath string

var profiles []Profile

func exit(err error, message string) {
	print.Message(message, "error")
	fmt.Println(err)
	os.Exit(1)
}

const serviceName = "gut"

// Init a config file for the profiles and load it into the config package
func init() {
	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		exit(err, "Sorry, I can't find your home directory 😓")
	}
	// Path to the config file
	configPath = filepath.Join(home, "/.gut/", "profiles.toml")

	// Check if .gut directory exists
	if _, err := os.Stat(filepath.Join(home, "/.gut/")); os.IsNotExist(err) {
		// Create .gut directory
		err = os.Mkdir(filepath.Join(home, "/.gut/"), 0755)
		if err != nil {
			exit(err, "I can't create the .gut directory 😓 at "+filepath.Join(home, "/.gut/"))
		}
	}

	// Create config file if it doesn't exist
	f, err := os.Open(configPath)
	if os.IsNotExist(err) {
		f, err = os.Create(configPath)
		if err != nil {
			exit(err, "I can't create the config file 😓 at "+configPath)
		}
	} else if err != nil {
		exit(err, "I can't open the config file 😓 at "+configPath)
	}
	f.Close()
	/* 	config.BindStruct("profile", &Profile{})

	   	config.AddDriver(toml.Driver)

	   	// Load config file
	   	err = config.LoadFiles(configPath)
	   	if err != nil {
	   		exit(err, "I can't load the config file 😓 at "+configPath)
	   	} */

	// Open file in read mode
	f, err = os.Open(configPath)
	if err != nil {
		exit(err, "I can't open the config file 😓 at "+configPath)
	}
	// Load config file and unmarshal it
	var data map[string]interface{}
	_, err = toml.NewDecoder(f).Decode(&data)
	if err != nil {
		exit(err, "I can't load the config file 😓 at "+configPath)
	}

	for key, val := range data {
		// Get password from keyring
		password, err := keyring.Get(serviceName, key)
		if err != nil {
			print.Message("The profile "+key+" doesn't have a password, I'll skip it", print.Warning)
			continue
		}
		val := val.(map[string]interface{})
		alias, ok := val["Alias"].(string)
		if !ok {
			print.Message("The profile "+key+" doesn't have an alias, I'll skip it", print.Warning)
			continue
		}
		website, ok := val["Website"].(string)
		if !ok {
			print.Message("The profile "+key+" doesn't have a website, I'll skip it", print.Warning)
			continue
		}
		username, ok := val["Username"].(string)
		if !ok {
			print.Message("The profile "+key+" doesn't have a username, I'll skip it", print.Warning)
			continue
		}
		email, ok := val["Email"].(string)
		if !ok {
			print.Message("The profile "+key+" doesn't have an email, I'll skip it", print.Warning)
			continue
		}

		// Add profile to the profiles array
		profiles = append(profiles, Profile{
			Id:       key,
			Alias:    alias,
			Username: username,
			Password: string(password),
			Website:  website,
			Email:    email,
		})
	}

}

// Save the config file
func saveFile() {
	// Open config file in write mode
	f, err := os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		exit(err, "I can't open the config file located at "+configPath+" 😓")
	}
	/* // Save config file
	_, err = config.DumpTo(f, "toml")
	if err != nil {
		exit(err, "I can't save the config file located at "+configPath+" 😓")
	} */
	// Create a new encoder
	encoder := toml.NewEncoder(f)

	// Convert profiles to a map of DiskProfile
	profilesMap := make(map[string]DiskProfile)
	for _, profile := range profiles {
		profilesMap[profile.Id] = DiskProfile{
			Alias:    profile.Alias,
			Website:  profile.Website,
			Username: profile.Username,
			Email:    profile.Email,
		}
	}
	// Encode the map
	err = encoder.Encode(profilesMap)
	if err != nil {
		exit(err, "I can't save the config file located at "+configPath+" 😓")
	}

}

// Add a profile to the config file and return the id
func AddProfile(profile Profile) string {
	id, err := nanoid.New()
	if err != nil {
		exit(err, "Sorry, I can't generate an id 😓")
	}

	err = keyring.Set(serviceName, id, profile.Password)
	if err != nil {
		exit(err, "Sorry, I can't save the password in the keyring 😓")
	}

	// When the function is called, the id is empty
	// So we set it here
	profile.Id = id
	profiles = append(profiles, profile)

	saveFile()
	return id
}

// Return the profiles array
func GetProfiles() *[]Profile {
	return &profiles
}

func RemoveProfile(id string) {
	// Remove profile from the database
	for i, profile := range profiles {
		if profile.Id == id {
			profiles = append(profiles[:i], profiles[i+1:]...)
			break
		}
	}
	// Remove password from the keyring
	err := keyring.Delete(serviceName, id)
	if err != nil {
		exit(err, "Sorry, I can't remove the password from the keyring 😓")
	}
	saveFile()
}

func CheckIfProfileExists(id string) bool {
	for _, profile := range profiles {
		if profile.Id == id {
			return true
		}
	}
	return false
}

func GetProfileIDFromPath(path string) string {
	// Open file in read mode
	f, err :=
		os.OpenFile(filepath.Join(path, ".gut"),
			os.O_RDONLY, 0755)

	if err != nil {
		defer f.Close()
		return ""

	} else {
		defer f.Close()
		// Close file at the end of the function

		// Create the schema
		profileIDSchema := SchemaGutConf{}

		// Decode ID in TOML
		t := toml.NewDecoder(f)
		_, err = t.Decode(&profileIDSchema)
		if err != nil {
			print.Message("Sorry, I can't read the .gut file 😓", print.Error)
			os.Exit(1)
		}
		return profileIDSchema.ProfileID
	}

}

func GetProfileFromPath(path string) (Profile, error) {
	id := GetProfileIDFromPath(path)
	if id == "" {
		return Profile{}, errors.New("no profile found in this directory")
	}
	for _, profile := range profiles {
		if profile.Id == id {
			return profile, nil
		}
	}
	return Profile{}, errors.New("no profile found globally")
}

func IsAliasAlreadyUsed(alias string) bool {
	for _, profile := range profiles {
		if profile.Alias == alias {
			return true
		}
	}
	return false

}

type SchemaGutConf struct {
	ProfileID string `toml:"profile_id"`
	UpdatedAt string `toml:"updated_at"`
}
