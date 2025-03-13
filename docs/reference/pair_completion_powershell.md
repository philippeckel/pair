# pair completion powershell

Generate the autocompletion script for powershell

## Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	pair completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```shell
pair completion powershell [flags]
```

## Options

```text
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

## Options inherited from parent commands

```text
  -c, --config string   config file path (default "/Users/philipp.eckel/.pair.json")
```

## See also

* [pair completion](pair_completion.md) - Generate the autocompletion script for the specified shell
