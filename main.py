import argparse
import subprocess
import os
import logging

# Setup logging
logging.basicConfig(level=logging.INFO, format='%(message)s')
logger = logging.getLogger(__name__)

COMMANDS_DIR = os.path.expanduser('~/.config/fabrun/commands')

def load_command(command_name):
    """
    Load the command from a specified file.

    :param command_name: Name of the command file (without extension).
    :return: The command as a string, or None if the file does not exist or is empty.
    """
    command_file = os.path.join(COMMANDS_DIR, command_name, 'command.md')
    
    if not os.path.isfile(command_file):
        logger.error(f"Command file '{command_file}' does not exist.")
        return None

    try:
        with open(command_file, 'r') as file:
            command = file.read().strip()
        if not command:
            logger.error(f"Command file '{command_file}' is empty.")
            return None
        return command
    except Exception as e:
        logger.error(f"Error reading command file '{command_file}': {e}")
        return None

def run_command(command):
    """
    Execute the given command using the system shell.

    :param command: The command to execute.
    """
    if not command:
        logger.error("No command to run.")
        return

    try:
        result = subprocess.run(command, shell=True, check=True, text=True, capture_output=True)
        print(result.stdout)
    except subprocess.CalledProcessError as e:
        logger.error(f"Error running command '{command}': {e.stderr}")
    except Exception as e:
        logger.error(f"Unexpected error running command '{command}': {e}")

def list_commands():
    """
    List all available commands in the COMMANDS_DIR.
    """
    if not os.path.isdir(COMMANDS_DIR):
        logger.error(f"Commands directory '{COMMANDS_DIR}' does not exist.")
        return []

    try:
        return [d for d in os.listdir(COMMANDS_DIR) if os.path.isdir(os.path.join(COMMANDS_DIR, d))]
    except Exception as e:
        logger.error(f"Error listing commands in '{COMMANDS_DIR}': {e}")
        return []

def main():
    """
    Main function to parse arguments and run the specified command.
    """
    parser = argparse.ArgumentParser(description="Run a command from a file.")
    parser.add_argument('command_name', type=str, nargs='?', help="The name of the command to run.")
    parser.add_argument('--list', '-l', action='store_true', help="List all available commands.")
    args = parser.parse_args()

    if args.list:
        commands = list_commands()
        if commands:
            print("Available commands:")
            for cmd in commands:
                print(f" - {cmd}")
        else:
            print("No commands available.")
    elif args.command_name:
        command = load_command(args.command_name)
        if command:
            run_command(command)
        else:
            logger.error(f"Command '{args.command_name}' not found or invalid in {COMMANDS_DIR}")
    else:
        parser.print_help()

if __name__ == "__main__":
    os.makedirs(COMMANDS_DIR, exist_ok=True)
    main()