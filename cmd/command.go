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
	conf := config.New()
	return &Command{
		Config: conf,
	}
}

func (cmd *Command) GetApplications() []Application {
	files, err := ioutil.ReadDir(cmd.Config.ApplicationPath)

	if err != nil {
		panic(err)
	}

	apps := make([]Application, 0)

	for _, file := range files {
		// TODO support linux
		if file.IsDir() && !cmd.Config.IsDarwin() {
			continue
		}

		name := file.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		apps = append(apps, Application{
			Name: name,
			Path: filepath.Join(cmd.Config.ApplicationPath, name),
		})
	}

	return apps
}

func (cmd *Command) RunApp(path string) error {
	var command *exec.Cmd
	if cmd.Config.IsDarwin() {
		command = exec.Command("open", path)
	}

	// TODO support linux
	if cmd.Config.IsLinux() {

	}

	return command.Run()
}

func (cmd *Command) Run() {
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
