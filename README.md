# cmd-runner

Run commands on remote hosts over ssh and write stdout to files

# install

Download compiled binary for your system

[Linux](https://github.com/mxssl/cmd-runner/releases/download/0.0.2/cmd-runner-linux-amd64)

[MacOS](https://github.com/mxssl/cmd-runner/releases/download/0.0.2/cmd-runner-darwin-amd64)

[Windows](https://github.com/mxssl/cmd-runner/releases/download/0.0.2/cmd-runner-windows-amd64)

Example
```
wget https://github.com/mxssl/cmd-runner/releases/download/0.0.2/cmd-runner-linux-amd64 -O cmd-runner
chmod +x cmd-runner
```

# run

1. Create config file - `config.toml`

Example:
```
username = "user"
password = "pass"

hosts = [
	"1.1.1.1",
	"2.2.2.2"
]

commands_file = "commands.txt"
```

2. Create file with commands that you want to run on remote hosts - `commands.txt`

Example:
```
apt-get update
apt-get upgrade -y
```

3. Start program `./cmd-runner start`

4. Program creates files with stdout - `hostname-output.txt`
