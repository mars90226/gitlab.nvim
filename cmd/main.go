package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/harrisoncramer/gitlab.nvim/cmd/app"
	"github.com/harrisoncramer/gitlab.nvim/cmd/app/git"
)

var pluginOptions app.PluginOptions
var Version = "unknown" // Set via ldflags

// The plugin passes its settings (including the auth token) through this environment
// variable, because command-line arguments are visible to every user on the machine.
const settingsEnvVar = "GITLAB_NVIM_SERVER_SETTINGS"

// readServerSettings reads the server configuration from the environment, falling back
// to the first command-line argument for older versions of the Lua plugin. The variable
// is unset so that child processes (e.g. git) do not inherit the auth token.
func readServerSettings(args []string) (string, error) {
	if settings, ok := os.LookupEnv(settingsEnvVar); ok {
		if err := os.Unsetenv(settingsEnvVar); err != nil {
			return "", err
		}
		return settings, nil
	}
	if len(args) >= 2 {
		return args[1], nil
	}
	return "", errors.New("must provide server configuration")
}

func main() {
	log.SetFlags(0)

	/* Set up before anything else, so that a server that is still initializing is
	also shut down when Neovim dies. */
	app.WatchForParentExit()

	settings, err := readServerSettings(os.Args)
	if err != nil {
		log.Fatalf("Failure reading plugin settings: %v", err)
	}

	err = json.Unmarshal([]byte(settings), &pluginOptions)
	app.SetPluginOptions(pluginOptions)
	app.SetVersion(Version)

	if err != nil {
		log.Fatalf("Failure parsing plugin settings: %v", err)
	}

	gitManager := git.Git{}
	gitData, err := git.NewGitData(pluginOptions.ConnectionSettings.Remote, pluginOptions.GitlabUrl, gitManager)

	if err != nil {
		log.Fatalf("Failure initializing plugin: %v", err)
	}

	client, err := app.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize Gitlab client: %v", err)
	}

	projectInfo, err := app.InitProjectSettings(client, gitData)
	if err != nil {
		log.Fatalf("Failed to initialize project settings: %v", err)
	}

	app.StartServer(client, projectInfo, gitData)
}
