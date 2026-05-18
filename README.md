# deploylog

Aggregates deployment events from multiple CI/CD sources into a unified audit timeline.

## Installation

```bash
go install github.com/yourusername/deploylog@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/deploylog.git && cd deploylog && go build ./...
```

## Usage

Configure your CI/CD sources in a `deploylog.yaml` file:

```yaml
sources:
  - type: github_actions
    token: $GITHUB_TOKEN
    repo: myorg/myapp
  - type: gitlab_ci
    token: $GITLAB_TOKEN
    project_id: 12345
```

Then run:

```bash
deploylog tail --since 24h
```

Example output:

```
2024-06-10 14:32:01  [github_actions]  myorg/myapp         DEPLOY  production  ✓ success  (triggered by: alice)
2024-06-10-13:15:44  [gitlab_ci]       myorg/backend       DEPLOY  staging     ✓ success  (triggered by: bob)
2024-06-10 11:02:19  [github_actions]  myorg/frontend      DEPLOY  production  ✗ failed   (triggered by: carol)
```

Export the timeline to JSON for auditing:

```bash
deploylog export --since 7d --format json > audit.json
```

## Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `--since` | `24h` | How far back to fetch events |
| `--format` | `text` | Output format: `text`, `json`, `csv` |
| `--source` | all | Filter by a specific source type |

## License

MIT — see [LICENSE](LICENSE) for details.