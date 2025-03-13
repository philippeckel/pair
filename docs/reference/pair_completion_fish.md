# pair completion fish

Generate the autocompletion script for fish

## Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	pair completion fish | source

To load completions for every new session, execute once:

	pair completion fish > ~/.config/fish/completions/pair.fish

You will need to start a new shell for this setup to take effect.


```shell
pair completion fish [flags]
```

## Options

```text
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

## Options inherited from parent commands

```text
  -c, --config string   config file path (default "/Users/philipp.eckel/.pair.json")
```

## See also

* [pair completion](pair_completion.md) - Generate the autocompletion script for the specified shell
