package cli

import (
	"flag"
	"fmt"
	"os"
)

func displayHelpMenu() string {
	return `
  Welcome to ViteCPP!
	run "vitecpp help" for help
	run "vitecpp create -path=<path>" to create a new project.
	run "vitecpp add -package=<package>" to add a library.
	run "vitecpp build" to build the project.
	run "vitecpp run" to run the executable.
	run "vitecpp dev" to compile and run.
`
}

func Execute() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	var projectPath string
	createCmd.StringVar(&projectPath, "path", ".", "Path to new project")

	packageCmd := flag.NewFlagSet("add", flag.ExitOnError)
	var packageName string
	packageCmd.StringVar(&packageName, "package", "", "Package to add")

	if len(os.Args) < 2 {
		_, err := fmt.Fprintln(os.Stderr, "ERROR: Not enough Arguments run `vitecpp help` for help")
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}

	switch os.Args[1] {
	case "help", "-h", "--help":
		fmt.Println(displayHelpMenu())

	case "create":
		createCmd.Parse(os.Args[2:])
		fmt.Println("Creating Project At", projectPath)

	case "add":
		packageCmd.Parse(os.Args[2:])
		if packageName == "" {
			_, err := fmt.Fprintln(os.Stderr, "ERROR: Package Name Not Provided")
			if err != nil {
				panic(err)
			}
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
		_, err := fmt.Fprintln(os.Stderr, "ERROR : Invalid command run `vitecpp help` for help")
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}
