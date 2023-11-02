# Contributing

Hey there! Thanks for your interest in contributing to this project.

## Table of Contents

- [Contributing](#contributing)
  - [Table of Contents](#table-of-contents)
  - [Getting Started](#getting-started)
    - [Reporting Bugs](#reporting-bugs)
    - [Suggesting Enhancements](#suggesting-enhancements)
    - [Pull Requests](#pull-requests)
  - [License](#license)
  - [Development guide](#development-guide)
    - [Setup](#setup)
    - [Folder structure](#folder-structure)
    - [Build for production](#build-for-production)
    - [Commit message format](#commit-message-format)
    - [Code formatting](#code-formatting)

## Getting Started

### Reporting Bugs

If you find a bug, please create an issue on GitHub. Make sure to include as much information as possible so that anyone can try to reproduce the bug.

### Suggesting Enhancements

If you have an idea for a new feature, please create an issue on GitHub.

### Pull Requests

Before submitting a pull request, let's discuss the changes you'd like to make. You can create an issue on GitHub to start the discussion.
Any changes are welcome as long as they align with the spirit of the project. Please make sure to update any tests and documentation as needed.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
Any contributions to this project are made under the same license.

## Development guide

### Setup

`gut` is a statically compiled binary written in Golang. To develop the project, you will only need Go installed on your machine.

To get started, clone the repository and run `go build -o gut main.go` to build the binary.

### Folder structure

Except for main.go, the source code is split into several folders within the `src` folder:

- `cmd`: contains the main entrypoints of the application. Each `.go` file in this folder is a subcommand of `gut`. No code should be written in this folder. It's only a declaration of the subcommands. Each declaration must then call the associated function in the controller package.
- `controller`: contains the business logic of each command. Each `.go` file in this folder is a subcommand of `gut`.
- `executor`: contains the code to interact with the git repository.
- `print`: helpers to print messages to the console.
- `profile`: contains the code to interact with the authentication profile.
- `prompt`: helpers to prompt the user for input.
- `telemetry`: handles the telemetry of the application.
- `provider`: contains the code to interact with some git hosting providers (GitHub, GitLab, Bitbucket, etc.).

### Build for production

To build the binary for production on all supported platforms, run `./script/release.sh` from the root of the project.

Before running the command, make sure you have tagged the last commit with the version number you want to release.

### Commit message format

This project uses [Gitmoji](https://gitmoji.dev/) to add emojis to commit messages. If you don't know how to use it,
don't hesitate to commit your code with `gut`.

### Code formatting

This project uses [gofmt](https://golang.org/cmd/gofmt/) to format the code.
