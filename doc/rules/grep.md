# Grep

This module allows you to configure assertions on file or directories content using the `grep` command line tool.
You need the `grep` executable in your path for this module to work.

Example:

```yaml
rules:
  - name: grepp-example
    grep:
      - path: "/tmp"
        pattern: "hello*"
        recursive: true
        match: true
```

Available options are:

- `path`: the file or directory to check (mandatory)
- `pattern`: the string pattern to search for (mandatory)
- `recursive`: recursively search subdirectories listed (grep `-r`option)
- `match`: if set to true, the module will be successful if the pattern is found. Default to false, where the module will be successful if the pattern is **not** found.
