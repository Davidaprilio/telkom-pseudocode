package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"telucode/internal/parser"
)

func CompileAndRun(filePath string) error {
	tmpPath, err := transpileToTempGo(filePath)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPath)

	cmd := exec.Command("go", "run", tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return fmt.Errorf("perintah 'go' tidak ditemukan di PATH")
		}
		return fmt.Errorf("eksekusi go run gagal: %w", err)
	}

	return nil
}

func CompileAndBuild(filePath, outputPath string) error {
	tmpPath, err := transpileToTempGo(filePath)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPath)

	if strings.TrimSpace(outputPath) == "" {
		outputPath = defaultBinaryName(filePath)
	}

	cmd := exec.Command("go", "build", "-o", outputPath, tmpPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return fmt.Errorf("perintah 'go' tidak ditemukan di PATH")
		}
		return fmt.Errorf("eksekusi go build gagal: %w", err)
	}

	fmt.Printf("binary berhasil dibuat: %s\n", outputPath)
	return nil
}

func CompileToGo(filePath, outputPath string) error {
	if filepath.Ext(filePath) != ".telu" {
		return fmt.Errorf("file input harus berekstensi .telu")
	}

	sourceBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("gagal membaca file '%s': %w", filePath, err)
	}

	goCode, err := parser.Transpile(string(sourceBytes))
	if err != nil {
		return fmt.Errorf("gagal transpile: %w", err)
	}

	if strings.TrimSpace(outputPath) == "" {
		outputPath, err = defaultGoOutputPath(filePath)
		if err != nil {
			return err
		}
	}

	if err := os.WriteFile(outputPath, []byte(goCode), 0644); err != nil {
		return fmt.Errorf("gagal menulis output Go '%s': %w", outputPath, err)
	}

	fmt.Printf("file Go berhasil dibuat: %s\n", outputPath)
	return nil
}

func transpileToTempGo(filePath string) (string, error) {
	if filepath.Ext(filePath) != ".telu" {
		return "", fmt.Errorf("file input harus berekstensi .telu")
	}

	sourceBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("gagal membaca file '%s': %w", filePath, err)
	}

	goCode, err := parser.Transpile(string(sourceBytes))
	if err != nil {
		return "", fmt.Errorf("gagal transpile: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "telupsc-*.go")
	if err != nil {
		return "", fmt.Errorf("gagal membuat file temporer: %w", err)
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.WriteString(goCode); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("gagal menulis file temporer: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("gagal menutup file temporer: %w", err)
	}

	return tmpPath, nil
}

func defaultBinaryName(filePath string) string {
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(base), ".exe") {
		return base + ".exe"
	}
	return base
}

func defaultGoOutputPath(filePath string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("gagal membaca current working directory: %w", err)
	}
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath)) + ".go"
	return filepath.Join(wd, base), nil
}
