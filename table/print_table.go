/**
 * @Author: Einic <einicyeo AT gmail.com>
 * @Description:
 * @File: table
 * @Version: 1.0.0
 * @Date: 2024/2/8 15:25
 * @BLOG:  https://www.infvie.com
 * @Project home page:
 *     @https://github.com/Einic/EnvoyinStack
 */

package table

import (
	"fmt"
	"github.com/Einic/cops/lib"
	AlterResource "github.com/Einic/cops/resources"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func PrintUpdateTable(updateSlice []lib.ResourceInfo) {
	// Create a new table
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	// Customize the table header style
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateRows = true
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	var BoldStyle = table.Style{
		Name: "BoldStyle",
		Box:  table.StyleBoxBold,
		//Color:   table.ColorOptionsDefault,
		Color: table.ColorOptions{
			Header: text.Colors{text.Bold, text.Bold},
		},
		Format:  table.FormatOptionsDefault,
		HTML:    table.DefaultHTMLOptions,
		Options: table.OptionsDefault,
		Title:   table.TitleOptionsDefault,
	}
	t.SetStyle(BoldStyle)
	t.SetAutoIndex(true)

	// Append the header row with bold formatting
	headerRow := table.Row{"DataTime", "WORKLOAD", "CONTAINERNAME", "WORKTYPE", "NAMESPACE", "Replicas", "Requests (CPU)", "Requests (Memory)", "Limits (CPU)", "Limits (Memory)", "PodQos", "RUNSTATUS", "ALTERSTATUS"}
	// Set the color and style for the header row
	t.AppendHeader(headerRow, rowConfigAutoMerge)

	// Customize the table
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Replicas", Transformer: transformReplicas},
		{Name: "Requests (CPU)", Transformer: transformColorfulValue},
		{Name: "Requests (Memory)", Transformer: transformColorfulValue},
		{Name: "Limits (CPU)", Transformer: transformColorfulValue},
		{Name: "Limits (Memory)", Transformer: transformColorfulValue},
	})

	// Append rows for each update
	for _, update := range updateSlice {
		t.AppendSeparator()
		t.AppendRow([]interface{}{
			update.DataTime,
			update.Workload,
			update.ContainerName,
			update.WorkType,
			update.Namespace,
			fmt.Sprintf("%d -> %d", update.CurrentReplicas, update.AlterReplicas),
			fmt.Sprintf("%s -> %s", update.CurrentRequestsCPU, update.AlterRequestsCPU),
			fmt.Sprintf("%s -> %s", update.CurrentRequestsMemory, update.AlterRequestsMemory),
			fmt.Sprintf("%s -> %s", update.CurrentLimitsCPU, update.AlterLimitsCPU),
			fmt.Sprintf("%s -> %s", update.CurrentLimitsMemory, update.AlterLimitsMemory),
			update.PodQos,
			update.RunStatus,
			AlterResource.GetStatusText(update.AlterStatus),
		})
	}

	// Render the table
	t.Render()
}

// Custom transformer for Replicas column
func transformReplicas(data interface{}) string {
	if replicas, ok := data.(string); ok {
		parts := strings.Split(replicas, " -> ")
		if len(parts) == 2 {
			oldReplicas, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
			newReplicas, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			if newReplicas > oldReplicas {
				return text.FgGreen.Sprint(replicas)
			} else if newReplicas < oldReplicas {
				return text.FgRed.Sprint(replicas)
			}
		}
	}
	return fmt.Sprintf("%v", data)
}

// Custom transformer for colorizing values
func transformColorfulValue(data interface{}) string {
	switch value := data.(type) {
	case string:
		parts := strings.Split(value, " -> ")
		if len(parts) == 2 {
			oldValue, oldUnit := parseValueWithUnit(parts[0])
			newValue, newUnit := parseValueWithUnit(parts[1])

			// If units are different, just return the original value without coloring
			if oldUnit != newUnit {
				return value
			}

			if newValue > oldValue {
				return text.FgGreen.Sprint(value)
			} else if newValue < oldValue {
				return text.FgRed.Sprint(value)
			}
		}
	}
	return fmt.Sprintf("%v", data)
}

// Helper function to parse value with unit
func parseValueWithUnit(valueWithUnit string) (float64, error) {
	// Trim leading and trailing spaces
	valueWithUnit = strings.TrimSpace(valueWithUnit)

	// Find the first non-digit index
	nonDigitIndex := strings.IndexFunc(valueWithUnit, func(r rune) bool {
		return !unicode.IsDigit(r)
	})

	// If no non-digit character found, assume no unit
	if nonDigitIndex == -1 {
		value, err := strconv.ParseFloat(valueWithUnit, 64)
		if err != nil {
			return 0, fmt.Errorf("error parsing value: %s", err)
		}
		return value, nil
	}

	// Parse value and unit
	valueStr := valueWithUnit[:nonDigitIndex]
	unit := valueWithUnit[nonDigitIndex:]

	// Parse value
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing value: %s", err)
	}

	// Handle unit conversion
	switch unit {
	case "m":
		// Convert millicores to cores
		value /= 1000
	case "Mi":
		// Convert Mebibytes to Megabytes
		value *= 1024
	// Add more cases for other units as needed
	default:
		return 0, fmt.Errorf("unknown unit: %s", unit)
	}

	return value, nil
}
