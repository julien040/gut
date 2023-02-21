package profile

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"github.com/julien040/gut/src/print"
	nanoid "github.com/matoous/go-nanoid/v2"

	keyringLinux "github.com/99designs/keyring"
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
	print.Message(message, print.Error)
	fmt.Println(err)
	os.Exit(1)
}

const serviceName = "gut"

var ring keyringLinux.Keyring

// Init a config file for the profiles and load it into the config package
func init() {
	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		exit(err, "Sorry, I can't find your home directory 😓")
	}
	// Path to the config file
	configPath = filepath.Join(home, "/.gut/", "profiles.toml")

	// Init keyring
	ring, _ = keyringLinux.Open(keyringLinux.Config{
		ServiceName:     serviceName,
		AllowedBackends: []keyringLinux.BackendType{keyringLinux.PassBackend},
	})

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
		var password string
		var err error

		password, err = retrievePassword(key)
		if err != nil {
			print.Message("I can't retrieve the password for the profile %s, I'll skip it", print.Warning, key)
			print.Message("Error: %s", print.Error, err.Error())
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
			Password: password,
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

// Check if an executable is in the path
//
// Intended to check if a command is available
func isExecInPath(executable string) bool {
	_, err := exec.LookPath(executable)
	return err == nil
}

// Check if the user is using a GUI on Linux. Returns true if yes.
//
// We have to check because gnome-keyring doesn't work on a server (requires to fill in a popup)
func checkGUIOnLinux() bool {
	// Try to find the Xorg executable
	// https://unix.stackexchange.com/a/237750 CC BY-SA 3.0

	// We use lookPath rather than "type Xorg" because I think type is a shell builtin
	// Go'll try to find the executable in the path to run it

	_, err := exec.LookPath("Xorg")
	fmt.Println(err)
	return err == nil

}

// Check if a folder exists
func checkFolderExists(path string) bool {
	// https://gist.github.com/mattes/d13e273314c3b3ade33f
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true

}

// Save the password in the keyring
//
// If the user is using a GUI on Linux, we use the default gnome-keyring.
// If not, we use pass.
// On other OS, we use the default keyring defined by zalando/go-keyring
func savePassword(id string, password string) {
	// https://github.com/julien040/gut/issues/14
	if runtime.GOOS == "linux" {
		// Check if the user is using a GUI
		// If yes, we use the default gnome-keyring
		guiAvailable := checkGUIOnLinux()
		if guiAvailable {
			// Check if gnome-keyring is installed
			if !isExecInPath("gnome-keyring-daemon") {
				print.Message("To install gnome-keyring-daemon, run: sudo apt install gnome-keyring", print.None)
				print.Message("If you use YUM, run: sudo yum install gnome-keyring", print.None)
				exit(nil, "I can't find gnome-keyring-daemon in your path. Please install it 😓")
			}

			// We use the default gnome-keyring
			err := keyring.Set(serviceName, id, password)
			if err != nil {
				exit(err, "I can't save the password in the keyring 😓")
			}

			// If not, we use pass.
		} else {
			// Check if pass is installed
			if isExecInPath("pass") {
				// We check if the password store exists
				homeDir, err := os.UserHomeDir()
				if err != nil {
					exit(err, "I can't get your home directory 😓")
				}
				dirPassFolder := path.Join(homeDir, ".password-store")
				if !checkFolderExists(dirPassFolder) {
					// We prompt the user to create the password store
					print.Message("Please set up a password store with pass. To do so, follow this guide: https://gut-cli.dev/error/setup-pass-store", print.None)

					os.Exit(1)
				} else {
					// We save the password in the password store
					err := ring.Set(keyringLinux.Item{
						Key:   id,
						Data:  []byte(password),
						Label: "Password for " + id,
					})
					if err != nil {
						exit(err, "I can't save the password in the keyring 😓")
					}

				}

				// If not, we explain to the user how to install it on his distro
			} else {
				exit(errors.New("pass not installed"), "Please install pass with your package manager (https://www.passwordstore.org/#download)")
			}
		}

	} else {
		err := keyring.Set(serviceName, id, password)
		if err != nil {
			exit(err, "I can't save the password in the keyring 😓")
		}
	}
}

func retrievePassword(id string) (string, error) {
	// https://github.com/julien040/gut/issues/14
	if runtime.GOOS == "linux" {
		// Check if the user is using a GUI
		// If yes, we use the default gnome-keyring
		guiAvailable := checkGUIOnLinux()
		if guiAvailable {
			// We use the default gnome-keyring
			password, err := keyring.Get(serviceName, id)
			return password, err
		} else {
			// If not, we use pass.
			// Check if pass is installed
			if isExecInPath("pass") {
				// We check if the password store exists
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return "", errors.New("unable to get your home directory")
				}
				dirPassFolder := path.Join(homeDir, ".password-store")
				if !checkFolderExists(dirPassFolder) {
					// We prompt the user to create the password store
					print.Message("Please set up a password store with pass. To do so, follow this guide: https://gut-cli.dev/error/setup-pass-store", print.None)
					return "", errors.New("pass is not set up")
				} else {
					// We retrieve the password from the password store
					password, err := ring.Get(id)
					if err != nil {
						print.Message("Unlock your password store with pass first", print.Info)
						print.Message("You can do it with the following command:", print.Info)
						print.Message("	pass show "+id, print.None)
						print.Message("To learn more about this, follow this guide: https://gut-cli.dev/error/unlock-pass-store", print.None)
						return "", errors.New("unable to retrieve the password from the keyring")
					}
					return string(password.Data), nil
				}
			} else {
				print.Message("To install pass, follow this guide: https://www.passwordstore.org/#download", print.None)
				return "", errors.New("pass not installed")
			}
		}

	} else {
		password, err := keyring.Get(serviceName, id)
		if err != nil {
			return "", err
		}
		return password, nil
	}
}

func deletePassword(id string) error {
	if runtime.GOOS == "linux" {
		// Check if the user is using a GUI
		// If yes, we use the default gnome-keyring
		guiAvailable := checkGUIOnLinux()
		if guiAvailable {
			// We use the default gnome-keyring
			return keyring.Delete(serviceName, id)
		} else {
			// If not, we use pass.
			// Check if pass is installed
			if isExecInPath("pass") {
				// We check if the password store exists
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return errors.New("unable to get your home directory")
				}
				dirPassFolder := path.Join(homeDir, ".password-store")
				if !checkFolderExists(dirPassFolder) {
					// We prompt the user to create the password store
					return errors.New("pass is not set up")
				} else {
					// We delete the password from the password store
					return ring.Remove(id)
				}
			} else {
				fmt.Println("To install pass, follow this guide: https://www.passwordstore.org/#download")
				return errors.New("pass not installed")
			}
		}

	} else {
		return keyring.Delete(serviceName, id)
	}
}

// Add a profile to the config file and return the id
func AddProfile(profile Profile) string {
	id, err := nanoid.New()
	if err != nil {
		exit(err, "Sorry, I can't generate an id 😓")
	}

	// Save password in the keyring
	savePassword(id, profile.Password)

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
	err := deletePassword(id)

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
