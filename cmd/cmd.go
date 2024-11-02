package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func IsImagePath(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".bmp"
}

type Command struct {
	Name string
	Args map[string]string
}
type Commands []Command

func ParseCommands(args []string) Commands {
	var commands []Command
	var currentCommand *Command

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			if currentCommand != nil {
				commands = append(commands, *currentCommand)
			}

			currentCommand = &Command{
				Name: strings.TrimPrefix(arg, "--"),
				Args: make(map[string]string),
			}
		} else if strings.HasPrefix(arg, "-") && currentCommand != nil {
			// -argument=value), split it by '='

			parts := strings.SplitN(arg, "=", 2)

			if len(parts) == 2 {
				key := strings.TrimPrefix(parts[0], "-")
				currentCommand.Args[key] = parts[1]
			} else {
				currentCommand.Args[parts[0]] = ""
			}
		}
	}

	if currentCommand != nil {
		commands = append(commands, *currentCommand)
	}

	return commands
}

type commandInfo struct {
	Name        string
	Usage       string
	Description string
	Arguments   []string
}

var availableCommands = []commandInfo{
	{"brightness", "--brightness -value=50 <bmp_image_path>", "Adjust brightness of the image.", []string{"-value=(int): Brightness percentage adjustment value."}},
	{"contrast", "--contrast -value=25 <bmp_image_path>", "Adjust contrast of the image.", []string{"-value=(int): Contrast adjustment value (-255 to 255)."}},
	{"negative", "--negative <bmp_image_path>", "Create a negative of the image.", []string{}},
	{"hflip", "--hflip <bmp_image_path>", "Flip the image horizontally.", []string{}},
	{"vflip", "--vflip <bmp_image_path>", "Flip the image vertically.", []string{}},
	{"dflip", "--dflip <bmp_image_path>", "Flip the image diagonally.", []string{}},
	{"shrink", "--shrink -value=2 <bmp_image_path>", "Shrink the image by a factor.", []string{"-value=(int): Shrink factor."}},
	{"enlarge", "--enlarge -value=2 <bmp_image_path>", "Enlarge the image by a factor.", []string{"-value=(int): Enlarge factor."}},
	{"adaptive", "--adaptive <bmp_image_path>", "Apply adaptive median noise removal filter to the image.", []string{"-min=(int): Minimal size of window size for filter defaults, to 3.", "-max=(int): Maximal size of window size for filter defaults, to 7."}},
	{"min", "--min -value=3 <bmp_image_path>", "Apply min noise removal filter.", []string{"-value=(int): Window size."}},
	{"max", "--max -value=3 <bmp_image_path>", "Apply max noise removal filter.", []string{"-value=(int): Window size."}},
	{"mse", "--mse <comparison_image_path> <bmp_image_path>", "Calculate Mean Square Error with a comparison image.", []string{}},
	{"pmse", "--pmse <comparison_image_path> <bmp_image_path>", "Calculate Peak Mean Square Error with a comparison image.", []string{}},
	{"snr", "--snr <comparison_image_path> <bmp_image_path>", "Calculate Signal to Noise Ratio with a comparison image.", []string{}},
	{"psnr", "--psnr <comparison_image_path> <bmp_image_path>", "Calculate Peak Signal to Noise Ratio with a comparison image.", []string{}},
	{"md", "--md <comparison_image_path> <bmp_image_path>", "Calculate Max Difference with a comparison image.", []string{}},
	{"help", "--help", "Show this help message.", []string{}},
}

func PrintHelp() {
	fmt.Println("Usage: go run main.go <command> [-argument=value [...]] <bmp_image_path> [<second_image_path>]")
	fmt.Println("\nAvailable commands:")

	for _, cmd := range availableCommands {
		fmt.Printf(" %s\n", cmd.Usage)
		fmt.Printf("   Description: %s\n", cmd.Description)

		if len(cmd.Arguments) > 0 {
			fmt.Println("   Arguments:")
			for _, arg := range cmd.Arguments {
				fmt.Printf("    %s\n", arg)
			}
		}
		fmt.Println()
	}
}

func GetOrDefault[T int | string | float64](val string, defaultValue T) T {
	if val == "" {
		return defaultValue
	}

	var result any
	switch any(defaultValue).(type) {
	case int:
		if num, err := strconv.Atoi(val); err == nil {
			result = num
		} else {
			result = defaultValue
		}
	case float64:
		if num, err := strconv.ParseFloat(val, 64); err == nil {
			result = num
		} else {
			result = defaultValue
		}
	case string:
		result = val

	default:
		result = defaultValue
	}

	return result.(T)
}
