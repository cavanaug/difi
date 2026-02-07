package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oug-t/difi/internal/config"
	"github.com/oug-t/difi/internal/ui"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	plain := flag.Bool("plain", false, "Print a plain, non-interactive summary and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: difi [flags] [target-branch]\n")
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  difi              # Diff against main\n")
		fmt.Fprintf(os.Stderr, "  difi develop      # Diff against develop\n")
		fmt.Fprintf(os.Stderr, "  difi HEAD~1       # Diff against last commit\n")
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("difi version %s\n", version)
		os.Exit(0)
	}

	target := "HEAD"
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}

	if *plain {
		cmd := exec.Command("git", "diff", "--name-status", fmt.Sprintf("%s...HEAD", target))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	cfg := config.Load()

	p := tea.NewProgram(ui.NewModel(cfg, target), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
