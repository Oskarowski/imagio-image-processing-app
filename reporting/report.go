package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func runCommand(command string, args []string) (string, error) {
	// Create the command
	cmd := exec.Command(command, args...)

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing command: %v, output: %s", err, output)
	}

	return string(output), nil
}

func main() {
	// base images
	// lenac := "../imgs/lenac.bmp"
	lenac_small := "../imgs/lenac_small.bmp"
	// lenag := "../imgs/lenag.bmp"

	// Images with impulse noise
	lenac_impulse_1_small := "../imgs/impulse_noise/lenac_impulse1_small.bmp"
	lenac_impulse_2_small := "../imgs/impulse_noise/lenac_impulse2_small.bmp"
	lenac_impulse_3_small := "../imgs/impulse_noise/lenac_impulse3_small.bmp"
	// lenag_impulse_1_small := "../imgs/impulse_noise/lenag_impulse1_small.bmp"
	// lenag_impulse_2_small := "../imgs/impulse_noise/lenag_impulse2_small.bmp"
	// lenag_impulse_3_small := "../imgs/impulse_noise/lenag_impulse3_small.bmp"

	// Images with normal distribution noise
	lenac_normal_1_small := "../imgs/noise_with_normal_distribution/lenac_normal1_small.bmp"
	lenac_normal_2_small := "../imgs/noise_with_normal_distribution/lenac_normal2_small.bmp"
	lenac_normal_3_small := "../imgs/noise_with_normal_distribution/lenac_normal3_small.bmp"
	// lenag_normal_1_small := "../imgs/noise_with_normal_distribution/lenag_normal1_small.bmp"
	// lenag_normal_2_small := "../imgs/noise_with_normal_distribution/lenag_normal2_small.bmp"
	// lenag_normal_3_small := "../imgs/noise_with_normal_distribution/lenag_normal3_small.bmp"

	// Images with uniform distribution noise
	lenac_uniform_1_small := "../imgs/noise_with_uniform_distribution/lenac_uniform1_small.bmp"
	lenac_uniform_2_small := "../imgs/noise_with_uniform_distribution/lenac_uniform2_small.bmp"
	lenac_uniform_3_small := "../imgs/noise_with_uniform_distribution/lenac_uniform3_small.bmp"
	// lenag_uniform_1_small := "../imgs/noise_with_uniform_distribution/lenag_uniform1_small.bmp"
	// lenag_uniform_2_small := "../imgs/noise_with_uniform_distribution/lenag_uniform2_small.bmp"
	// lenag_uniform_3_small := "../imgs/noise_with_uniform_distribution/lenag_uniform3_small.bmp"

	commands := [][]string{
		// {"--contrast", "-value=30", lenac_impulse_2_small},
		// {"--negative", lenac_impulse_3_small},

		{"--adaptive", lenac_small, lenac_impulse_1_small},
		{"--adaptive", lenac_small, lenac_impulse_2_small},
		{"--adaptive", lenac_small, lenac_impulse_3_small},

		{"--adaptive", lenac_small, lenac_normal_1_small},
		{"--adaptive", lenac_small, lenac_normal_2_small},
		{"--adaptive", lenac_small, lenac_normal_3_small},

		{"--adaptive", lenac_small, lenac_uniform_1_small},
		{"--adaptive", lenac_small, lenac_uniform_2_small},
		{"--adaptive", lenac_small, lenac_uniform_3_small},

		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_impulse_1_small, lenac_small},
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_impulse_2_small, lenac_small},
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_impulse_3_small, lenac_small},

		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_normal_1_small, lenac_small},
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_normal_2_small, lenac_small},
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_normal_3_small, lenac_small},

		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_uniform_1_small, lenac_small},
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_uniform_2_small, lenac_small},
		{"--mse", "--pmse", "--snr", "--psnr", "--md", lenac_uniform_3_small, lenac_small},
	}

	command := "go"

	currentTime := time.Now()
	timestamp := currentTime.Format("20060102150405") // Format: YYYYMMDDHHMMSS
	outputFileName := fmt.Sprintf("command_output_%s.txt", timestamp)

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile.Close()

	for _, args := range commands {
		fullArgs := append([]string{"run", "../main.go"}, args...)

		startTime := time.Now()
		output, err := runCommand(command, fullArgs)
		if err != nil {
			fmt.Println(err)
			return
		}
		duration := time.Since(startTime)

		_, err = fmt.Fprintf(outputFile, "Execution time for whole command '%s': %v\nOutput:\n%s\n\n", args, duration, output)
		if err != nil {
			fmt.Printf("Error writing to output file: %v\n", err)
			return
		}
	}

	fmt.Printf("Command execution results saved to %v\n", outputFileName)
}
