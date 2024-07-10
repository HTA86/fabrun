import sys
import subprocess
import os

PATTERNS_DIR = os.path.expanduser('~/.config/fabric/patterns')

def load_command(command_name):
    command_file = os.path.join(PATTERNS_DIR, command_name, 'command.md')
    if not os.path.isfile(command_file):
        return None
    
    with open(command_file, 'r') as file:
        command = file.read().strip()
    return command

def run_command(command, command_name):
    full_command = f"{command} | fabric --stream --pattern {command_name}"
    try:
        result = subprocess.run(full_command, shell=True, check=True, text=True)
        print(result.stdout)
    except subprocess.CalledProcessError as e:
        print(f"Error running command '{full_command}': {e.stderr}")

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: fabrun <command-name>")
        sys.exit(1)
    
    command_name = sys.argv[1]
    command = load_command(command_name)

    if command:
        run_command(command, command_name)
    else:
        print(f"Command '{command_name}' not found in {PATTERNS_DIR}")