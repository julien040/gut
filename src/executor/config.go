package executor

import (
	"errors"

	"github.com/go-git/go-git/v5/config"
)

// Get a config flag from the specified scope
//
// Return "" if the flag is not found
func getConfigFlagForGivenScope(scope config.Scope, section string, flag string) (string, error) {
	conf, err := config.LoadConfig(scope)
	if err != nil {
		return "", err
	}

	if !conf.Raw.HasSection(section) {
		return "", errors.New("section not found")
	}

	if !conf.Raw.Section(section).HasOption(flag) {
		return "", errors.New("flag not found")
	}

	return conf.Raw.Section(section).Option(flag), nil

}

// Get a config flag from the global and system scope
// in this order
// Return "" if the flag is not found
func getConfigFlag(section string, flag string) (string, error) {
	res, err := getConfigFlagForGivenScope(config.GlobalScope, section, flag)
	if err != nil {
		res, err = getConfigFlagForGivenScope(config.SystemScope, section, flag)
		if err != nil {
			return "", err
		}
		return res, nil
	}
	return res, nil

}

// Get the set editor
func GetConfigEditor() (string, error) {
	return getConfigFlag("core", "editor")
}
