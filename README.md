# Pair - Git co-authors management tool

## Overviewz

Pair is a command-line tool that simplifies managing Git co-authors for pair programming sessions. It allows you to maintain a list of co-authors and easily add them to your Git commits, helping you give proper credit to everyone contributing to your code.

## Installation

## Features

* Maintain a roster of frequent collaborators
* Add co-authors to commits with simple commands
* Interactive fuzzy-search selection of co-authors
* View active co-authors at any time
* Works with Git's commit template mechanism
* Supports both global and project-specific co-author lists

## Usage

```shell
# List all available co-authors
pair list

# Show currently active co-authors
pair show

# Add co-authors by alias or index
pair add jane john

# Remove a co-author
pair remove john

# Interactively select co-authors
pair select

# Interactively remove co-authors
pair unselect

# Clear all co-authors
pair clear

# Initialize with sample config
pair init
```
