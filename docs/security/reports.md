# Security Reports

BrandKit generates JSON security reports following the [multi-agent-spec](https://github.com/agentplexus/multi-agent-spec) team-report format.

## Generating Reports

### CLI

```bash
brandkit security-scan-all brands/ --report=security-report.json
```

With project metadata:

```bash
brandkit security-scan-all brands/ \
  --report=security-report.json \
  --project=myproject \
  --version=1.0.0
```

### Library

```go
import "github.com/grokify/brandkit/svg/security"

results, _ := security.DirectoryRecursive("brands/")
report := security.GenerateReport(results, "myproject", "1.0.0")

jsonBytes, _ := report.ToJSON()
os.WriteFile("report.json", jsonBytes, 0644)
```

## Report Structure

```json
{
  "$schema": "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/report/team-report.schema.json",
  "title": "SVG SECURITY SCAN REPORT",
  "project": "brandkit",
  "version": "0.4.0",
  "phase": "SECURITY VALIDATION",
  "status": "GO",
  "generated_at": "2024-01-15T10:30:00Z",
  "generated_by": "brandkit security-scan",
  "summary_blocks": [...],
  "teams": [...],
  "footer_blocks": [...]
}
```

## Status Values

| Status | Meaning |
|--------|---------|
| `GO` | No threats detected, safe to proceed |
| `NO-GO` | Critical or high severity threats found |
| `WARN` | Medium or low severity threats found |
| `SKIP` | Scan was skipped |

## Summary Block

The summary block provides high-level statistics:

```json
{
  "summary_blocks": [
    {
      "type": "kv_pairs",
      "pairs": [
        {"key": "Files Scanned", "value": "166"},
        {"key": "Secure Files", "value": "138"},
        {"key": "Files with Threats", "value": "28"},
        {"key": "Total Threats", "value": "28"}
      ]
    }
  ]
}
```

## Team Sections

Each threat category is represented as a "team" section:

```json
{
  "teams": [
    {
      "id": "script-detection",
      "name": "Script Detection",
      "status": "GO",
      "tasks": [
        {
          "id": "scan",
          "status": "GO",
          "detail": "No threats detected"
        }
      ]
    },
    {
      "id": "event-handler-detection",
      "name": "Event Handler Detection",
      "status": "NO-GO",
      "tasks": [
        {
          "id": "scan",
          "status": "NO-GO",
          "severity": "critical",
          "detail": "3 threat(s) detected"
        }
      ],
      "content_blocks": [
        {
          "type": "list",
          "title": "Findings",
          "items": [
            {
              "icon": "🔴",
              "text": "brands/malicious/icon.svg: event handler attribute"
            }
          ]
        }
      ]
    }
  ]
}
```

## Team Categories

| Team ID | Name | Severity |
|---------|------|----------|
| `script-detection` | Script Detection | critical |
| `event-handler-detection` | Event Handler Detection | critical |
| `external-ref-detection` | External Reference Detection | high |
| `xml-entity-detection` | XML Entity Detection | high |
| `animation-detection` | Animation Detection | medium |
| `style-block-detection` | Style Block Detection | low |
| `link-detection` | Link Detection | medium |

## Footer Block (Action Items)

When threats are found, the footer contains recommended actions:

```json
{
  "footer_blocks": [
    {
      "type": "kv_pairs",
      "title": "ACTION ITEMS",
      "pairs": [
        {
          "icon": "🔴",
          "key": "1",
          "value": "Remove all script elements and event handlers (CRITICAL)"
        },
        {
          "icon": "🟡",
          "key": "2",
          "value": "Remove DOCTYPE and ENTITY declarations (HIGH)"
        }
      ]
    }
  ]
}
```

## Using Reports in CI

### GitHub Actions

```yaml
- name: Security scan
  run: brandkit security-scan-all brands/ --report=security-report.json

- name: Upload report
  uses: actions/upload-artifact@v3
  if: always()
  with:
    name: security-report
    path: security-report.json

- name: Check status
  run: |
    STATUS=$(jq -r '.status' security-report.json)
    if [ "$STATUS" = "NO-GO" ]; then
      echo "Security scan failed with status: $STATUS"
      exit 1
    fi
```

### Parsing with jq

```bash
# Get overall status
jq -r '.status' security-report.json

# Count total threats
jq -r '.summary_blocks[0].pairs[3].value' security-report.json

# List all failing categories
jq -r '.teams[] | select(.status != "GO") | .name' security-report.json

# Get all findings
jq -r '.teams[].content_blocks[]?.items[]?.text' security-report.json
```

## Report Schema

Reports conform to the multi-agent-spec team-report schema:

```
https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/report/team-report.schema.json
```

This enables integration with tools that understand the multi-agent-spec format.
