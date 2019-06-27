package served

import (
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "served",
	Short: "A useful service installer",
	Long:  "Install programs as a deamon service on major platforms.",
}

var installCmd *cobra.Command
var uninstallCmd *cobra.Command
var configCmd *cobra.Command

func init() {
	initCmd()
	rootCmd.AddCommand(installCmd, uninstallCmd, configCmd)
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

	initConfigCmd()
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

func initConfigCmd() {
	var username string
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	cfgInit := &Config{}
	var out string
	configInitCmd := &cobra.Command{
		Use:   "init",
		Short: "Init a new config file",
		Long:  "Init a new config file with special arguments",
		Run: func(cmd *cobra.Command, args []string) {
			err := startInitConfig(cfgInit, out)
			if err != nil {
				exitWithErr(err)
			}
			fmt.Printf("init config successfully out: %v \n", out)
		},
	}
	configInitCmd.Flags().StringVarP(&cfgInit.Name, "name", "n", "", "program name")
	configInitCmd.MarkFlagRequired("name")
	configInitCmd.Flags().StringVarP(&cfgInit.UserName, "user", "u", username, "run program as username")
	configInitCmd.Flags().StringVarP(&cfgInit.WorkingDirectory, "wd", "", ".", "program work directory")
	configInitCmd.Flags().StringSliceVarP(&cfgInit.Arguments, "args", "", nil, "program arguments")
	configInitCmd.Flags().StringVarP(&out, "out", "o", "config.yaml", "new config file path")

	var editFile string
	defaultValue := "!!!@@@"
	var osxOptSlice, posixOptSlice []string

	var cfgEdit = &Config{}

	configEditCmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a special config file",
		Long:  "Edit a special config file with special arguments",
		Run: func(cmd *cobra.Command, args []string) {
			err := startEditConfig(cfgEdit, osxOptSlice, posixOptSlice, editFile)
			if err != nil {
				exitWithErr(err)
			}
			fmt.Printf("edit config successfully: %v \n", editFile)
		},
	}
	configEditCmd.Flags().StringVarP(&cfgEdit.Name, "name", "n", defaultValue, "program name")
	configEditCmd.Flags().StringVarP(&cfgEdit.DisplayName, "display", "", defaultValue, "program display name")
	configEditCmd.Flags().StringVarP(&cfgEdit.Description, "desc", "", defaultValue, "program description")
	configEditCmd.Flags().StringVarP(&cfgEdit.UserName, "user", "u", defaultValue, "run program as username")
	configEditCmd.Flags().StringVarP(&cfgEdit.WorkingDirectory, "wd", "", defaultValue, "program work directory")
	configEditCmd.Flags().StringVarP(&cfgEdit.ChRoot, "chroot", "", defaultValue, "program change root")
	configEditCmd.Flags().StringSliceVarP(&cfgEdit.Arguments, "args", "", []string{}, "program arguments")
	configEditCmd.Flags().StringSliceVarP(&cfgEdit.Dependencies, "dep", "", []string{}, "program dependencies")
	configEditCmd.Flags().StringSliceVarP(&osxOptSlice, "osx", "", []string{}, "osx option")
	configEditCmd.Flags().StringSliceVarP(&posixOptSlice, "posix", "", []string{}, "posix option")
	configEditCmd.Flags().StringVarP(&editFile, "file", "f", "config.yaml", "program config file to edit")

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Operate config file",
		Long:  "Init or edit special config file",
	}
	configCmd.AddCommand(configInitCmd, configEditCmd)
}

func startInitConfig(c *Config, out string) error {
	if c == nil {
		return fmt.Errorf("nil config")
	}
	c.DisplayName = c.Name
	c.Description = c.Name

	if fileExists(out) {
		return fmt.Errorf("config file exitsed: %v", out)
	}

	return writeConfig2File(c, out)
}

func startEditConfig(c *Config, osxOpt, posixOpt []string, filename string) error {
	if c == nil {
		return fmt.Errorf("nil config")
	}
	originCfg, err := ParseConfig(filename)
	if err != nil {
		return err
	}
	defaultValue := "!!!@@@"
	if c.Name != defaultValue {
		if c.Name == "" {
			return fmt.Errorf("name can not be empty")
		}
		originCfg.Name = c.Name
	}
	updateValue(&originCfg.DisplayName, &c.DisplayName, defaultValue)
	updateValue(&originCfg.Description, &c.Description, defaultValue)
	updateValue(&originCfg.WorkingDirectory, &c.WorkingDirectory, defaultValue)
	updateValue(&originCfg.UserName, &c.UserName, defaultValue)
	updateValue(&originCfg.ChRoot, &c.ChRoot, defaultValue)

	if !isDefaultSlice(c.Arguments) {
		originCfg.Arguments = c.Arguments
	}
	if !isDefaultSlice(c.Dependencies) {
		originCfg.Dependencies = c.Dependencies
	}

	if !isDefaultSlice(osxOpt) {
		originCfg.OSXOpt = goutil.Map{}
		for _, item := range osxOpt {
			kv := strings.SplitN(item, ":", 2)
			if len(kv) != 2 {
				return fmt.Errorf("invalid oxs option: %v", item)
			}
			originCfg.OSXOpt.Set(kv[0], kv[1])
		}
	}

	if !isDefaultSlice(posixOpt) {
		originCfg.POSIXOpt = goutil.Map{}
		for _, item := range posixOpt {
			kv := strings.SplitN(item, ":", 2)
			if len(kv) != 2 {
				return fmt.Errorf("invalid posix option: %v", item)
			}
			originCfg.POSIXOpt.Set(kv[0], kv[1])
		}
	}

	return writeConfig2File(originCfg, filename)
}

func updateValue(origin, newer *string, defaultValue string)  {
	if *newer != defaultValue {
		*origin = *newer
	}
}

func isDefaultSlice(s []string) bool {
	if s != nil && len(s) == 0 {
		return true
	}
	return false
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

func fileExists(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}
	return true
}
