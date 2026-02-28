# Brand Assets

BrandKit includes a library of 55+ brand icons with standardized SVG variants.

## File Conventions

Each brand directory contains up to 3 standardized variants:

| File | Description |
|------|-------------|
| `icon_orig.svg` | Original source icon (as obtained) |
| `icon_white.svg` | White on transparent background |
| `icon_color.svg` | Color variant on transparent background |

## Available Brands

### AI & ML

| Brand | Directory | Variants |
|-------|-----------|----------|
| Anthropic | `anthropic/` | orig, white, color |
| Anthropic Claude | `anthropic-claude/` | orig, white, color |
| Deepgram | `deepgram/` | orig, white |
| ElevenLabs | `elevenlabs/` | orig, white, color |
| Google Gemini | `google-gemini/` | orig, white, color |
| Langfuse | `langfuse/` | orig, white, color |
| Ollama | `ollama/` | orig, white, color |
| OpenAI | `openai/` | orig, white, color |
| Opik | `opik/` | orig, white, color |
| Phoenix | `phoenix/` | orig, white, color |
| xAI | `xai/` | orig, white, color |

### Cloud & Infrastructure

| Brand | Directory | Variants |
|-------|-----------|----------|
| AWS | `aws/` | orig, white, color |
| AWS AgentCore | `aws-agentcore/` | orig, white |
| AWS CDK | `aws-cdk/` | orig, white, color |
| AWS Kiro | `aws-kiro/` | orig, white, color |
| AWS Security Lake | `aws-security-lake/` | orig, white, color |
| Azure | `azure/` | orig, white, color |
| Docker | `docker/` | orig, white, color |
| Google GCP | `google-gcp/` | orig, white, color |
| Helm | `helm/` | orig, white, color |
| Kubernetes | `kubernetes/` | orig, white, color |
| Pulumi | `pulumi/` | orig, white, color |

### Developer Tools

| Brand | Directory | Variants |
|-------|-----------|----------|
| Bolt | `bolt/` | orig, white, color |
| Cursor | `cursor/` | orig, white, color |
| GitHub | `github/` | orig, white, color |
| GitHub Copilot | `github-copilot/` | orig, white, color |
| Lovable | `lovable/` | orig, white, color |
| Postman | `postman/` | orig, white, color |
| Replit | `replit/` | orig, white, color |
| v0 | `v0/` | orig, white, color |
| Windsurf | `windsurf/` | orig, white, color |

### Programming Languages

| Brand | Directory | Variants |
|-------|-----------|----------|
| Dart | `dart/` | orig, white, color |
| Flutter | `flutter/` | orig, white, color |
| Go | `go/` | orig, white, color |
| JavaScript | `javascript/` | orig, white, color |
| Kotlin | `kotlin/` | orig, white, color |
| Python | `python/` | orig, white, color |

### Frameworks & Libraries

| Brand | Directory | Variants |
|-------|-----------|----------|
| Bootstrap | `bootstrap/` | orig, white, color |
| React | `react/` | orig, white, color |
| Spring | `spring/` | orig, white, color |

### Security & Auth

| Brand | Directory | Variants |
|-------|-----------|----------|
| Cedar | `cedar/` | orig, white, color |
| CrowdStrike | `crowdstrike/` | orig, white, color |
| OAuth2 | `oauth2/` | orig, white, color |
| OCSF | `ocsf/` | orig, white, color |
| Saviynt | `saviynt/` | orig, white, color |

### Databases

| Brand | Directory | Variants |
|-------|-----------|----------|
| Datadog | `datadog/` | orig, white, color |
| PostgreSQL | `postgresql/` | orig, white, color |

### APIs & Services

| Brand | Directory | Variants |
|-------|-----------|----------|
| OpenAPI | `openapi/` | orig, white, color |
| SerpAPI | `serpapi/` | orig, white, color |
| Serper | `serper/` | orig, white, color |
| Twilio | `twilio/` | orig, white, color |

### Operating Systems

| Brand | Directory | Variants |
|-------|-----------|----------|
| Linux | `linux/` | orig, white, color |
| macOS | `macos/` | orig, white, color |
| Windows | `windows/` | orig, white, color |

### Communication

| Brand | Directory | Variants |
|-------|-----------|----------|
| WhatsApp | `whatsapp/` | orig, white, color |

## Using Brand Assets

### CLI

```bash
# Create white variant from original
brandkit white brands/react/icon_orig.svg -o brands/react/icon_white.svg

# Create color variant from original
brandkit color brands/react/icon_orig.svg -o brands/react/icon_color.svg
```

### Go Library

```go
import "github.com/grokify/brandkit"

// Get embedded icon
icon, err := brandkit.GetIconWhite("react")
if err != nil {
    log.Fatal(err)
}

// List available icons
icons := brandkit.ListIcons()
for _, name := range icons {
    fmt.Println(name)
}

// Check if icon exists
if brandkit.IconExists("react") {
    fmt.Println("React icon available")
}
```

## Quality Standards

All brand icons meet these standards:

- **Pure Vector** — No embedded raster data
- **Security Scanned** — No malicious elements
- **Centered** — Content centered in viewBox
- **Transparent Background** — No solid backgrounds
- **Consistent Sizing** — Optimized viewBox

## Contributing

To add a new brand:

1. Add original SVG to `brands/<name>/icon_orig.svg`
2. Generate variants:
   ```bash
   brandkit white brands/<name>/icon_orig.svg -o brands/<name>/icon_white.svg
   brandkit color brands/<name>/icon_orig.svg -o brands/<name>/icon_color.svg
   ```
3. Verify all files:
   ```bash
   brandkit verify brands/<name>/
   brandkit security-scan brands/<name>/
   ```
