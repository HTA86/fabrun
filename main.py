import argparse
import subprocess
from pathlib import Path
import logging

# Setup logging
logging.basicConfig(level=logging.INFO, format='%(message)s')
logger = logging.getLogger(__name__)

VERSION = "0.2.0"
COMMANDS_DIR = Path.home() / '.config' / 'fabrun' / 'commands'
DESCRIPTION_FILENAME = 'about.md'
COMMAND_FILENAME = 'command.md'

def load_command(command_name: str) -> str:
    """
    Load the command from a specified file.

    :param command_name: Name of the command file (without extension).
    :return: The command as a string, or None if the file does not exist or is empty.
    """
    command_file = COMMANDS_DIR / command_name / COMMAND_FILENAME
    
    if not command_file.is_file():
        logger.error(f"Error: Command file '{command_file}' does not exist.")
        return None

    try:
        command = command_file.read_text().strip()
        if not command:
            logger.error(f"Error: Command file '{command_file}' is empty.")
            return None
        return command
    except Exception as e:
        logger.error(f"Error reading command file '{command_file}': {e}")
        return None

def load_description(command_name: str) -> str:
    """
    Load the description from a specified file.

    :param command_name: Name of the description file (without extension).
    :return: The description as a string, or a default message if the file does not exist.
    """
    description_file = COMMANDS_DIR / command_name / DESCRIPTION_FILENAME
    
    if not description_file.is_file():
        return "No description available."

    try:
        description = description_file.read_text().strip()
        if not description:
            return "No description available."
        return description
    except Exception as e:
        logger.error(f"Error reading description file '{description_file}': {e}")
        return "Error reading description."

def run_command(command: str):
    """
    Execute the given command using the system shell and stream the output.

    :param command: The command to execute.
    """
    if not command:
        logger.error("No command to run.")
        return

    try:
        result = subprocess.run(command, shell=True, check=True, text=True)
        if result.returncode != 0:
            logger.error(f"Error running command '{command}': {result.stderr}")
    except Exception as e:
        logger.error(f"Unexpected error running command '{command}': {e}")

def list_commands():
    """
    List all available commands in the COMMANDS_DIR.
    
    :return: List of available command names and their descriptions.
    """
    if not COMMANDS_DIR.is_dir():
        logger.error(f"Error: Commands directory '{COMMANDS_DIR}' does not exist.")
        return []

    try:
        return [d.name for d in COMMANDS_DIR.iterdir() if d.is_dir()]
    except Exception as e:
        logger.error(f"Error listing commands in '{COMMANDS_DIR}': {e}")
        return []

def format_command_list(commands):
    """
    Format the command list with descriptions for better readability.
    
    :param commands: List of command names.
    """
    max_command_length = max(len(cmd) for cmd in commands)
    max_width = 80
    separator = '-' * max_width
    for cmd in commands:
        description = load_description(cmd)
        line_length = max_command_length + len(description) + 4
        if line_length > max_width:
            description = description[:max_width - max_command_length - 7] + '...'
        print(f" - \033[1m{cmd.ljust(max_command_length)}\033[0m\t\033[92m{description}\033[0m")
        print(f"{separator}")

def main():
    """
    Main function to parse arguments and run the specified command.
    """
    parser = argparse.ArgumentParser(
        prog='fabrun',
        description="Execute predefined commands from 'command.md' files based on folder names.",
        epilog="Example usage:\n"
               "  fabrun git_commit           # Run the git_commit command\n"
               "  fabrun --list               # List all available commands\n"
               "  fabrun --help               # Show this help message",
        formatter_class=argparse.RawTextHelpFormatter
    )
    parser.add_argument('command_name', type=str, nargs='?', help="The name of the folder containing the command.md file to run.")
    parser.add_argument('--list', '-l', action='store_true', help="List all available commands.")
    parser.add_argument('--version', '-v', action='version', version=f'%(prog)s {VERSION}', help="Show the program's version.")
    args = parser.parse_args()

    if args.list:
        commands = list_commands()
        if commands:
            print("Available commands:")
            format_command_list(commands)
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
    COMMANDS_DIR.mkdir(parents=True, exist_ok=True)
    main()