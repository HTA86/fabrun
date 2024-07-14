# `FabRun`

### **Disclaimer:** FabRun is currently under development and is not ready for production use. Use at your own risk.

FabRun is a simple utility designed to run predefined commands stored in command.md files. Initially created to simplify running long and complex [Fabric](https://github.com/danielmiessler/fabric) commands, FabRun now serves as a general-purpose tool for executing any custom commands with shorter, more memorable names.

## Folder Structure

FabRun uses a straightforward folder structure to organize your commands. Each folder name corresponds to the command you will use with FabRun. Inside each folder, the command.md file contains the actual command to be executed.

Here is an example of how you can organize your commands:
```
~/.config/fabrun/commands/
├── git_commit/
│   └── command.md
├── deploy/
│   └── command.md
├── backup/
│   └── command.md
└── diff/
    └── command.md
```

### In this structure:
* Each folder name (git_commit, deploy, backup, diff) corresponds to the command you will use with FabRun.
* Inside each folder, the command.md file contains the actual command to be executed.


## Purpose

The primary purpose of FabRun is to make it easier to run lengthy and intricate [Fabric](https://github.com/danielmiessler/fabric) commands. By storing these commands in a structured way, FabRun allows you to execute them effortlessly from the command line, reducing the need to remember and type out long command sequences. Additionally, FabRun can be used to manage and execute any custom command you need, providing a more convenient way to run your scripts.

## Installation

### Via Homebrew

1. Add the tap:
```sh
brew tap HTA86/fabrun
```

2. Install FabRun:
```sh
brew install fabrun
```

3. Verify installation:
```sh
fabrun --version
```

## Update

### Via Homebrew

1. Update:
```sh
brew update
```

### Update Commands
To update the commands from the GitHub repository, use the --update-commands flag:
```sh
fabrun --update-commands
```




## Usage

#### Copy all Commands
```sh
cp -r commands/ ~/.config/fabrun/commands/
```

#### Listing Available Commands
```sh
fabrun --list
```

#### Running a Command
```sh
fabrun <command_name>
```

This command runs the command defined in `~/.config/fabrun/commands/diff/command.md`.


#### Create a custom Command

Store your commands in command.md files within folders. The folder name will be used to identify and run the corresponding command.

#### Example

1.	Create a folder and a command.md file for your command:

```sh
mkdir -p ~/.config/fabrun/commands/diff
echo "fabric --stream --pattern create_git_diff_commit" > ~/.config/fabrun/commands/diff/command.md
```

2.	Run the command using FabRun:

```sh
fabrun diff
```


## Develop
#### Build
```sh
# Bygg för macOS på Intel/AMD64
GOOS=darwin GOARCH=amd64 go build -o releases/fabrun-darwin-amd64

# Bygg för macOS på M1 (ARM64)
GOOS=darwin GOARCH=arm64 go build -o releases/fabrun-darwin-arm64

# Bygg för Linux på AMD64
GOOS=linux GOARCH=amd64 go build -o releases/fabrun-linux-amd64

# Bygg för Windows på AMD64
GOOS=windows GOARCH=amd64 go build -o releases/fabrun-windows-amd64.exe
```

#### Use build script
Make the script executable and run it:
```sh
chmod +x build.sh
./build.sh
```
