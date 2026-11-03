package builder

import (
	"fmt"
	"os"
	"os/exec"
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
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("errror running camke build %w", err)
	}

	cmd = exec.Command("cmake", "--build", "build")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
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
