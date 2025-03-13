# Installation

There are several ways to install Pair on your system. Choose the method that works best for your environment.

## Homebrew (macOS and Linux)

The easiest way to install Pair on macOS and Linux is via Homebrew:

```shell
brew install philippeckel/tap/pair
```

To update to the latest version:

```shell
brew upgrade pair
```

## Build from source

If you prefer to build from source or need the latest development version, you can compile Pair yourself:

### Prerequisites

* Go 1.18 or newer
* Git

### Steps

* Clone the repository:

```shell
git clone https://github.com/philippeckel/pair.git && cd pair
```

* Build the binary:

```shell
go build -o pair
```

* Install the binary to your Go binary path, usually `~/go/bin`:

```shell
go install
```

## Verify Installation

After installation, verify that Pair is working correctly:

```shell
pair --version
```

This should display the current version of Pair.
