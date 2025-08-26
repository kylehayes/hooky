package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "1.3.0"

func main() {
	var (
		configFile = flag.String("config", "hooky.yaml", "Path to configuration file")
		install    = flag.Bool("install", false, "Install hooks")
		uninstall  = flag.Bool("uninstall", false, "Uninstall hooks")
		list       = flag.Bool("list", false, "List available hooks")
		verbose    = flag.Bool("verbose", false, "Enable verbose output")
		showVersion = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("hooky version %s\n", version)
		return
	}

	manager := NewHookManager(*configFile, *verbose)

	switch {
	case *install:
		if err := manager.InstallHooks(); err != nil {
			fmt.Fprintf(os.Stderr, "Error installing hooks: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Hooks installed successfully")

	case *uninstall:
		if err := manager.UninstallHooks(); err != nil {
			fmt.Fprintf(os.Stderr, "Error uninstalling hooks: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Hooks uninstalled successfully")

	case *list:
		if err := manager.ListHooks(); err != nil {
			fmt.Fprintf(os.Stderr, "Error listing hooks: %v\n", err)
			os.Exit(1)
		}

	default:
		flag.Usage()
	}
}