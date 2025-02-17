package tui

import (
	"image-processing/analysis"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderComparisonResults(entries []analysis.CharacteristicsEntry) string {
	if len(entries) == 0 {
		return ""
	}

	metricWidth := 10
	descWidth := 70
	resultWidth := 18
	imgWidth := 45

	metricHeader := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Width(metricWidth).Render("Metric")
	resultHeader := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Width(resultWidth).Render("Result")
	img1Header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Width(imgWidth).Render("Img1")
	img2Header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Width(imgWidth).Render("Img2")
	descHeader := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Width(descWidth).Render("Description")

	sep := " | "
	header := strings.Join([]string{
		metricHeader, resultHeader, img1Header, img2Header, descHeader,
	}, sep)

	separator := strings.Repeat("-", lipgloss.Width(header))

	var rows []string
	for _, entry := range entries {
		parts := strings.Split(entry.Result, ":")
		var resultDisplay string
		if len(parts) > 1 {
			resultDisplay = strings.TrimSpace(parts[1])
		} else {
			resultDisplay = entry.Result
		}

		img2 := entry.Img2Name
		if strings.TrimSpace(img2) == "" {
			img2 = "N/A"
		}

		metricCell := lipgloss.NewStyle().Width(metricWidth).MaxWidth(metricWidth).Inline(true).Render(entry.MetricMethod)
		resultCell := lipgloss.NewStyle().Width(resultWidth).MaxWidth(resultWidth).Inline(true).Render(resultDisplay)
		img1Cell := lipgloss.NewStyle().Width(imgWidth).MaxWidth(imgWidth).Inline(true).Render(entry.Img1Name)
		img2Cell := lipgloss.NewStyle().Width(imgWidth).MaxWidth(imgWidth).Inline(true).Render(img2)
		descCell := lipgloss.NewStyle().Width(descWidth).MaxWidth(descWidth).Inline(true).Render(entry.Description)

		row := strings.Join([]string{
			metricCell, resultCell, img1Cell, img2Cell, descCell,
		}, sep)
		rows = append(rows, row)
	}

	table := strings.Join(append([]string{header, separator}, rows...), "\n")
	table = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(table)

	return table
}
