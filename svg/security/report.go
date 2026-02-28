package security

import (
	"encoding/json"
	"fmt"
	"time"
)

// Status represents the Go/No-Go status for a report.
type Status string

const (
	StatusGo   Status = "GO"
	StatusNoGo Status = "NO-GO"
	StatusWarn Status = "WARN"
	StatusSkip Status = "SKIP"
)

// TeamReport represents the full security scan report.
type TeamReport struct {
	Schema        string            `json:"$schema,omitempty"`
	Title         string            `json:"title,omitempty"`
	Project       string            `json:"project"`
	Version       string            `json:"version"`
	Phase         string            `json:"phase"`
	Tags          map[string]string `json:"tags,omitempty"`
	SummaryBlocks []ContentBlock    `json:"summary_blocks,omitempty"`
	Teams         []TeamSection     `json:"teams"`
	FooterBlocks  []ContentBlock    `json:"footer_blocks,omitempty"`
	Status        Status            `json:"status"`
	GeneratedAt   string            `json:"generated_at"`
	GeneratedBy   string            `json:"generated_by,omitempty"`
}

// TeamSection represents a section of the report for a specific check category.
type TeamSection struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Status        Status         `json:"status"`
	Verdict       string         `json:"verdict,omitempty"`
	Tasks         []TaskResult   `json:"tasks,omitempty"`
	ContentBlocks []ContentBlock `json:"content_blocks,omitempty"`
}

// TaskResult represents the result of a single check task.
type TaskResult struct {
	ID       string `json:"id"`
	Status   Status `json:"status"`
	Severity string `json:"severity,omitempty"`
	Detail   string `json:"detail,omitempty"`
}

// ContentBlock represents rich content in the report.
type ContentBlock struct {
	Type    string     `json:"type"`
	Title   string     `json:"title,omitempty"`
	Pairs   []KVPair   `json:"pairs,omitempty"`
	Items   []ListItem `json:"items,omitempty"`
	Content string     `json:"content,omitempty"`
}

// KVPair represents a key-value pair in content blocks.
type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Icon  string `json:"icon,omitempty"`
}

// ListItem represents an item in a list content block.
type ListItem struct {
	Text   string `json:"text"`
	Icon   string `json:"icon,omitempty"`
	Status Status `json:"status,omitempty"`
}

// GenerateReport creates a TeamReport from scan results.
func GenerateReport(results []*Result, project, version string) *TeamReport {
	report := &TeamReport{
		Schema:      "https://raw.githubusercontent.com/agentplexus/multi-agent-spec/main/schema/report/team-report.schema.json",
		Title:       "SVG SECURITY SCAN REPORT",
		Project:     project,
		Version:     version,
		Phase:       "SECURITY VALIDATION",
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		GeneratedBy: "brandkit security-scan",
		Teams:       []TeamSection{},
	}

	// Count totals
	totalFiles := len(results)
	secureFiles := 0
	threatsByType := make(map[ThreatType]int)
	var allThreats []Threat

	for _, r := range results {
		if r.IsSuccess() {
			secureFiles++
		}
		for _, t := range r.Threats {
			threatsByType[t.Type]++
			allThreats = append(allThreats, t)
		}
	}

	// Determine overall status
	report.Status = StatusGo
	if len(allThreats) > 0 {
		// Check if any critical/high threats
		hasCritical := threatsByType[ThreatScript] > 0 || threatsByType[ThreatEventHandler] > 0
		hasHigh := threatsByType[ThreatExternalRef] > 0 || threatsByType[ThreatXMLEntity] > 0

		if hasCritical || hasHigh {
			report.Status = StatusNoGo
		} else {
			report.Status = StatusWarn
		}
	}

	// Summary block
	report.SummaryBlocks = []ContentBlock{
		{
			Type: "kv_pairs",
			Pairs: []KVPair{
				{Key: "Files Scanned", Value: formatInt(totalFiles)},
				{Key: "Secure Files", Value: formatInt(secureFiles)},
				{Key: "Files with Threats", Value: formatInt(totalFiles - secureFiles)},
				{Key: "Total Threats", Value: formatInt(len(allThreats))},
			},
		},
	}

	// Create team sections for each threat category
	threatCategories := []struct {
		id         string
		name       string
		threatType ThreatType
		severity   string
	}{
		{"script-detection", "Script Detection", ThreatScript, "critical"},
		{"event-handler-detection", "Event Handler Detection", ThreatEventHandler, "critical"},
		{"external-ref-detection", "External Reference Detection", ThreatExternalRef, "high"},
		{"xml-entity-detection", "XML Entity Detection", ThreatXMLEntity, "high"},
		{"animation-detection", "Animation Detection", ThreatAnimation, "medium"},
		{"style-block-detection", "Style Block Detection", ThreatStyleBlock, "low"},
		{"link-detection", "Link Detection", ThreatLink, "medium"},
	}

	for _, cat := range threatCategories {
		count := threatsByType[cat.threatType]
		section := TeamSection{
			ID:   cat.id,
			Name: cat.name,
		}

		if count == 0 {
			section.Status = StatusGo
			section.Tasks = []TaskResult{
				{
					ID:     "scan",
					Status: StatusGo,
					Detail: "No threats detected",
				},
			}
		} else {
			// Determine status based on severity
			switch cat.severity {
			case "critical", "high":
				section.Status = StatusNoGo
			case "medium":
				section.Status = StatusWarn
			default:
				section.Status = StatusWarn
			}

			section.Tasks = []TaskResult{
				{
					ID:       "scan",
					Status:   section.Status,
					Severity: cat.severity,
					Detail:   formatInt(count) + " threat(s) detected",
				},
			}

			// Add content block with threat details
			var items []ListItem
			for _, r := range results {
				for _, t := range r.Threats {
					if t.Type == cat.threatType {
						icon := "🔴"
						if cat.severity == "medium" {
							icon = "🟡"
						} else if cat.severity == "low" {
							icon = "🟢"
						}
						items = append(items, ListItem{
							Icon: icon,
							Text: r.FilePath + ": " + t.Description,
						})
					}
				}
			}
			if len(items) > 0 {
				// Limit items to avoid huge reports
				if len(items) > 10 {
					items = items[:10]
					items = append(items, ListItem{
						Icon: "...",
						Text: "and more...",
					})
				}
				section.ContentBlocks = []ContentBlock{
					{
						Type:  "list",
						Title: "Findings",
						Items: items,
					},
				}
			}
		}

		report.Teams = append(report.Teams, section)
	}

	// Footer with action items if threats found
	if len(allThreats) > 0 {
		var actionItems []KVPair
		actionNum := 1

		if threatsByType[ThreatScript] > 0 || threatsByType[ThreatEventHandler] > 0 {
			actionItems = append(actionItems, KVPair{
				Icon:  "🔴",
				Key:   formatInt(actionNum),
				Value: "Remove all script elements and event handlers (CRITICAL)",
			})
			actionNum++
		}
		if threatsByType[ThreatExternalRef] > 0 {
			actionItems = append(actionItems, KVPair{
				Icon:  "🔴",
				Key:   formatInt(actionNum),
				Value: "Remove external references and foreignObject elements (HIGH)",
			})
			actionNum++
		}
		if threatsByType[ThreatXMLEntity] > 0 {
			actionItems = append(actionItems, KVPair{
				Icon:  "🟡",
				Key:   formatInt(actionNum),
				Value: "Remove DOCTYPE and ENTITY declarations (HIGH)",
			})
			actionNum++
		}
		if threatsByType[ThreatAnimation] > 0 {
			actionItems = append(actionItems, KVPair{
				Icon:  "🟡",
				Key:   formatInt(actionNum),
				Value: "Remove animation elements for static images (MEDIUM)",
			})
			actionNum++
		}
		if threatsByType[ThreatStyleBlock] > 0 {
			actionItems = append(actionItems, KVPair{
				Icon:  "🟢",
				Key:   formatInt(actionNum),
				Value: "Consider inlining styles and removing style blocks (LOW)",
			})
			actionNum++
		}
		if threatsByType[ThreatLink] > 0 {
			actionItems = append(actionItems, KVPair{
				Icon:  "🟡",
				Key:   formatInt(actionNum),
				Value: "Remove anchor elements for static images (MEDIUM)",
			})
		}

		if len(actionItems) > 0 {
			report.FooterBlocks = []ContentBlock{
				{
					Type:  "kv_pairs",
					Title: "ACTION ITEMS",
					Pairs: actionItems,
				},
			}
		}
	}

	return report
}

// ToJSON converts the report to JSON bytes.
func (r *TeamReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// formatInt converts an integer to string.
func formatInt(n int) string {
	return fmt.Sprintf("%d", n)
}
