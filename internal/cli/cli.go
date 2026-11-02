package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/AmlanJSarmah/ViteCPP/internal/scaffold"
)

func displayHelpMenu() string {
	return `
  Welcome to ViteCPP!
	run "vitecpp help" for help
	run "vitecpp create -name=<name>" to create a new project.
	run "vitecpp add -package=<package>" to add a library.
	run "vitecpp build" to build the project.
	run "vitecpp run" to run the executable.
	run "vitecpp dev" to compile and run.
`
}

func Execute() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	var projectName string
	createCmd.StringVar(&projectName, "name", ".", "new project name")

	packageCmd := flag.NewFlagSet("add", flag.ExitOnError)
	var packageName string
	packageCmd.StringVar(&packageName, "package", "", "package to add")

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "ERROR: not enough arguments run `vitecpp help` for help")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "help", "-h", "--help":
		fmt.Println(displayHelpMenu())

	case "create":
		createCmd.Parse(os.Args[2:])
		err := scaffold.CreateProject(projectName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR creating project :", err)
			os.Exit(1)
		}

	case "add":
		packageCmd.Parse(os.Args[2:])
		if packageName == "" {
			fmt.Fprintln(os.Stderr, "ERROR: package name not provided")
			os.Exit(1)
		}
		fmt.Println("Installing", packageName)

	case "build":
		fmt.Println("Building...")

	case "run":
		fmt.Println("Running...")

	case "dev":
		fmt.Println("Building & Running...")

	default:
		fmt.Fprintln(os.Stderr, "ERROR : invalid command run `vitecpp help` for help")
		os.Exit(1)
	}
}
