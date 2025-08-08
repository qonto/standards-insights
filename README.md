# Standards Insights

Standards Insights allows you to follow your applications quality by executing rules (assertions) on them. Standards Insights can be run as a simple CLI or in a server mode where rules are periodically evaluated on projects and results exposed as Prometheus metrics that you can then used to build powerful reporting dashboards.

It also supports automatic discovery of projects to check. For example, if you're deploying your projects using [ArgoCD](https://argo-cd.readthedocs.io/en/stable/), Standards Insights can automatically discover them.

![Grafana Dashboard showing Standards Insights metrics](doc/img/grafana.png?raw=true "Grafana dashboard with example data")

## Quick start

### Installation

- Static binaries are available in the [Github release page](https://github.com/qonto/standards-insights/releases).
- A Docker Image is available on [ECR](https://gallery.ecr.aws/qonto/standards-insights).
- An Helm chart is available on [ECR](https://gallery.ecr.aws/qonto/standards-insights-chart). You should use `ociA://public.ecr.aws/qonto` as repository in Helm in order to fetch it. The supported values are listed in the `chart/values.yaml` file.

### Release

To create a new project release, just create and push a tag:

```
git tag -a <version>
git push --tags
```

### Configuration

Standards Insights is configured through a yaml configuration file. The `config.example.yaml` file is a commented example of the available options.

### CLI mode

Create a new `config.yaml` file containing:

```yaml
interval: 30

labels:
  - team

rules:
  - name: golang-project
    files:
      - path: "go.mod"
        exists: true
  - name: go-version-latest
    files:
      - path: "go.mod"
        contains: "go 1.21"
  - name: go-main
    files:
      - path: "main.go"
        contains: "^(func main)"

checks:
  - name: go-version-latest
    operator: and
    labels:
      category: upgrade
      severity: critical
    rules:
      - go-version-latest
  - name: go-main
    operator: and
    labels:
      category: code-structure
      severity: minor
    rules:
      - go-main

groups:
  - name: golang
    checks:
      - go-version-latest
      - go-main
    when:
      - golang-project
```

We defined in this configuration 3 `rules`. They are assertions that will be validated later on your projects. In this example, the 3 rules are using the `files` module:

- `golang-project`: verify that the file `go.mod` exists
- `go-version-latest`: verify that the `go.mod` file contains the string `go 1.21`
- `go-main`: verify that the file `main.go` contains the pattern `^(func main)`

Rules modules are documented in the `doc/rules` directory.

We then define `checks` in the configuration. Checks are used to group similar rules together. They can also be configured with optional labels (`category` and `severity` in this example) that we will used later, when Standards Insights is run in `daemon` mode.
In this example, the two checks contain only one rule each (`go-version-latest` and `go-main`).
Checks supports 2 different operators to evaluate `rules`: `and` and `or`.
`and` is the default operator and the check result will be successful if all the rules evaluate to true.
`or` operator will mark the check as successful if at least one rule evaluates to true.

Finally, we configure `groups`: groups contains checks to execute on a project when a condition is true. In this example, the group named `golang` will only be executed on projects where `golang-project` is valid (so, when the `go.mod` file exists). Thanks to the `when` condition, we can specify which check should be executed on which project.

Normally groups and checks within it are executed at project level. It is possible to configure a group to be executed at file level. In order to so, you have to set the `apply-to-files` to true under the `files` block as show below:

```yaml
  - name: paseto-group
    files:
      apply-to-files: true
      files-pattern: "controller.rb"
    checks:
      - paseto-migration-check
    when:
      - ruby-controller
      - ruby-projects
```

When `apply-to-files` is set to true all the files of a project are evaluated. But if you want to evaluate only some files, you can do so by setting `apply-to-files` to true and `files-pattern` to the pattern of the files you want to evaluate. The configuration example above will execute the `paseto-migration-check` check on all the files matching the `files-pattern` pattern. Therefore there will be one check result per file.
Only for certain rules makes sense to apply them at file level. Check the `doc/rules` for more information. For assigning ownership to a file, we support [CODEOWNERS](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners) file parsing. If the project has a CODEOWNERS file, the ownership will be extracted and added to the file labels. If it's not the case, the owner will be extracted as usual according to the provider.

If you navigate in a Golang project (for example, the Standards Insights repository itself) and run `standards-insights run --config config.yaml`, you will get:

```
== Project standards-insights
✅ Check go-version-latest PASS (labels: map[category:upgrade severity:critical])
✅ Check go-main PASS (labels: map[category:code-structure severity:minor])
```

We see that both checks were successfully executed on this project.

You can also fetch a remote configuration file by passing to `--config` an HTTP url. In server mode, Standards Insights can expose its configuration file on `/config`.

### Server mode

In server mode, Standards Insights will periodically:

- Clone or pull the configured repositories from a Git server.
- Execute the checks on them if they matches the groups `when` clause.
- Report the checks results in Prometheus format on `/metrics`.

In order to discover projects to clone, you need to configure in your configuration file `providers`. Standards Insights supports a `static` provider and an `argocd` provider. Check the `doc/providers` for more information about providers. Let's use the static one for this example.

Add in your configuration these values:

```yaml
http:
  host: 127.0.0.1
  port: 3000

providers:
  static:
    - name: "standards-insights"
      url: https://github.com/qonto/standards-insights.git
      branch: "main"
      path: "/tmp/projects/standards-insights"
      labels:
        team: sre

# git:
#   basic-auth:
#     username: oauth2
#     password: pass # You can use GIT_BASIC_AUTH_PASSWORD as well
#   private-key:
#     path: /home/mcorbin/.ssh/my_ssh_private_key
#     password: key-password # you can use GIT_PRIVATE_KEY_PASSWORD

interval: 300

labels:
  - category
  - severity

```

- The `http` block is the http server configuration.
- The `provider` block is used to configure providers. The `static` provider allows you to statically reference Git. projects to fetch and to configure where they should be cloned. Optional labels can be passed per project.
- The `git` block can be used to configure Git authentication mechanisms.
- The `interval` configuration is the interval in seconds between Standards Insights check loop.
- The `labels` block contains a list of labels to extract from checks and projects and to use in the generated Prometheus metrics.

In summary, Standards Insights will with this configuration, every 300 seconds,  clone the `https://github.com/qonto/standards-insights.git` repository and execute the configured checks on it.

Let's run the server with `standards-insights server --config config.yaml`. You should see in the logs:

```
time=2023-10-05T13:28:34.204+02:00 level=INFO msg="starting HTTP server on 127.0.0.1:3000"
time=2023-10-05T13:28:34.204+02:00 level=INFO msg="starting daemon"
time=2023-10-05T13:28:34.204+02:00 level=INFO msg="checking projects"
time=2023-10-05T13:28:34.886+02:00 level=INFO msg="checking project standards-insights"
time=2023-10-05T13:28:34.886+02:00 level=INFO msg="projects checked"
```

You can now open `http://localhost:3000/metrics` to see the application metrics. You will see alongside other metrics:

```
check_result_success{category="code-structure",name="go-main",project="standards-insights",severity="minor",team="sre"} 1
check_result_success{category="upgrade",name="go-version-latest",project="standards-insights",severity="critical",team="sre"} 1
```

Both checks were successful, so the value of the metric is `1`. In case of failure, the value would be `0`.

