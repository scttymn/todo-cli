# Todo CLI

A Git-integrated CLI tool for managing branch-specific todo lists. Keep track of tasks for each feature branch you're working on.

## Overview

Todo CLI creates and manages markdown-based todo lists that are tied to your Git branches. Each feature branch gets its own todo list, helping you stay organized across different features and projects.

## Installation

```bash
go build -o todo
# Move to your PATH if desired
sudo mv todo /usr/local/bin/
```

## Quick Start

```bash
# Create or switch to a todo list (creates feature/my-feature branch)
todo list my-feature

# Add todos to the current list
todo add "Implement user authentication"
todo add "Write unit tests"
todo add "Update documentation"

# View current list progress  
todo progress

# Mark items as completed
todo check 1
todo check 2

# View all lists and their progress
todo list

# Switch to a different list
todo list main
```

## Commands

### `todo list [list-name]`
Create, switch to, or view todo lists.

- `todo list` - Show all available lists with progress
- `todo list <name>` - Switch to or create a list (creates `feature/<name>` branch)
- `todo list --delete <name>` - Delete a list and its branch
- `todo list -d <name>` - Short form of delete

### `todo add <item>`
Add a new todo item to the current list.

```bash
todo add "Fix the login bug"
todo add "Refactor authentication module"
```

### `todo check <number>`
Mark a todo item as completed.

```bash
todo check 1
todo check 3
```

### `todo uncheck <number>`
Mark a todo item as incomplete.

```bash
todo uncheck 2
```

### `todo progress [list-name]`
Show progress for todo lists.

- `todo progress` - Show current list progress
- `todo progress <name>` - Show progress for specific list  
- `todo progress --all` - Show progress for all lists
- `todo progress -a` - Short form of --all

### `todo version`
Display the CLI version.

## How It Works

Todo CLI integrates with Git to provide branch-specific todo lists:

1. **Branch Integration**: Each todo list corresponds to a Git branch following the pattern `feature/<list-name>`
2. **File Storage**: Todo items are stored in markdown files at `.todo/<list-name>.md`
3. **Git Safety**: Warns about uncommitted changes before switching branches
4. **Automatic Creation**: Lists and branches are created automatically when needed

## File Structure

```
project/
├── .todo/
│   ├── my-feature.md
│   ├── bug-fixes.md
│   └── main.md
└── [your project files]
```

Each `.md` file contains a standard markdown todo list:

```markdown
# Todo List for my-feature

- [ ] Implement user authentication
- [x] Write unit tests  
- [ ] Update documentation
```

## Examples

### Working on a New Feature

```bash
# Start working on authentication feature
todo list authentication

# Add some todos
todo add "Set up OAuth provider"
todo add "Create user model"  
todo add "Add login/logout routes"

# Work on tasks, mark them complete
todo check 1
todo check 2

# Check progress
todo progress
# Output: Todo list for branch 'authentication':
#         1. [x] Set up OAuth provider
#         2. [x] Create user model
#         3. [ ] Add login/logout routes
#         Progress: 2/3 completed
```

### Managing Multiple Lists

```bash
# See all your lists
todo list
# Output: Lists:
#         authentication - 2/3 completed (67%)
#         bug-fixes - 1/5 completed (20%)
#         main - 0/2 completed (0%)

# Switch between lists
todo list bug-fixes
todo list main
```

### Cleaning Up Completed Work

```bash
# Delete a completed feature list
todo list --delete authentication
# Prompts: Are you sure you want to delete list 'authentication'? 
#          This will remove both the branch and todo file. (y/N): y
# Output: Successfully deleted list 'authentication'
```

## Safety Features

- **Uncommitted Changes Warning**: Warns before switching branches if you have uncommitted changes
- **Delete Confirmation**: Requires confirmation before deleting lists
- **Branch Protection**: Cannot delete the list you're currently working on
- **Git Repository Check**: Provides helpful messages when not in a git repository

## Requirements

- Go 1.19+
- Git (must be in a git repository to use)

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b my-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Submit a pull request

## License

MIT License - see LICENSE file for details.