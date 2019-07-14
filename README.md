# Kubernetes Job UI

A web ui for triggering manual runs of Kubernetes CronJobs.  Supports a configuration annotation that can be added to CronJobs that allows parameters to be set for the job run.  Previous job runs can be viewed in the UI.

## Configuration

```json
{
  "options": [
    {
      "envvar": "FOOBAR",
      "values": [
        "FOO",
        "BAR",
        "foobar"
      ],
      "default": "foobar",
      "description": "An option to select what type of foo"
    },
    {
      "envvar": "VERY_LONG_BOOLEAN_VALUE",
      "values": [
        "true",
        "false"
      ],
      "default": true,
      "description": "Should we do something?"
    }
  ]
}
```

## Limitations

1. Includes a known race condition.  Future versions will fix this issue.
2. There are not currently any ways to limit which CronJobs are visible in the ui.
3. Logs are not streamed to ui.  Page needs to be manually refreshed.
