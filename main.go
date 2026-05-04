package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Definisi Warna
const (
	RED    = "\033[0;31m"
	GREEN  = "\033[0;32m"
	CYAN   = "\033[0;36m"
	YELLOW = "\033[1;33m"
	BLUE   = "\033[0;34m"
	PURPLE = "\033[0;35m"
	NC     = "\033[0m"
)

type Target struct {
	os   string
	arch string
	ext  string
}

func showBanner() {
	fmt.Printf("%s", PURPLE)
	fmt.Println("          __            __")
	fmt.Println("         /  \\          /  \\")
	fmt.Println("        |    \\        /    |")
	fmt.Println("        |     \\      /     |")
	fmt.Println("         \\     \\    /     /")
	fmt.Println("          \\     \\  /     /")
	fmt.Println("           \\     \\/     /")
	fmt.Println("            \\    /\\    /")
	fmt.Println("             \\  /  \\  /")
	fmt.Println("              \\/    \\/")
	fmt.Println("      [ PACKAGESOFT - GO BUILDER ]")
	fmt.Println("      [ Author: Muhammad Quwais Saputra ]")
	fmt.Printf("%s\n", NC)
}

// Bypass Termux kernel sandbox crash (SIGSYS: bad system call)
// Kita menggunakan os.Stat manual alih-alih exec.LookPath bawaan Go
func findGoPath() string {
	// Prioritas 1: Path absolut default Termux
	termuxPath := "/data/data/com.termux/files/usr/bin/go"
	if _, err := os.Stat(termuxPath); err == nil {
		return termuxPath
	}

	// Prioritas 2: Loop manual melalui environment $PATH untuk OS lain (Linux/Mac/Win)
	pathEnv := os.Getenv("PATH")
	for _, p := range strings.Split(pathEnv, ":") {
		fullPath := filepath.Join(p, "go")
		if stat, err := os.Stat(fullPath); err == nil && !stat.IsDir() {
			return fullPath
		}
	}
	
	// Fallback ke string "go" biasa jika semuanya gagal
	return "go"
}

func main() {
	if len(os.Args) < 2 {
		showBanner()
		fmt.Printf("%sUsage: go run packagesoft.go <filename.go>%s\n", YELLOW, NC)
		return
	}

	sourceFile := os.Args[1]
	binName := strings.TrimSuffix(filepath.Base(sourceFile), ".go")
	outDir := "packagesoft_build"

	// Buat direktori output
	os.MkdirAll(outDir, 0755)

	showBanner()
	fmt.Printf("%s[*] Compiling Go Source: %s%s%s\n", BLUE, YELLOW, sourceFile, NC)
	fmt.Println("------------------------------------------------")

	// Daftar platform target
	targets := []Target{
		{"linux", "amd64", ""},
		{"linux", "386", ""},
		{"linux", "arm64", ""},
		{"linux", "arm", ""},
		{"windows", "amd64", ".exe"},
		{"windows", "386", ".exe"},
		{"darwin", "amd64", ""}, // macOS Intel
		{"darwin", "arm64", ""}, // macOS M1/M2/M3
	}

	// Ambil path aman untuk binary go
	goBin := findGoPath()

	for _, t := range targets {
		suffix := fmt.Sprintf("%s_%s", t.os, t.arch)
		outFile := filepath.Join(outDir, fmt.Sprintf("%s_%s%s", binName, suffix, t.ext))

		fmt.Printf("%s[>] Building for %-15s... %s", CYAN, suffix, NC)

		// Eksekusi build dengan command bypass
		cmd := exec.Command(goBin, "build", "-ldflags", "-s -w", "-o", outFile, sourceFile)
		cmd.Env = append(os.Environ(),
			"GOOS="+t.os,
			"GOARCH="+t.arch,
			"CGO_ENABLED=0", // Paksa mode statis sepenuhnya
		)

		err := cmd.Run()

		if err == nil {
			fmt.Printf("%sSUCCESS%s\n", GREEN, NC)
		} else {
			fmt.Printf("%sFAILED%s\n", RED, NC)
		}
	}

	fmt.Println("------------------------------------------------")
	fmt.Printf("%s[+] Build finished! Binaries are in '%s/'%s\n", GREEN, outDir, NC)
}
