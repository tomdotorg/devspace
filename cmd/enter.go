package cmd

import (
	"github.com/devspace-cloud/devspace/pkg/devspace/cloud"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/configutil"
	"github.com/devspace-cloud/devspace/pkg/devspace/kubectl"
	"github.com/devspace-cloud/devspace/pkg/devspace/services"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/spf13/cobra"
)

// EnterCmd is a struct that defines a command call for "enter"
type EnterCmd struct {
	flags *EnterCmdFlags
}

// EnterCmdFlags are the flags available for the enter-command
type EnterCmdFlags struct {
	selector      string
	namespace     string
	labelSelector string
	container     string
	switchContext bool
	pick          bool
	config        string
}

func init() {
	cmd := &EnterCmd{
		flags: &EnterCmdFlags{},
	}

	cobraCmd := &cobra.Command{
		Use:   "enter",
		Short: "Open a shell to a container",
		Long: `
#######################################################
################## devspace enter #####################
#######################################################
Execute a command or start a new terminal in your 
devspace:

devspace enter
devspace enter -p # Select pod to enter
devspace enter bash
devspace enter -s my-selector
devspace enter -c my-container
devspace enter bash -n my-namespace
devspace enter bash -l release=test
#######################################################`,
		Run: cmd.Run,
	}
	rootCmd.AddCommand(cobraCmd)

	cobraCmd.Flags().StringVarP(&cmd.flags.selector, "selector", "s", "", "Selector name (in config) to select pod/container for terminal")
	cobraCmd.Flags().StringVarP(&cmd.flags.container, "container", "c", "", "Container name within pod where to execute command")
	cobraCmd.Flags().StringVarP(&cmd.flags.labelSelector, "label-selector", "l", "", "Comma separated key=value selector list (e.g. release=test)")
	cobraCmd.Flags().StringVarP(&cmd.flags.namespace, "namespace", "n", "", "Namespace where to select pods")
	cobraCmd.Flags().BoolVar(&cmd.flags.switchContext, "switch-context", false, "Switch kubectl context to the DevSpace context")
	cobraCmd.Flags().BoolVarP(&cmd.flags.pick, "pick", "p", false, "Select a pod to stream logs from")
	cobraCmd.Flags().StringVar(&cmd.flags.config, "config", configutil.ConfigPath, "The DevSpace config file to load (default: '.devspace/config.yaml'")
}

// Run executes the command logic
func (cmd *EnterCmd) Run(cobraCmd *cobra.Command, args []string) {
	// Set config root
	if configutil.ConfigPath != cmd.flags.config {
		configutil.ConfigPath = cmd.flags.config
	}

	// Set config root
	configExists, err := configutil.SetDevSpaceRoot()
	if err != nil {
		log.Fatal(err)
	}
	if !configExists {
		log.Fatal("Couldn't find a DevSpace configuration. Please run `devspace init`")
	}

	log.StartFileLogging()

	// Configure cloud provider
	err = cloud.Configure(log.GetInstance())
	if err != nil {
		log.Fatalf("Unable to configure cloud provider: %v", err)
	}

	// Get kubectl client
	kubectl, err := kubectl.NewClientWithContextSwitch(cmd.flags.switchContext)
	if err != nil {
		log.Fatalf("Unable to create new kubectl client: %v", err)
	}

	// Start terminal
	err = services.StartTerminal(kubectl, cmd.flags.selector, cmd.flags.container, cmd.flags.labelSelector, cmd.flags.namespace, cmd.flags.pick, args, make(chan error), log.GetInstance())
	if err != nil {
		log.Fatal(err)
	}
}
