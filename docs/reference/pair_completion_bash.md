# pair completion bash

Generate the autocompletion script for bash

## Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(pair completion bash)

To load completions for every new session, execute once:

#### Linux:

	pair completion bash > /etc/bash_completion.d/pair

#### macOS:

	pair completion bash > $(brew --prefix)/etc/bash_completion.d/pair

You will need to start a new shell for this setup to take effect.


```shell
pair completion bash
```

## Options

```text
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

## Options inherited from parent commands

```text
  -c, --config string   config file path (default "/Users/philipp.eckel/.pair.json")
```

## See also

* [pair completion](pair_completion.md) - Generate the autocompletion script for the specified shell
