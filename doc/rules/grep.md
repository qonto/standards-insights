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

- `path`: the file or directory to check. The path is not mandatory when the rule is used within a check of a group that is applied to files. In this case, the rule will be applied to all the project files that match the pattern.
- `pattern`: the string pattern to search for (mandatory)
- `exclude-dir`: the string pattern representing command-line directories to skip
- `recursive`: recursively search subdirectories listed (grep `-r` option)
- `match`: if set to true, the module will be successful if the pattern is found. Default to false, where the module will be successful if the pattern is **not** found
- `extended-regexp`: if set to true, the pattern is treated as an extended regular expression (grep `-E` option)
- `null-data`: if set to true, treat input data as sequences of lines terminated by a zero-byte instead of a newline. Useful with multi-line matching. (grep `-z` option)
- `insensitive-case`: if set to true, perform case-insensitive matching (grep `-i` option)
