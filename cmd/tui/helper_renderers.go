package tui

import (
	"imagio/analysis"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderComparisonResults(entries []analysis.CharacteristicsEntry) string {
	if len(entries) == 0 {
		return ""
	}

	termWidth := m.terminalSize.width

	borderPadding := 8
	tableWidth := termWidth - borderPadding
	sep := " | "

	metricWidth := tableWidth * 10 / 100
	resultWidth := tableWidth * 15 / 100
	imgWidth := tableWidth * 20 / 100
	descWidth := tableWidth - (metricWidth + resultWidth + imgWidth*2 + len(sep)*4)

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	metricHeader := headerStyle.Width(metricWidth).Render("Metric")
	resultHeader := headerStyle.Width(resultWidth).Render("Result")
	img1Header := headerStyle.Width(imgWidth).Render("Img1")
	img2Header := headerStyle.Width(imgWidth).Render("Img2")
	descHeader := headerStyle.Width(descWidth).Render("Description")

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

		metricCell := lipgloss.NewStyle().Width(metricWidth).MaxWidth(metricWidth).Render(entry.MetricMethod)
		resultCell := lipgloss.NewStyle().Width(resultWidth).MaxWidth(resultWidth).Render(resultDisplay)
		img1Cell := lipgloss.NewStyle().Width(imgWidth).MaxWidth(imgWidth).Render(entry.Img1Name)
		img2Cell := lipgloss.NewStyle().Width(imgWidth).MaxWidth(imgWidth).Render(img2)
		descCell := lipgloss.NewStyle().Width(descWidth).MaxWidth(descWidth).Render(entry.Description)

		row := strings.Join([]string{
			metricCell, resultCell, img1Cell, img2Cell, descCell,
		}, sep)
		rows = append(rows, row)
	}

	table := strings.Join(append([]string{header, separator}, rows...), "\n")
	table = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(table)

	return table
}
