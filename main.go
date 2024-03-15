package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
)

type Item struct {
	Code string `json:"code" yaml:"code"`
	Desc string `json:"desc" yaml:"desc"`
}

type ListOfItem []Item

func (l *ListOfItem) ToSliceString() []string {
	result := make([]string, len(*l))
	for i := range *l {
		code := (*l)[i].Code
		desc := l.cutStr((*l)[i].Desc)

		if code != "" && desc != "" {
			result[i] = code + ": " + desc
		} else if code != "" {
			result[i] = code
		}
	}
	return result
}

func (l *ListOfItem) cutStr(str string) string {
	maxLength := 64
	if len(str) <= maxLength {
		return str
	}
	return str[0:maxLength-3] + "..."
}

type Config struct {
	Types  []Item `json:"types" yaml:"types"`
	Scopes []Item `json:"scopes" yaml:"scopes"`
}

var configPath string

func init() {
	u, _ := user.Current()
	flag.StringVar(&configPath, "c", filepath.Join(u.HomeDir, ".gitcz/config.json"), "configuration file path")
	flag.Parse()
}

func main() {
	if err := checkStageFile(); err != nil {
		pterm.Info.Println("No staged file: " + err.Error())
		os.Exit(1)
	}

	conf, err := loadConfigFile()
	if err != nil {
		pterm.Info.Println("Fail to load config file: " + err.Error())
		os.Exit(1)
	}
	defer saveConfig(conf)

	chosenType, err := chooseType(conf.Types)
	if err != nil {
		pterm.Info.Println("Fail to choose type: " + err.Error())
		os.Exit(1)
	}

	chosenScope, err := chooseScope(conf.Scopes)
	if err != nil {
		pterm.Info.Println("Fail to choose scope: " + err.Error())
		os.Exit(1)
	}

	title, err := writeTitle()
	if err != nil {
		pterm.Info.Println("Fail to write title: " + err.Error())
		os.Exit(1)
	}

	msg := buildCommitMessage(chosenType, chosenScope, title)
	if err := executeCommit(msg); err != nil {
		pterm.Info.Println("Fail to commit: " + err.Error())
		os.Exit(1)
	}
}

func loadConfigFile() (*Config, error) {
	c := new(Config)

	err := createConfigFile(configPath)
	if err != nil {
		return c, err
	}

	str, err := openFile(configPath)
	if err != nil {
		return c, err
	}

	json.Unmarshal([]byte(str), &c)

	if len(c.Types) == 0 {
		c.Types = defaultTypes()
	} else {
		types := make([]Item, 0)
		copy(types, c.Types)

		c.Types = nil
		for i := range types {
			if types[i].Code != "" {
				c.Types = append(c.Types, types[i])
			}
		}
	}

	if len(c.Scopes) == 0 {
		c.Scopes = defaultScopes()
	} else {
		scopes := make([]Item, 0)
		copy(scopes, c.Types)

		c.Scopes = nil
		for i := range scopes {
			if scopes[i].Code != "" {
				c.Scopes = append(c.Types, scopes[i])
			}
		}
	}
	return c, nil
}

func saveConfig(c *Config) {
	b, err := json.Marshal(c)
	if err != nil {
		pterm.Error.Println("Cannot convert config to json: " + err.Error())
	}
	err = os.WriteFile(configPath, b, os.ModePerm)
	if err != nil {
		pterm.Error.Println("Cannot save config: " + err.Error())
	}
}

func createConfigFile(path string) error {
	paths := strings.Split(path, string(os.PathSeparator))
	if len(paths) > 1 {
		err := os.MkdirAll(strings.Join(paths[:len(paths)-1], string(os.PathSeparator)), os.ModePerm)
		if err != nil {
			return err
		}
	}
	_, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	return err
}

func openFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func checkStageFile() error {
	cmd := exec.Command("git", "diff", "--cached")
	b, err := cmd.Output()
	if err != nil {
		return err
	}
	if len(string(b)) == 0 {
		return fmt.Errorf("cached diff is empty")
	}
	return nil
}

func chooseType(listOfType ListOfItem) (string, error) {
	selected, err := pterm.DefaultInteractiveSelect.
		WithOptions(listOfType.ToSliceString()).
		Show("commit type")
	if err != nil {
		return "", err
	}
	split := strings.SplitN(selected, ":", 2)
	if len(split) == 0 {
		return "", fmt.Errorf("error selected scope of change is empty")
	}
	return strings.Trim(split[0], " "), nil
}

func chooseScope(listOfScope ListOfItem) (string, error) {
	selected, err := pterm.DefaultInteractiveSelect.
		WithOptions(listOfScope.ToSliceString()).
		Show("scope of change")
	if err != nil {
		return "", err
	}
	split := strings.SplitN(selected, ":", 2)
	if len(split) == 0 {
		return "", fmt.Errorf("error selected scope of change is empty")
	}
	return strings.Trim(split[0], " "), nil
}

func writeTitle() (string, error) {
	return pterm.DefaultInteractiveTextInput.Show("Title of commit")
}

func buildCommitMessage(chosenType, chosenScope, title string) string {
	return fmt.Sprintf("%s(%s): %s", chosenType, chosenScope, title)
}

func executeCommit(msg string) error {
	cmd := exec.Command("git", "commit", "-m", msg)
	if err := cmd.Run(); err != nil {
		return err
	}
	pterm.Info.Println(msg)
	return nil
}

func defaultTypes() []Item {
	return []Item{
		{
			Code: "feat",
			Desc: "introduce new functionality.",
		},
		{
			Code: "fix",
			Desc: "address and resolve issues or bugs.",
		},
		{
			Code: "refac",
			Desc: "improve or organize code structure without changing the behavior.",
		},
		{
			Code: "docs",
			Desc: "update documentation or comments within the code itself.",
		},
		{
			Code: "clean",
			Desc: "remove unused code or redundant code.",
		},
		{
			Code: "deps",
			Desc: "update dependencies ensuring compatibility.",
		},
		{
			Code: "config",
			Desc: "modify config, such as scripts, environment, or CI/CD.",
		},
		{
			Code: "opt",
			Desc: "optimize code or algorithms for better performance or efficiency.",
		},
		{
			Code: "style",
			Desc: "update code style, such as guidelines, indentation, naming conventions, or formatting.",
		},
		{
			Code: "local",
			Desc: "add or update localization files.",
		},
		{
			Code: "test",
			Desc: "add, update, or fix tests to ensure code quality and functionality.",
		},
		{
			Code: "revert",
			Desc: "undo previous commit changes.",
		},
		{
			Code: "merge",
			Desc: "merge changes from one branch into another.",
		},
		{
			Code: "sec",
			Desc: "address security vulnerabilities or weaknesses.",
		},
		{
			Code: "setup",
			Desc: "setup the initial project structure, development tools or environment.",
		},
		{
			Code: "debug",
			Desc: "commits for troubleshooting issues.",
		},
	}
}

func defaultScopes() []Item {
	return []Item{
		{
			Code: "environment",
			Desc: "changes to project settings, config, or dependencies, updates to local, staging, or production, as well as changes to env variables or config files.",
		},
		{
			Code: "file",
			Desc: "involve modifications to individual files within the codebase, such as adding, editing, or deleting files.",
		},
		{
			Code: "directory",
			Desc: "changes to entire directories or folders within the project structure, including additions, modifications, or removals of directories and their contents.",
		},
		{
			Code: "database",
			Desc: "involve changes to the database schema, migrations, queries, or configurations, including additions, modifications, or removals of database tables, columns, indexes, or constraints.",
		},
		{
			Code: "server",
			Desc: "changes to server configurations, settings, or infrastructure, including updates to server configurations, deployments, server-side scripts, or server-related dependencies.",
		},
	}
}
