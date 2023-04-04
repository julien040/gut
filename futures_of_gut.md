# The future of gut

## Introduction
First, I want to thank everyone for the enthusiasm gut received lately.  
On Thursday 30 March, gut hit the #2 spot on Hacker News and has since gained over 250 stars. I’m so grateful, thank you!

As a reminder, gut is an alternative CLI for Git. It provides a consistent naming of commands and a useful set of features to make git easier to use.
In this blog post, I’m going to discuss my plans for this project.

## Status of the launch 🎉
For rational people who love statistics, here is a sum up of the launch
- The GitHub repo gained 276 stars ⭐️
- Gut was called 259 times on people who activated telemetry.
- 62% of users run on Linux
- The USA is the country with the most gut users.
- The GItHub repo has been visited over 3000 times
I have received a lot of feedback and I look forward to improving gut.

## Main pain points

### Rewriting history is a feature

I knew rewriting history and running `git push -f` could lead to conflict. So, I made the assumption: avoid rewriting the history at all cost. But I was wrong.  
As multiple [people](https://news.ycombinator.com/item?id=35372990) stated, it’s a feature rather than a possible foot gun.
In the next releases, `gut` will allow change the history safely and with warning so that beginners don’t nuke their repositories.

### Git Submodules

Gut has no support at all for submodules. Even if submodules are an advanced aspect of git, gut must still provide some level of support.
Or alternatively, it could simply display a message saying it’s unsupported instead of a [baffling message](https://github.com/julien040/gut/issues/37).

### Secret management

To reduce the complexity, I thought removing `git add` was a great idea because any beginner would just do `git add -A`.  
But at no time, I thought a credential could be committed because the user forgot to add it to the `.gitignore`
I need to find a solution to this issue.

## Idea

*The following are only suggestion and may not be implemented in gut*
Reading through all the comments gave me many ideas for the future of gut

### `gut switch` behaviour

Gut's syntax is based on subcommands, e.g. all profile-related functions are located under `gut profile ...`. However, having two different commands to manage branches—`gut switch` and `gut branch ...`—doesn't make sense. As a result, `gut switch` will now have a new behavior; it will be capable of switching to a branch or a commit.

Additionally, if your working tree is dirty, you won't be able to switch to an existing branch. To address this, `gut` could, in the future, stash, switch, and stash pop to commit your changes to another branch.

### Rename `gut undo`

In version 0.2.5, `gut undo` hard resets your working tree to the HEAD. But is this the most accurate name for this command? In a graphical user interface, 'undo' generally implies undoing the last action rather than reverting to the latest save. Perhaps a more suitable name exists for this command - one that more accurately describes its function.

### Conventional commit

Not everyone wants emojis in their commit message. Or some people use [Conventional commits.](https://www.conventionalcommits.org/en/v1.0.0/)  
so `gut save` should also support Conventional commit.

## Feedback

Despite being in alpha, `gut` can be downloaded. To do so, check the tutorial [here](https://github.com/julien040/gut#installation).
If you have any feedback, suggestions, or bug reports, feel free to open an issue on GitHub. I'm eager to hear your thoughts!

## Conclusion

I'm still full of gratitude for the enthusiasm `gut` has been met with. 
Despite the amount of work that remains before its v1.0.0 release, I'm certain that it can be done. 
To stay abreast of new developments, you can watch the repo on GitHub and receive notifications for new releases.
Furthermore, if you are confident enough to work on this project, all contributors are welcome.

Julien



