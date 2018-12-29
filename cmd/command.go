package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	pipeline "github.com/mattn/go-pipeline"
	"github.com/skanehira/gol/config"
)

var (
	fzfMode  *bool
	listMode *bool
	spec     *string
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

func (cmd *Command) getApplications(dir string) (apps []Application) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		name := file.Name()

		//  exlusion dotfiles
		if strings.HasPrefix(name, ".") {
			continue
		}

		// if dir
		if file.IsDir() && !strings.HasSuffix(name, ".app") {
			apps = append(apps, cmd.getApplications(filepath.Join(dir, name))...)
		} else {
			apps = append(apps, Application{
				Name: name,
				Path: filepath.Join(dir, name),
			})
		}
	}

	return apps
}

func (cmd *Command) runApp(path string) {
	var command *exec.Cmd

	switch cmd.Config.OS {
	case config.MacOS:
		command = exec.Command("open", path)

	// TODO support linux
	case config.Linux:

	// TODO support windows
	case config.Windows:

	}

	if err := command.Run(); err != nil {
		panic(err)
	}
}

func (cmd *Command) parseArgs() {
	fzfMode = flag.Bool("f", false, "fzf mode")
	listMode = flag.Bool("l", false, "output application list")
	spec = flag.String("s", "", "specified application")
	flag.Parse()
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

func getPathFromAppName(apps []Application, name string) string {
	name = strings.ToLower(name)
	for _, app := range apps {
		if strings.Contains(strings.ToLower(app.Name), name) {
			return app.Path
		}
	}
	return ""
}

func (cmd *Command) Run() {
	// parse args
	cmd.parseArgs()

	// get application info
	apps := cmd.getApplications(cmd.Config.ApplicationPath)

	// get specified apps
	if *spec != "" {
		for _, name := range strings.Split(*spec, ",") {
			if name != "" {
				cmd.runApp(getPathFromAppName(apps, name))
			}
		}
		return
	}

	// use fzf
	if *fzfMode {
		paths := make(map[string]string)
		var appNames string

		for _, app := range apps {
			name := app.Name
			paths[name] = app.Path
			appNames += name + "\n"
		}

		selected, err := pipeline.Output(
			[]string{"echo", strings.TrimRight(appNames, "\n")},
			[]string{"fzf"},
		)
		if err != nil {
			if strings.Contains(err.Error(), "exit status") {
				return
			}
			panic(err)
		}

		for _, app := range strings.Split(strings.TrimRight(string(selected), "\n"), "\n") {
			cmd.runApp(paths[string(app)])
		}
		return
	}

	// list applications
	if *listMode {
		for _, app := range apps {
			fmt.Println(app.Name)
		}
		return
	}

	// use default
	prompt := promptui.Select{
		Label: "Applications",
		Templates: &promptui.SelectTemplates{
			Label:    `{{ . | green }}`,
			Active:   `{{ .Name | red }}`,
			Inactive: ` {{ .Name | cyan }}`,
			Selected: `{{ .Name | yellow }}`,
		},
		Items: apps,
		Size:  50,
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
		panic(err)
	}

	cmd.runApp(apps[i].Path)

}
