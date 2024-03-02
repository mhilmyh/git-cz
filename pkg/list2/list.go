package list2

import "github.com/charmbracelet/bubbles/list"

type Item struct {
	code, title, desc string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title + i.desc }

func NewListItemTypeOfChange() []list.Item {
	return []list.Item{
		Item{
			code:  "feat",
			title: "Feature Addition",
			desc:  "Introduce new functionality or features to the project.",
		},
		Item{
			code:  "fix",
			title: "Bug Fix",
			desc:  "Address and resolve issues or bugs in the codebase.",
		},
		Item{
			code:  "refactor",
			title: "Refactoring",
			desc:  "Improve the code structure, readability, or performance without changing its external behavior.",
		},
		Item{
			code:  "docs",
			title: "Documentation Update",
			desc:  "Update or add documentation, such as README files, comments, or documentation within the code itself.",
		},
		Item{
			code:  "deps",
			title: "Dependency Update",
			desc:  "Update dependencies to newer versions, ensuring compatibility and security.",
		},
		Item{
			code:  "clean",
			title: "Code Cleanup",
			desc:  "Remove unused code, refactor redundant code, or improve code organization without changing its functionality.",
		},
		Item{
			code:  "config",
			title: "Configuration Change",
			desc:  "Modify project configurations, such as build scripts, environment settings, or CI/CD configurations.",
		},
		Item{
			code:  "optimize",
			title: "Optimization",
			desc:  "Optimize code or algorithms for better performance or efficiency.",
		},
		Item{
			code:  "style",
			title: "Code Style Changes",
			desc:  "Enforce or update code style guidelines, such as indentation, naming conventions, or code formatting.",
		},
		Item{
			code:  "local",
			title: "Localization",
			desc:  "Add or update translations or localization files.",
		},
		Item{
			code:  "test",
			title: "Testing",
			desc:  "Add, update, or fix tests to ensure code quality and functionality.",
		},
		Item{
			code:  "revert",
			title: "Revert",
			desc:  "Undo previous changes, reverting the codebase to a previous state.",
		},
		Item{
			code:  "merge",
			title: "Merge",
			desc:  "Merge changes from one branch into another, typically seen in pull request merges.",
		},
		Item{
			code:  "security",
			title: "Security Fix",
			desc:  "Address security vulnerabilities or weaknesses in the codebase.",
		},
		Item{
			code:  "env",
			title: "Environment Setup",
			desc:  "Initialize or set up the development environment, such as configuring development tools, setting up local or remote development servers, or establishing project-specific environment variables.",
		},
		Item{
			code:  "debug",
			title: "Debugging",
			desc:  "Address issues identified during debugging sessions. These commits typically involve fixing errors, resolving unexpected behavior, or adding logging or debugging statements to aid in diagnosing and fixing problems within the codebase.",
		},
	}
}

func NewListItemScopeOfChange() []list.Item {
	return []list.Item{
		Item{
			code:  "domain",
			title: "Domain",
			desc:  "...",
		},
		Item{
			code:  "provider",
			title: "Provider",
			desc:  "...",
		},
		Item{
			code:  "handler",
			title: "Handler",
			desc:  "...",
		},
	}
}
