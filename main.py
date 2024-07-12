import argparse
import subprocess
import os

COMMANDS_DIR = os.path.expanduser('~/.config/fabrun/commands')

def load_command(command_name):
    """
    Load the command from a specified file.

    :param command_name: Name of the command file (without extension).
    :return: The command as a string, or None if the file does not exist or is empty.
    """
    command_file = os.path.join(COMMANDS_DIR, command_name, 'command.md')
    
    if not os.path.isfile(command_file):
        print(f"Command file '{command_file}' does not exist.")
        return None

    try:
        with open(command_file, 'r') as file:
            command = file.read().strip()
        if not command:
            print(f"Command file '{command_file}' is empty.")
            return None
        return command
    except Exception as e:
        print(f"Error reading command file '{command_file}': {e}")
        return None

def run_command(command):
    """
    Execute the given command using the system shell.

    :param command: The command to execute.
    """
    if not command:
        print("No command to run.")
        return

    try:
        result = subprocess.run(command, shell=True, check=True, text=True, capture_output=True)
        print(result.stdout)
    except subprocess.CalledProcessError as e:
        print(f"Error running command '{command}': {e.stderr}")
    except Exception as e:
        print(f"Unexpected error running command '{command}': {e}")

def main():
    """
    Main function to parse arguments and run the specified command.
    """
    parser = argparse.ArgumentParser(description="Run a command from a file.")
    parser.add_argument('command_name', type=str, help="The name of the command to run.")
    args = parser.parse_args()

    command = load_command(args.command_name)

    if command:
        run_command(command)
    else:
        print(f"Command '{args.command_name}' not found or invalid in {COMMANDS_DIR}")

if __name__ == "__main__":
    os.makedirs(COMMANDS_DIR, exist_ok=True)
    main()