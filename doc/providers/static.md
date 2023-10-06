# Static provider

This provider returns a statically configured list of projects.

Example:

```yaml
providers:
  static:
    - name: "standards-insights"
      url: https://github.com/qonto/standards-insights.git
      branch: "main"
      path: "/tmp/projects/standards-insights"
      labels:
        team: sre
```

In this example, the project will be fetch from `https://github.com/qonto/standards-insights.git` on the branch `main` and cloned into the path `/tmp/projects/standards-insights`.
Labels are optional and will be used if configured in the Prometheus metrics exposed by the HTTP server.
