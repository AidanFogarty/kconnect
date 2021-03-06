## kconnect ls

Query the user's connection history

### Synopsis


Query and display the user's connection history entries, including entry IDs and
aliases.

Each time kconnect creates a new kubectl context to connect to a Kubernetes
cluster, it saves the settings for the new connection as an entry in the user's
connection history.  The user can then reconnect using those same settings later
via the connection history entry's ID or alias.


```bash
kconnect ls [flags]
```

### Examples

```bash

  # Display all connection history entries as a table
  kconnect ls

  # Display all connection history entries as YAML
  kconnect ls --output yaml

  # Display a specific connection history entry by entry id
  kconnect ls --filter id=01EM615GB2YX3C6WZ9MCWBDWBF

  # Display a specific connection history entry by its alias
  kconnect ls --filter alias=mydev

  # Display all connection history entries that have "dev" in its alias
  kconnect ls --filter alias=*dev*

  # Display all connection history entries for the EKS managed cluster provider
  kconnect ls --filter cluster-provider=eks

  # Display all connection history entries for entries with namespace kube-system
  kconnect ls --filter namespace=kube-system

  # Reconnect using the connection history entry alias
  kconnect to mydev

```

### Options

```bash
      --filter string             filter to apply to import. Can specify multiple filters by using commas, and supports wilcards (*)
  -h, --help                      help for ls
      --history-location string   Location of where the history is stored. (default "$HOME/.kconnect/history.yaml")
  -k, --kubeconfig string         Location of the kubeconfig to use. (default "$HOME/.kube/config")
      --max-history int           Sets the maximum number of history items to keep (default 100)
      --no-history                If set to true then no history entry will be written
  -o, --output string             Output format for the results (default "table")
```

### Options inherited from parent commands

```bash
      --config string      Configuration file for application wide defaults. (default "$HOME/.kconnect/config.yaml")
      --no-input           Explicitly disable interactivity when running in a terminal
      --no-version-check   If set to true kconnect will not check for a newer version
  -v, --verbosity int      Sets the logging verbosity. Greater than 0 is debug and greater than 9 is trace.
```

### SEE ALSO

* [kconnect](index.md)	 - The Kubernetes Connection Manager CLI


> NOTE: this page is auto-generated from the cobra commands
