# served

Install programs as a deamon service on major platforms.

## Suport platform

* Windows XP
* Linux/(systemd | Upstart | SysV)
* OSX/Launchd

## Installation

### Download Reease

### Build

```shell
go get github.com/hzxiao/served
cd $GOPATH/src/github.com/hzxiao/served/cmd
go install
```



## Usage

1. Init a new program config file

```shell
served config init --name prg --out ./config.yaml
```

There will be a new yaml format file in current dir, with conten in it:

```yaml
# Required name of the service. No spaces suggested.
name: prg

# Display name, spaces allowed.
display_name: prg

# Long description of service.
description: prog's desc

# Run as username.
username: your username

arguments:

# Array of service dependencies.
# Not yet fully implemented on Linux or OS X:
# 1. Support linux-systemd dependencies, just put each full line as the
#     element of the string array, such as
#     "After=network.target syslog.target"
#     "Requires=syslog.target"
#     Note, such lines will be directly appended into the [Unit] of
#     the generated service config file, will not check their correctness.
dependencies:


# The following fields are not supported on Windows.
# Initial working directory.
working_directory: .
chroot: 


# System specific options and default value.
#* OS X
#   LaunchdConfig: ""    - Use custom launchd config
#   KeepAlive: true
#   RunAtLoad: false
#   UserService: false   - Install as a current user service.
#   SessionCreate: false - Create a full user session.
osx_opt:

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
posix_opt:
```

2. Edit config file (option)

```shell
# edit program display name
served config edit -f ./config.yaml --display="my prog"

# edit program arguments
served config edit -f ./config.yaml --args="-c conf.perpoties"

# see more info about edit command
sserved config edit -h
```



3. Install program as deamon service

```shell
served install --config ./config.yaml
```



4. uninstall service

```shell
served uninstall --name prg
```

