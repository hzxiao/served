package served

import (
	"github.com/hzxiao/goutil"
	"github.com/kardianos/service"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"text/template"
)

// Config provides the setup for a Service. The Name field is required.
type Config struct {
	Name        string   `yaml:"name"`         // Required name of the service. No spaces suggested.
	DisplayName string   `yaml:"display_name"` // Display name, spaces allowed.
	Description string   `yaml:"description"`  // Long description of service.
	UserName    string   `yaml:"username"`     // Run as username.
	Arguments   []string `yaml:"arguments"`    // Run with arguments.

	// Optional field to specify the executable for service.
	// If empty the current executable is used.
	Executable string `yaml:"executable"`

	// Array of service dependencies.
	// Not yet fully implemented on Linux or OS X:
	//  1. Support linux-systemd dependencies, just put each full line as the
	//     element of the string array, such as
	//     "After=network.target syslog.target"
	//     "Requires=syslog.target"
	//     Note, such lines will be directly appended into the [Unit] of
	//     the generated service config file, will not check their correctness.
	Dependencies []string `yaml:"dependencies"`

	// The following fields are not supported on Windows.
	WorkingDirectory string `yaml:"working_directory"` // Initial working directory.
	ChRoot           string `yaml:"chroot"`

	// System specific options.
	//  * OS X
	//    - LaunchdConfig string ()      - Use custom launchd config
	//    - KeepAlive     bool   (true)
	//    - RunAtLoad     bool   (false)
	//    - UserService   bool   (false) - Install as a current user service.
	//    - SessionCreate bool   (false) - Create a full user session.
	OSXOpt goutil.Map `yaml:"osx_opt"`
	//  * POSIX
	//    - SystemdScript string ()                 - Use custom systemd script
	//    - UpstartScript string ()                 - Use custom upstart script
	//    - SysvScript    string ()                 - Use custom sysv script
	//    - RunWait       func() (wait for SIGNAL)  - Do not install signal but wait for this function to return.
	//    - ReloadSignal  string () [USR1, ...]     - Signal to send on reaload.
	//    - PIDFile       string () [/run/prog.pid] - Location of the PID file.
	//    - LogOutput     bool   (false)            - Redirect StdErr & StdOut to files.
	//    - Restart       string (always)           - How shall service be restarted.
	//    - SuccessExitStatus string ()             - The list of exit status that shall be considered as successful,
	//                                                in addition to the default ones.
	POSIXOpt goutil.Map `yaml:"posix_opt"`
}

func ParseConfig(filename string) (*Config, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

type program struct {
	*Config
}

func (prg *program) Start(s service.Service) error {
	return nil
}

func (prg *program) Stop(s service.Service) error {
	return nil
}

func writeConfig2File(c *Config, filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.New("").Parse(configTmpl)
	if err != nil {
		return err
	}

	return t.Execute(f, c)
}

var configTmpl = `# Required name of the service. No spaces suggested.
name: {{.Name}}

# Display name, spaces allowed.
display_name: {{.DisplayName}}

# Long description of service.
description: {{.Description}}

# Run as username.
username: {{.UserName}}

arguments:{{range $i, $arg := .Arguments}} 
  - {{$arg}} {{end}}

# Array of service dependencies.
# Not yet fully implemented on Linux or OS X:
# 1. Support linux-systemd dependencies, just put each full line as the
#     element of the string array, such as
#     "After=network.target syslog.target"
#     "Requires=syslog.target"
#     Note, such lines will be directly appended into the [Unit] of
#     the generated service config file, will not check their correctness.
dependencies:{{range $i, $dep := .Dependencies}} 
  - {{$dep}} {{end}}


# The following fields are not supported on Windows.
# Initial working directory.
working_directory: {{.WorkingDirectory}}
chroot: {{.ChRoot}}


# System specific options and default value.
#* OS X
#   LaunchdConfig: ""    - Use custom launchd config
#   KeepAlive: true
#   RunAtLoad: false
#   UserService: false   - Install as a current user service.
#   SessionCreate: false - Create a full user session.
osx_opt:{{range $k, $v := .OSXOpt}} 
  {{$k}}: {{$v}} {{end}}

# * POSIX
#   SystemdScript: ""                 - Use custom systemd script
#   UpstartScript: ""                 - Use custom upstart script
#   SysvScript: ""                    - Use custom sysv script
#   ReloadSignal: "USR1, ..."         - Signal to send on reaload.
#   PIDFile: ""                       - Location of the PID file.
#   LogOutput: false                  - Redirect StdErr & StdOut to files.
#   Restart: always                   - How shall service be restarted.
#   SuccessExitStatus: ""             - The list of exit status that shall be considered as successful,
#                                     in addition to the default ones.
posix_opt:{{range $k, $v := .POSIXOpt}} 
  {{$k}}: {{$v}} {{end}}`
