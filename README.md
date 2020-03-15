# Kubernetes Job UI

A web ui for triggering manual runs of Kubernetes CronJobs.  Supports a configuration annotation that can be added to CronJobs that allows parameters to be set for the job run.  Previous job runs can be viewed in the UI.

# Configuration

## Process configuration

`NAMESPACE` Specify which namespace to return jobs from. Defaults to all namespaces
`CONFIGURED_ONLY` Whether to only return jobs with a configuration block

## Job configuration block

A job can be configured by adding a configuration block as an annotation to the cronjob

```yaml
  annotations:
    kubernetes-job-runner.io/config: |
      shortDescription: Deploy the current k8s-manifests configuration of the specified helm chart
      folder: folder2
      longDescription: |
        Choose a helm chart and specify a version from github to deploy.
        Versions are validated at runtime to ensure they exist.
        Third line to make this longer and longer.
        Fourth line to ensure it's longer.
      options:
        - envVar: LIST
          type: list
          values:
            - FOO
            - BAR
            - foobar
          default: foobar
          description: An option to select what type of foo
          container: hello
        - envVar: INT
          default: '1'
          type: int
          description: int to deploy
          container: goodbye
        - envVar: FLOAT
          default: '1.1'
          type: float
          descriptions: float to deploy
          container: goodbye
        - envVar: REGEX
          default: ''
          type: regex
          regex: "^int-.*"
          description: regex to deploy
          container: goodbye
        - envVar: TEXTAREA
          default: ''
          type: textarea
          description: freeform text
          container: goodbye
        - envVar: BOOL
          type: bool
          values:
            - 'true'
            - 'false'
          default: 'true'
          description: Should we do something?
          container: hello
```

## Limitations

1. There is no authentication so anyone with access to the web page can trigger jobs.
2. The backend currently polls for updates to jobs rather than listening for events.
