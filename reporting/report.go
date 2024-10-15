package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type OperationDetail struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Duration    string `json:"duration"`
	Result      string `json:"result,omitempty"`
}

type Result struct {
	Command            []string          `json:"command"`
	TotalExecutionTime string            `json:"total_execution_time"`
	TotalOperationTime string            `json:"total_operation_time"`
	Operations         []OperationDetail `json:"operations"`
}

func parseOutput(output string) (Result, error) {
	var operations []OperationDetail
	var totalOperationTime string

	lines := strings.Split(output, "\n")

	var currentOp OperationDetail

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Command:") {

			if currentOp.Name != "" {
				operations = append(operations, currentOp)
			}

			currentOp = OperationDetail{}
			currentOp.Name = strings.TrimSpace(strings.Split(line, ":")[1])

		} else if strings.HasPrefix(line, "Description:") {
			currentOp.Description = strings.TrimSpace(strings.Split(line, ":")[1])

		} else if strings.HasPrefix(line, "Result:") {
			currentOp.Result = strings.TrimSpace(strings.Split(line, ":")[2])

		} else if strings.HasPrefix(line, "Duration:") {
			currentOp.Duration = strings.TrimSpace(strings.Split(line, ":")[1])

		} else if strings.HasPrefix(line, "Total operation time:") {
			totalOperationTime = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	if currentOp.Name != "" {
		operations = append(operations, currentOp)
	}

	return Result{
		Operations:         operations,
		TotalOperationTime: totalOperationTime,
	}, nil
}

func runCommand(command string, args []string) (Result, error) {
	cmd := exec.Command(command, args...)

	// Capture output
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return Result{}, fmt.Errorf("error executing command: %v, output: %s", err, outputBytes)
	}

	output := string(outputBytes)

	return parseOutput(output)
}

func main() {
	// base images
	lenac := "../imgs/lenac.bmp"
	lenag := "../imgs/lenag.bmp"

	// Images with impulse noise
	lenac_impulse_1 := "../imgs/impulse_noise/lenac_impulse1.bmp"
	lenac_impulse_2 := "../imgs/impulse_noise/lenac_impulse2.bmp"
	lenac_impulse_3 := "../imgs/impulse_noise/lenac_impulse3.bmp"
	// lenag_impulse_1 := "../imgs/impulse_noise/lena_impulse1.bmp"
	// lenag_impulse_2 := "../imgs/impulse_noise/lena_impulse2.bmp"
	lenag_impulse_3 := "../imgs/impulse_noise/lena_impulse3.bmp"

	// Images with normal distribution noise
	lenac_normal_1 := "../imgs/noise_with_normal_distribution/lenac_normal1.bmp"
	lenac_normal_2 := "../imgs/noise_with_normal_distribution/lenac_normal2.bmp"
	lenac_normal_3 := "../imgs/noise_with_normal_distribution/lenac_normal3.bmp"
	// lenag_normal_1 := "../imgs/noise_with_normal_distribution/lena_normal1.bmp"
	// lenag_normal_2 := "../imgs/noise_with_normal_distribution/lena_normal2.bmp"
	lenag_normal_3 := "../imgs/noise_with_normal_distribution/lena_normal3.bmp"

	// Images with uniform distribution noise
	lenac_uniform_1 := "../imgs/noise_with_uniform_distribution/lenac_uniform1.bmp"
	lenac_uniform_2 := "../imgs/noise_with_uniform_distribution/lenac_uniform2.bmp"
	lenac_uniform_3 := "../imgs/noise_with_uniform_distribution/lenac_uniform3.bmp"
	// lenag_uniform_1 := "../imgs/noise_with_uniform_distribution/lena_uniform1.bmp"
	// lenag_uniform_2 := "../imgs/noise_with_uniform_distribution/lena_uniform2.bmp"
	lenag_uniform_3 := "../imgs/noise_with_uniform_distribution/lena_uniform3.bmp"

	commands := [][]string{

		// Adaptive median filters
		{"--adaptive", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_impulse_1},
		{"--adaptive", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_impulse_2},
		{"--adaptive", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_impulse_3},

		{"--adaptive", "--mse", "--pmse", "--snr", "--psnr", "--md", lenag, lenag_impulse_3},

		{"--adaptive", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_normal_1},
		{"--adaptive", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_normal_2},
		{"--adaptive", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_normal_3},

		{"--adaptive", "--mse", "--pmse", "--snr", "--psnr", "--md", lenag, lenag_normal_3},

		{"--adaptive", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_uniform_1},
		{"--adaptive", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_uniform_2},
		{"--adaptive", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_uniform_3},

		{"--adaptive", "--mse", "--pmse", "--snr", "--psnr", "--md", lenag, lenag_uniform_3},

		// Max filters
		{"--max", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_impulse_1},
		{"--max", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_impulse_2},
		{"--max", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_impulse_3},

		{"--max", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenag, lenag_impulse_3},

		{"--max", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_normal_1},
		{"--max", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_normal_2},
		{"--max", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_normal_3},

		{"--max", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenag, lenag_normal_3},

		{"--max", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_uniform_1},
		{"--max", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_uniform_2},
		{"--max", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_uniform_3},

		{"--max", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenag, lenag_uniform_3},

		// Min filters
		{"--min", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_impulse_1},
		{"--min", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_impulse_2},
		{"--min", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_impulse_3},

		{"--min", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenag, lenag_impulse_3},

		{"--min", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_normal_1},
		{"--min", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_normal_2},
		{"--min", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_normal_3},

		{"--min", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenag, lenag_normal_3},

		{"--min", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_uniform_1},
		{"--min", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_uniform_2},
		{"--min", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenac, lenac_uniform_3},

		{"--min", "-value=5", "--mse", "--pmse", "--snr", "--psnr", "--md", lenag, lenag_uniform_3},

		// General analysis for provided images
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_impulse_1, lenac},
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_impulse_2, lenac},
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_impulse_3, lenac},

		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenag_impulse_3, lenag},

		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_normal_1, lenac},
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_normal_2, lenac},
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_normal_3, lenac},

		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenag_normal_3, lenag},

		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_uniform_1, lenac},
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_uniform_2, lenac},
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_uniform_3, lenac},

		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenag_uniform_3, lenag},
	}

	command := "go"

	startReportTime := time.Now()
	timestamp := startReportTime.Format("20060102150405") // Format: YYYYMMDDHHMMSS

	var results []Result

	for _, args := range commands {
		fullArgs := append([]string{"run", "../main.go"}, args...)

		startTime := time.Now()
		result, err := runCommand(command, fullArgs)
		if err != nil {
			fmt.Println(err)
			return
		}
		duration := time.Since(startTime)

		result.Command = args
		result.TotalExecutionTime = duration.String()

		results = append(results, result)
	}

	resultsJSON, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		log.Fatalf("Error marshalling results to JSON: %v", err)
	}

	outputFileName := fmt.Sprintf("results_%s.json", timestamp)
	err = os.WriteFile(outputFileName, resultsJSON, 0644)
	if err != nil {
		panic(err)
	}

	totalDuration := time.Since(startReportTime)

	fmt.Printf("Command execution results saved to %v\n", outputFileName)
	fmt.Printf("Total execution time for all commands: %v\n", totalDuration)
}
