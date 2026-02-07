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
		w := os.Stderr
		fmt.Fprintln(w, "Usage: difi [flags] [target-branch]")
		fmt.Fprintln(w, "\nFlags:")
		flag.PrintDefaults()
		fmt.Fprintln(w, "\nExamples:")
		fmt.Fprintln(w, "  difi             # Diff against default")
		fmt.Fprintln(w, "  difi develop     # Diff against develop")
		fmt.Fprintln(w, "  difi HEAD~1      # Diff against last commit")
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
		// Uses --name-status for a concise, machine-readable summary suitable for CI
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
