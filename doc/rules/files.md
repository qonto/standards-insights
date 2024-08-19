# Files

This module allows you to configure assertions on files.

Example:

```yaml
rules:
  - name: golang-projects
    files:
      - path: "go.mod"
        exists: true
```

Available options are:

- `path`: the file path to check (mandatory). Notice: since the path is mandatory, it does not make sense to use this rule within a check used in a group that is applied to files.
- `exists`: should the file exists or not.
- `contains`: a regular expression that will be checked against the file content. The rule will be successful if the regular expression find a match.
- `not-contains`: a regular expression that will be checked against the file content. The rule will be successful if the regular expression doesn't find a match.
