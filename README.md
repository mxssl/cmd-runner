# cmd-runner

Execute commands on remote hosts asynchronously over ssh

## Installation

Download compiled binary for your system

[Linux](https://github.com/mxssl/cmd-runner/releases/download/0.0.4/cmd-runner-linux-amd64)

[MacOS](https://github.com/mxssl/cmd-runner/releases/download/0.0.4/cmd-runner-darwin-amd64)

Example

```bash
wget https://github.com/mxssl/cmd-runner/releases/download/0.0.4/cmd-runner-linux-amd64 -O cmd-runner
mv cmd-runner /usr/local/bin/cmd-runner
chmod +x /usr/local/bin/cmd-runner/cmd-runner
```

## How to run

### Usecase 1

#### Run commands from local file

1. Create config file - `config.toml`

Example:

```toml
# Credentials
username = "root"
password = "password"

# SSH private and public keys
ssh_private_key = "/home/user/.ssh/id_rsa"
ssh_public_key = "/home/user/.ssh/id_rsa.pub"

# SSH port
ssh_port = "22"

# Connection method: "key" or "password"
connection_method = "key"

# Remote hosts
hosts = [
	"1.1.1.1",
	"2.2.2.2",
	"3.3.3.3"
]

# File with commands for "cmd-runner start" command
commands_file = "commands.txt"

# > Full < source and destination path to file for "cmd-runner copy" command 
source_path = "/opt/scripts/script.sh"
destination_path = "/tmp/script.sh"
```

2. Create file with commands that you want to run on remote hosts - `commands.txt`

Example:

```bash
apt-get update
apt-get upgrade -y
```

3. Start program `cmd-runner start`

If you want files with stdout per host then start `cmd-runner` with key:
```
cmd-runner start --file
```

`cmd-runner` will create files `hostname-output.txt`

### Usecase 2

#### Copy local file to remote hosts

1. The same as `Usecase 1`

2. Start program `cmd-runner copy`

### Usecase 3

#### Combine Usecase 1 and Usecase 2

You can copy bash script to remote hosts with `cmd-runner copy` and run this script with `cmd-runner start`
