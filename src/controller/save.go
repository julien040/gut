package controller

import (
	"errors"
	"fmt"

	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/spf13/cobra"
)

type emoji struct {
	Emoji       string
	Code        string
	Description string
}

// Some emojis are missing.
// Should be fixed later
var gitEmoji = []emoji{
	{"🎉", ":tada:", "Initial commit"},
	{"✨", ":sparkles:", "Introduce new features"},
	{"🐛", ":bug:", "Fix a bug"},
	{"🔥", ":fire:", "Remove code or files"},
	{"🚑", ":ambulance:", "Critical hotfix"},
	{"📝", ":memo:", "Add or update documentation"},
	{"🎨", ":art:", "Improve the format/structure of the code"},
	{"⚡️", ":zap:", "Improve performance"},
	{"🔒", ":lock:", "Fix security issues"},
	{"🔖", ":bookmark:", "Release / Version tags"},
	{"🚀", ":rocket:", "Deploy stuff"},
	{"💄", ":lipstick:", "Add or update the UI and style files"},
	{"🎨", ":art:", "Improve structure / format of the code"},
	{"🚧", ":construction:", "Work in progress"},
	{"💚", ":green_heart:", "Fix CI Build"},
	{"⬇️", ":arrow_down:", "Downgrade dependencies"},
	{"⬆️", ":arrow_up:", "Upgrade dependencies"},
	{"📌", ":pushpin:", "Pin dependencies to specific versions"},
	{"♻️", ":recycle:", "Refactor code"},
	{"➖", ":heavy_minus_sign:", "Remove a dependency"},
	{"🐧", ":penguin:", "Fix something on Linux"},
	{"🍎", ":apple:", "Fix something on macOS"},
	{"🏁", ":checkered_flag:", "Fix something on Windows"},
	{"🤖", ":robot:", "Fix something on Android"},
	{"🍏", ":green_apple:", "Fix something on iOS"},
	{"🔧", ":wrench:", "Add or update configuration files"},
	{"🌐", ":globe_with_meridians:", "Internationalization and localization"},
	{"✏️", ":pencil2:", "Fix typos"},
	{"💩", ":poop:", "Write bad code that needs to be improved"},
	{"⏪", ":rewind:", "Revert changes"},
	{"🔀", ":twisted_rightwards_arrows:", "Merge branches"},
	{"📦", ":package:", "Add or update compiled files or packages"},
	{"👽", ":alien:", "Update code due to external API changes"},
	{"🚚", ":truck:", "Move or rename resources (e.g.: files, paths, routes)"},
	{"📄", ":page_facing_up:", "Add or update license"},
	{"💥", ":boom:", "Introduce breaking changes"},
	{"🍱", ":bento:", "Add or update assets"},
	{"♿️", ":wheelchair:", "Improve accessibility"},
	{"💡", ":bulb:", "Add or update comments in source code"},
	{"🍻", ":beers:", "Write code drunkenly"},
	{"💬", ":speech_balloon:", "Add or update text and literals"},
	{"🗃", ":card_file_box:", "Perform database related changes"},
	{"🔊", ":loud_sound:", "Add or update logs"},
	{"🔇", ":mute:", "Remove logs"},
	{"👥", ":busts_in_silhouette:", "Add or update contributor(s)"},
	{"🚸", ":children_crossing:", "Improve user experience / usability"},
	{"🏗", ":building_construction:", "Make architectural changes"},
	{"📱", ":iphone:", "Work on responsive design"},
	{"🤡", ":clown_face:", "Mock things"},
	{"🥚", ":egg:", "Add or update an easter egg"},
	{"🙈", ":see_no_evil:", "Add or update a .gitignore file"},
	{"📸", ":camera_flash:", "Add or update snapshots"},
	{"⚗", ":alembic:", "Experiment new things"},
	{"🔍", ":mag:", "Improve SEO"},
	{"🏷", ":label:", "Add or update types (Flow, TypeScript)"},
	{"🌱", ":seedling:", "Add or update seed files"},
	{"🚩", ":triangular_flag_on_post:", "Add, update, or remove feature flags"},
	{"🥅", ":goal_net:", "Catch errors"},
	{"💫", ":dizzy:", "Add or update animations and transitions"},
	{"🗑", ":wastebasket:", "Deprecate code that needs to be cleaned up"},
	{"🛂", ":passport_control:", "Work on code related to authorization"},
	{"🩹", ":adhesive_bandage:", "Simple fix for a non-critical issue"},
	{"🔖", ":bookmark_tabs:", "Release / Version tags"},
	{"👷", ":construction_worker:", "Add or update CI build system"},
	{"💸", ":moneybag:", "Add or update financial, legal, or business documentation"},
	{"📦", ":package:", "Add or update compiled files or packages"},
	{"🦺", ":safety_vest:", "Add or update security"},
	{"📈", ":chart_with_upwards_trend:", "Add or update analytics or track code"},
}

func Save(cmd *cobra.Command, args []string) {
	// Get where the command has been called
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't get the current working directory", err)
	}

	// Check if the current directory is a git repository
	checkIfGitRepoInitialized(wd)

	// Check if the user config is set
	verifUserConfig(wd)

	// Get the flag from the cmd
	title := cmd.Flag("title").Value.String()
	message := cmd.Flag("message").Value.String()

	var answers struct {
		Type        int
		Titre       string
		Description string
	}

	var qs = []*survey.Question{
		{
			Name:     "Type",
			Prompt:   &survey.Select{Message: "Select a category", Options: emojiList(), PageSize: 12, Help: "Gut uses emojis to categorize your commits. Select an emoji that best describes your commit"},
			Validate: survey.Required,
		},
	}
	if title == "" || len(title) > 50 {
		qs = append(qs, &survey.Question{
			Name:     "Titre",
			Prompt:   &survey.Input{Message: "Title of your commit (max 50 chars)", Help: "Ask yourself what you did in this commit | Use active voice | Avoid using 'and' or 'or'"},
			Validate: titleValidation,
		})
	} else {
		answers.Titre = title
	}
	if message == "" {
		qs = append(qs, &survey.Question{
			Name:   "Description",
			Prompt: &survey.Multiline{Message: "Describe your commit (optional)", Help: "Write a description of your commit. Explain why you did this commit and assume that you are explaining to a colleague who knows nothing about the codebase"},
		},
		)
	} else {
		answers.Description = message
	}

	err = survey.Ask(qs, &answers)
	if err != nil {
		exitOnError("Sorry, I can't get your answers", err)
	}

	// Append the title to the body because git only accept a message.
	// However, it's common that the first line is the title and the rest the body
	commitMessage := computeCommitMessage(answers)

	// Commit the changes
	Result, err := executor.Commit(wd, commitMessage)
	if err != nil {
		exitOnError("Error while committing", err)
	}
	print.Message("\n\nChanges updated successfully with commit hash: "+Result.Hash, print.Success)
	fmt.Printf("%d files changed, %d insertions(+), %d deletions(-)\n", Result.FilesUpdated, Result.FilesAdded, Result.FilesDeleted)

}

func emojiList() []string {
	var emojis []string
	for _, e := range gitEmoji {
		emojis = append(emojis, e.Emoji+" "+e.Description)
	}
	return emojis
}

func titleValidation(s interface{}) error {
	val, ok := s.(string)
	if !ok || len(val) == 0 {
		return errors.New("for easy retrieval, don't forget to add a title")
	} else if len(val) > 50 {
		return errors.New("it's recommended to keep the title under 50 characters")
	}
	return nil
}

func computeCommitMessage(answers struct {
	Type        int
	Titre       string
	Description string
}) string {
	var message string
	message += gitEmoji[answers.Type].Emoji + " " + answers.Titre + "\n" + answers.Description
	return message
}

// Check if the username and email for commits are set.
//
// If not, prompt the user to set it
func verifUserConfig(path string) {
	// Get User config for the path
	username, email, err := executor.GetUserConfig(path)
	if err != nil {
		exitOnError("Sorry, I can't get your user config", err)
	}
	if username != "" && email != "" {
		return
	}
	print.Message("Hi there, I'm missing some information about you. Let's fix that!", print.Info)
	var answers struct {
		Username string
		Email    string
	}
	var qs []*survey.Question
	if username == "" {
		qs = append(qs, &survey.Question{
			Name:   "Username",
			Prompt: &survey.Input{Message: "Username", Help: "How do you want to be called in commits messages?"},

			Validate: survey.Required,
		})
	}
	if email == "" {
		qs = append(qs, &survey.Question{
			Name:     "Email",
			Prompt:   &survey.Input{Message: "Email", Help: "What is your email?(Be careful, this is public on GitHub)"},
			Validate: survey.Required,
		})
	}
	answers.Username = username
	answers.Email = email
	err = survey.Ask(qs, &answers)
	if err != nil {
		exitOnError("Sorry, I can't get your answers", err)
	}
	err = executor.SetUserConfig(path, answers.Username, answers.Email)
	if err != nil {
		exitOnError("Sorry, I can't set your user config", err)
	}

}
