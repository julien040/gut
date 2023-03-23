# Gut

> ⚠️ Gut is still in alpha which means that there may be some features that are missing

Gut is a CLI designed to make Git easier to use.

If you have years of muscle memory, gut isn’t probably for you.

## Table of contents

- [Example](#example)
- [Features](#features)
- [Main Useful Commands](#main-useful-commands)
- [Installation](#installation)
  - [Windows](#windows)
  - [MacOS](#macos)
  - [Apt-get](#apt-get)
  - [Yum](#yum)
  - [Build from source](#build-from-source)
- [Principles](#principles)
    - [Integration with cloud](#integration-with-cloud)
    - [No rewriting of history](#no-rewriting-of-history)
    - [Staging area isn’t a thing](#staging-area-isnt-a-thing)
    - [Detached HEAD isn’t scary](#detached-head-isnt-scary)
    - [Great user experience](#great-user-experience)
    - [Coexist with Git](#coexist-with-git)
- [Documentation](#documentation)
- [FAQ](#faq)
    - [Why was this project built?](#why-was-this-project-built)
    - [How can I contact the developer?](#how-can-i-contact-the-developer)
    - [Can I contribute](#can-i-contribute)
- [License](#license)
- [Authors](#authors)
- [Contributing](#contributing)


## Example

```bash
cd my-awesome-project

# Init a new git repo
gut init

# Do some changes
touch my-billion-dollar-idea.txt

# Commit your new file
gut save # Alias of gut commit

# Sync your changes with the upstream
gut sync
```

## Features

- Built-in credentials manager
- Consistent naming of commands
- Integration with cloud platforms (merge and diff opens in the web UI)
- `gut fix` command helps you fix your mistakes with git
- `.gitignore` template downloader
- Simplified authentication with GitHub

## Main Useful Commands

- `gut save` - Commits changes using gitmoji
- `gut sync` - Syncs changes with your remote repository
- `gut goto` - Lets you rewind the state of your project to a particular commit by temporarily modifying the working tree
- `gut fix` - Helps you fix your mistakes with git
- `gut revert` - Reverts your project to a previous state to fix a bug introduced n commits ago
- `gut undo` - Discards changes made since the last commit
- `gut ignore` - Downloads templates of `.gitignore`
- `gut whereami` - Shows where your HEAD points to (no more `rev-parse`)
- `gut switch` - Creates a new branch or switches to an existing one

## Installation

### Windows

To install gut on Windows, run

```bash
scoop bucket add gut https://github.com/julien040/gut-scoop
scoop install gut/gut
```

### MacOS

To install gut on macOS, open the Terminal and run

```bash
brew tap julien040/gut && brew install gut
```

### Apt-get

```bash
echo "deb [trusted=yes] https://apt.gut-cli.dev /" | sudo tee /etc/apt/sources.list.d/gut.list
sudo apt-get update
sudo apt-get install gut
```

### Yum

```bash
sudo tee /etc/yum.repos.d/gut.repo <<EOF
[gut]
name=Gut Repository
baseurl=https://yum.gut-cli.dev/
enabled=1
gpgcheck=0
EOF
sudo yum update
sudo yum install gut
```

### Build from source

You need to have go installed on your machine

```bash
go install github.com/julien040/gut@latest
```

## Principles

### Integration with cloud

We have several tools like GitHub, GitLab, or BitBucket. So, why not use them to their fullest? When you attempt to merge a branch, gut will open a page to create a pull request. If you want to compare two commits, gut will open the compare view in your favorite repository hosting.

### No rewriting of history

Gut will never allow you to modify the history pushed to a remote repository. 

If you want to cancel your changes, run `gut revert`. It will create a new commit containing the state of the commit selected. 

If you made a typo in your commit message, gut will only allow you to change it if it hasn’t been pushed yet. The same thing applies if you attempt to amend files of a commit.

### Staging area isn’t a thing

Everything git tracks will be saved in your next commit. 

Gut tracks all files unless they are listed on your `.gitignore`. After all, it exists for a reason.

### Detached HEAD isn’t scary

With gut, getting into a detached HEAD is pretty trivial. Just run `gut goto <commit id>` to change your working tree according to that commit. 
But gut won’t leave you there. 

In detached HEAD, every operation is blocked until you do something with that commit. You can come back to a branch or create a new one from that commit. 
By blocking operations, you won’t create several commits before realizing they aren’t linked to anything.

### Great user experience

If you make a mistake, gut will try to figure it out and prompt you again. 
And if gut can’t help you, it will do its best to guide you to solve the issue.

I believe that when learning new technology, it's best to start with a high-level understanding and then gradually delve deeper over time. Specifically, when it comes to Git, I find that it can seem complex right from the beginning (such as setting `user.email` and `user.name`)

### Coexist with Git

While **`gut`** is a useful CLI tool, it's not intended to replace the **`git`** CLI. In fact, if Git were to disappear, **`gut`** would no longer function since it heavily relies on Git commands internally.

That being said, I believe that **`gut`** is an excellent choice for simple tasks, while Git can still be utilized for more complex tasks.

## Documentation

[Documentation](https://gut-cli.dev/docs)

## FAQ

### Why was this project built?

In my two years of learning how to code, I found `git` to be extremely frustrating. I was always scared of doing the wrong thing and not being able to revert it. This is why I built `gut` - so that everyone can use `git` without the headaches.

### How can I contact the developer?

To discuss a new feature you would like to see, open a new discussion on GitHub. 

For a bug, open a new issue. 

For anything related to security, commercial or press, send an email to [contact@julienc.me](mailto:contact@julienc.me).

### Can I contribute?

Of course, you can!

## Roadmap

- Create and delete tags.
- Conflict resolution.
- `gut restore` to checkout specific files.
- `gut commit sparsely` to create a commit with specified files rather than all files.
- `gut time-machine` to go back in time (e.g. reverse a pull).
- Open a new discussion if you want your feature to be added here!

## License

[MIT](https://choosealicense.com/licenses/mit/)

## Authors

- [@julien040](https://github.com/julien040)

## Contributing

Contributions are always welcome!

See `contributing.md` for ways to get started.

Please adhere to this project's `code of conduct`.
