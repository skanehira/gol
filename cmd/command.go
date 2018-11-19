package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/skanehira/gol/config"
)

type Application struct {
	Name string
	Path string
}

type Command struct {
	Config *config.Config
}

func New() *Command {
	return &Command{
		Config: config.New(),
	}
}

func (cmd *Command) GetApplications() []Application {
	return dirWalk(cmd.Config.ApplicationPath)
}

func dirWalk(dir string) []Application {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var apps []Application
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if file.IsDir() && !strings.HasSuffix(name, ".app") {
			apps = append(apps, dirWalk(filepath.Join(dir, name))...)
			continue
		}
		apps = append(apps, Application{
			Name: name,
			Path: filepath.Join(dir, name),
		})
	}

	return apps
}

func (cmd *Command) RunApp(path string) error {
	var command *exec.Cmd
	if cmd.Config.IsDarwin() {
		command = exec.Command("open", path)
	}

	// TODO support linux and
	if cmd.Config.IsLinux() {

	}

	// TODO support windows
	if cmd.Config.IsWindows() {

	}

	return command.Run()
}

func (cmd *Command) Run() {
	for {
		apps := cmd.GetApplications()

		prompt := promptui.Select{
			Label: "Applications",
			Templates: &promptui.SelectTemplates{
				Label:    `{{ . | green }}`,
				Active:   `{{ .Name | red }}`,
				Inactive: ` {{ .Name | cyan }}`,
				Selected: `{{ .Name | yellow }}`,
			},
			Items: apps,
			Size:  20,
			Searcher: func(input string, index int) bool {
				item := apps[index]
				name := strings.Replace(strings.ToLower(item.Name), " ", "", -1)
				input = strings.Replace(strings.ToLower(input), " ", "", -1)

				return strings.Contains(name, input)
			},
			StartInSearchMode: true,
		}

		i, _, err := prompt.Run()

		if err != nil {
			if isEOF(err) || isInterrupt(err) {
				os.Exit(0)
			}
			fmt.Println(err)
			os.Exit(-1)
		}

		cmd.RunApp(apps[i].Path)

	}
}

func isEOF(err error) bool {
	if err == promptui.ErrEOF {
		return true
	}

	return false
}

func isInterrupt(err error) bool {
	if err == promptui.ErrInterrupt {
		return true
	}

	return false
}
