# Configuration File

Pair supports configuration through a YAML file located in your home directory. This allows you to customize its behavior to match your preferences.

## Configuration file location

Pair looks for configuration in the following locations (in order of priority):

1. `~/.config/pair/config.yaml` (recommended)
2. Current working directory

## Creating Your Configuration File

To set up your configuration:

```shell
# Create the config directory if it doesn't exist
mkdir -p ~/.config/pair

# Create or edit the configuration file
touch ~/.config/pair/config.yaml
```

## Available Configuration Options

```shell
# ~/.config/pair/config.yaml

# Disable colored output (default: false)
no_color: false
```
