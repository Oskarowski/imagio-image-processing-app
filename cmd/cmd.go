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

type CommandInfo struct {
	Name        string
	Usage       string
	Description string
	Arguments   []string
}

var AvailableCommands = []CommandInfo{
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
	{"histogram", "--histogram <bmp_image_path>", "Generate and save a graphical representation of the histogram of the image.", []string{}},
	{"hrayleigh", "--hrayleigh -min=0 -max=255 -alpha=\"0.2\" <bmp_image_path>", "Apply Rayleigh transformation to the image.", []string{"-min=(int): Minimum value in the range [0, 255].", "-max=(int): Maximum value in the range [0, 255], must be greater than min.", "-alpha=(float): Alpha value for transformation. Note: Quote float values (e.g., -alpha=\"0.5\")."}},
	{"cmean", "--cmean <bmp_image_path>", "Calculate the mean intensity from the histogram of the image.", []string{}},
	{"cvariance", "--cvariance <bmp_image_path>", "Calculate the variance intensity from the histogram of the image.", []string{}},
	{"cstdev", "--cstdev <bmp_image_path>", "Calculate the standard deviation from the histogram of the image.", []string{}},
	{"cvarcoi", "--cvarcoi <bmp_image_path>", "Calculate the coefficient of variation (type I) from the histogram.", []string{}},
	{"casyco", "--casyco <bmp_image_path>", "Calculate the asymmetry coefficient from the histogram.", []string{}},
	{"cflatco", "--cflatco <bmp_image_path>", "Calculate the flattening coefficient from the histogram.", []string{}},
	{"cvarcoii", "--cvarcoii <bmp_image_path>", "Calculate the coefficient of variation (type II) from the histogram.", []string{}},
	{"centropy", "--centropy <bmp_image_path>", "Calculate the entropy from the histogram of the image.", []string{}},
	{"sedgesharp", "--sedgesharp -mask=\"edge1\" <bmp_image_path>", "Apply edge sharpening with the specified mask.", []string{"-mask=(string): The name of the mask to use."}},
	{"okirsf", "--okirsf <bmp_image_path>", "Apply Kirsch edge detection to the image.", []string{}},
	{"dilation", "--dilation -se=<structuring_element> <bmp_image_path>", "Apply dilation operation using the specified structuring element.", []string{"-se=(string): Name of SE based on structure_elements.json."}},
	{"erosion", "--erosion -se=<structuring_element> <bmp_image_path>", "Apply erosion operation using the specified structuring element.", []string{"-se=(string): Name of SE based on structure_elements.json."}},
	{"opening", "--opening -se=<structuring_element> <bmp_image_path>", "Apply opening operation using the specified structuring element.", []string{"-se=(string): Name of SE based on structure_elements.json."}},
	{"closing", "--closing -se=<structuring_element> <bmp_image_path>", "Apply closing operation using the specified structuring element.", []string{"-se=(string): Name of SE based on structure_elements.json."}},
	{"hmt", "--hmt -se1=<foreground_se> -se2=<background_se> <bmp_image_path>", "Perform hit-or-miss transformation using foreground and background structuring elements.", []string{
		"-se1=(string): Path to or inline definition of the foreground structuring element.",
		"-se2=(string): Path to or inline definition of the background structuring element.",
	}},
	{"thinning", "--thinning <bmp_image_path>", "Apply thinning operation to the image.", []string{}},
	{"regiongrow", "--regiongrow -seeds=<seeds> -metric=<metric> -threshold=<value> <bmp_image_path>", "Perform region growing segmentation on the image.", []string{
		"-seeds=(string): List of seed points as [x,y][x,y][x,y].",
		"-metric=(int): Distance metric ('0 - Euclidean', '1 - Manhattan', '2 - Chebyshev').",
		"-threshold=(double): Similarity threshold for region growing.",
	}},
	{"help", "--help", "Show this help message.", []string{}},
}

func PrintHelp() {
	fmt.Println("Usage: go run main.go <command> [-argument=value [...]] <bmp_image_path> [<second_image_path>]")
	fmt.Println("\nAvailable commands:")

	for _, cmd := range AvailableCommands {
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

func (commands Commands) Includes(name string) bool {
	for _, cmd := range commands {
		if cmd.Name == name {
			return true
		}
	}
	return false
}
