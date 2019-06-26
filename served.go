package served

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "served",
	Short: "A useful service installer",
	Long:  "Install programs as a deamon service on major platforms.",
}

var installCmd *cobra.Command
var uninstallCmd *cobra.Command

func init() {
	initCmd()
	rootCmd.AddCommand(installCmd, uninstallCmd)
}

var (
	cfgFile string
	name    string
)

func initCmd() {
	installCmd = &cobra.Command{
		Use:   "install",
		Short: "install service",
		Long:  "install service with special config file",
		Run:   install,
	}
	installCmd.Flags().StringVarP(&cfgFile, "config", "", "", "config file to install service")
	installCmd.MarkFlagRequired("config")

	uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "uninstall service",
		Long:  "uninstall service with service name",
		Run:   uninstall,
	}
	uninstallCmd.Flags().StringVarP(&name, "name", "", "", "name of service to uninstall")
	uninstallCmd.MarkFlagRequired("name")
}

func install(cmd *cobra.Command, args []string) {
	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		exitWithErr(err)
	}
	prg := &program{cfg}

	srvCfg := &service.Config{
		Name:             prg.Name,
		DisplayName:      prg.DisplayName,
		Description:      prg.Description,
		UserName:         prg.UserName,
		Arguments:        prg.Arguments,
		Dependencies:     prg.Dependencies,
		WorkingDirectory: prg.WorkingDirectory,
		ChRoot:           prg.ChRoot,
	}
	srvCfg.Option = service.KeyValue{}
	for k, v := range prg.OSXOpt {
		srvCfg.Option[k] = v
	}
	for k, v := range prg.POSIXOpt {
		srvCfg.Option[k] = v
	}

	s, err := service.New(prg, srvCfg)
	if err != nil {
		exitWithErr(err)
	}

	err = s.Install()
	if err != nil {
		fmt.Printf("install %v fail with %v platform\n", prg.Name, s.Platform())
		exitWithErr(err)
	}

	fmt.Printf("install %v successfully with %v platform\n", prg.Name, s.Platform())
}

func uninstall(cmd *cobra.Command, args []string) {
	prg := &program{&Config{Name: name}}
	srvCfg := &service.Config{Name: prg.Name}

	s, err := service.New(prg, srvCfg)
	if err != nil {
		exitWithErr(err)
	}

	err = s.Uninstall()
	if err != nil {
		fmt.Printf("uninstall %v fail with %v platform\n", prg.Name, s.Platform())
		exitWithErr(err)
	}

	fmt.Printf("uninstall %v successfully with %v platform\n", prg.Name, s.Platform())
}

func exitWithErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func Run() error {
	return rootCmd.Execute()
}

