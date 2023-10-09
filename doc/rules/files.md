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

- `path`: the file path to check (mandatory).
- `exists`: should the file exists or not.
- `contains`: a regular expression that will be checked against the file content. The rule will be successful if the regular expression find a match.
- `not-contains`: a regular expression that will be checked against the file content. The rule will be successful if the regular expression doesn't find a match.
