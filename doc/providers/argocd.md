# ArgoCD project

This provider discovers projects to fetch by listing ArgoCD applications and by extracting the Spec.Source.RepoURL values from the applications.

Example:

```yaml
  argocd:
    url: "your.argocd.url"
    projects:
      - my-project
    selector: "app=foo"
    base-path: "/tmp/base"
```

You should also set the `ARGOCD_AUTH_TOKEN` to authenticate against ArgoCD.

In this example, standards Insights will fetch ArgoCD projects from the ArgoCD API available at `your.argocd.url` by listing all applications in the ArgoCD project `my-project`. Applications will also be filtered thanks to the selector, and Standards Insights will clone them on the `base-path` directory (the final project path will be `/tmp/base/<project-name>`).
