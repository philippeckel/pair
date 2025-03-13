# pair completion zsh

Generate the autocompletion script for zsh

## Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(pair completion zsh)

To load completions for every new session, execute once:

#### Linux:

	pair completion zsh > "${fpath[1]}/_pair"

#### macOS:

	pair completion zsh > $(brew --prefix)/share/zsh/site-functions/_pair

You will need to start a new shell for this setup to take effect.


```shell
pair completion zsh [flags]
```

## Options

```text
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

## Options inherited from parent commands

```text
  -c, --config string   config file path (default "/Users/philipp.eckel/.pair.json")
```

## See also

* [pair completion](pair_completion.md) - Generate the autocompletion script for the specified shell
