package tui

import (
	"fmt"
	"image-processing/analysis"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TODO the header and tbody are off by one space
func renderComparisonResults(entries []analysis.ImageComparisonEntry) string {
	if len(entries) == 0 {
		return ""
	}

	metricWidth := 10
	descWidth := 40
	resultWidth := 20
	imgWidth := 15

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	cellStyle := lipgloss.NewStyle().Padding(0, 1)

	header := fmt.Sprintf("%-*s | %-*s | %-*s | %-*s | %-*s",
		metricWidth, "Metric",
		descWidth, "Description",
		resultWidth, "Result",
		imgWidth, "Img1",
		imgWidth, "Img2",
	)
	separator := strings.Repeat("-", len(header))

	var b strings.Builder
	b.WriteString(headerStyle.Render(header) + "\n")
	b.WriteString(headerStyle.Render(separator) + "\n")

	for _, entry := range entries {
		parts := strings.Split(entry.Result, ":")
		var resultDisplay string
		if len(parts) > 1 {
			resultDisplay = strings.TrimSpace(parts[1])
		} else {
			resultDisplay = entry.Result
		}

		line := fmt.Sprintf("%-*s | %-*s | %-*s | %-*s | %-*s",
			metricWidth, entry.MetricMethod,
			descWidth, entry.Description,
			resultWidth, resultDisplay,
			imgWidth, entry.Img1Name,
			imgWidth, entry.Img2Name,
		)
		b.WriteString(cellStyle.Render(line) + "\n")
	}

	return b.String()
}
