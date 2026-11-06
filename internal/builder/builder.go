package builder

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func checkIfProgramExists(program string) error {
	_, err := exec.LookPath(program)
	if err != nil {
		return fmt.Errorf("cannot find %v please install if not installed already.", program)
	}
	return nil
}

func checkIfFileExists(file string) error {
	_, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("cannot find %v\n%w", file, err)
	}
	return nil
}

func executeCmakeBuild() error {
	cmd := exec.Command("cmake", "-B", "build", "-S", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("errror running camke build %w", err)
	}

	cmd = exec.Command("cmake", "--build", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error running cmake %w", err)
	}
	return nil
}

func Build() error {
	// Check if CMake exists
	err := checkIfProgramExists("cmake")
	if err != nil {
		return err
	}
	// Check if CMakeLists.txt exist
	err = checkIfFileExists("CMakeLists.txt")
	if err != nil {
		return err
	}
	// Exec
	err = executeCmakeBuild()
	if err != nil {
		return err
	}
	return nil
}

func getProjectName() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get project name %w", err)
	}
	return filepath.Base(cwd), nil
}

func executeWindows(name string) error {
	projectName := name + ".exe"
	cmd := exec.Command(filepath.Join("build", projectName))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cannot run project %w", err)
	}
	return nil
}

func executeMacAndLinux(name string) error {
	cmd := exec.Command(filepath.Join("build", name))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cannot run project %w", err)
	}
	return nil
}

func handleExecutionBasedOnOs(os string, name string) error {
	switch os {
	case "windows":
		err := executeWindows(name)
		if err != nil {
			return err
		}
		return nil
	case "darwin", "linux":
		err := executeMacAndLinux(name)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("unsupported os")
	}
}

func Run(targetName ...string) error {
	var name string
	if len(targetName) > 0 && targetName[0] != "" {
		name = targetName[0]
	} else {
		var err error
		name, err = getProjectName()
		if err != nil {
			return err
		}
	}
	err := handleExecutionBasedOnOs(runtime.GOOS, name)
	if err != nil {
		return err
	}
	return nil
}
