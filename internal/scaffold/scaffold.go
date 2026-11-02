package scaffold

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var GITIGNORE = `build/
compiler_commands.json
.cache/
`

var CODE = `#include <iostream>

int main() {
	std::cout << "Hello from ViteCpp" << std::endl;
	return 0;
}
`

var CMAKE = `cmake_minimum_required(VERSION 3.20)

# Project name (will match the directory or app name)
project({{PROJECT_NAME}} LANGUAGES CXX)

# Require C++20
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Generates compile_commands.json for clangd/LSP support in VS Code, Neovim, etc.
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Define our executable
add_executable({{PROJECT_NAME}}
    src/main.cpp
)

# Allow #include "..." from the include directory
target_include_directories({{PROJECT_NAME}} PRIVATE include)
`

// Check if CMake Project already exists in current directory.
func checkIfProjectAlreadyExist() (bool, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false, errors.New("Failed to Read Contents of current directorys")
	}

	for _, entry := range entries {
		if entry.Name() == "CMakeLists.txt" {
			return true, nil
		}
	}

	return false, nil
}

func constructProjectPaths(name string) (string, string, string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", "", "", errors.New("Failed to get current directory name")
	}
	var rootPath string
	var srcPath string
	var includePath string
	if name == "." {
		srcPath = filepath.Join(currentDir, "src")
		includePath = filepath.Join(currentDir, "include")
		rootPath = currentDir
		return rootPath, srcPath, includePath, nil
	} else {
		srcPath = filepath.Join(currentDir, name, "src")
		includePath = filepath.Join(currentDir, name, "include")
		rootPath = filepath.Join(currentDir, name)
		return rootPath, srcPath, includePath, nil
	}
}

func createDirectory(srcPath string, includePath string) error {
	err := os.MkdirAll(srcPath, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Failed to create new directory ")
	}

	err = os.MkdirAll(includePath, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Failed to create new directory")
	}
	return nil
}

func createTemplateFiles(rootPath string, srcPath string) error {
	cMakeFilePath := filepath.Join(rootPath, "CMakeLists.txt")
	projectName := filepath.Base(rootPath)
	CMAKE = strings.ReplaceAll(CMAKE, "{{PROJECT_NAME}}", projectName)
	cMakeContent := []byte(CMAKE)
	err := os.WriteFile(cMakeFilePath, cMakeContent, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Failed to create CMakeLists.txt")
	}

	gitIgnoreFilePath := filepath.Join(rootPath, ".gitignore")
	gitIgnoreContent := []byte(GITIGNORE)
	err = os.WriteFile(gitIgnoreFilePath, gitIgnoreContent, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Failed to create .gitignore")
	}

	codePath := filepath.Join(srcPath, "main.cpp")
	codeContent := []byte(CODE)
	err = os.WriteFile(codePath, codeContent, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Failed to create main.cpp")
	}

	return nil
}

func CreateProject(name string) error {
	if name == "." {
		ifExist, err := checkIfProjectAlreadyExist()
		if err != nil {
			return err
		}
		if ifExist {
			return errors.New("Project Already Exists in Current Directory")
		}
	}

	projectRootPath, projectSrcPath, projectIncludePath, err := constructProjectPaths(name)
	if err != nil {
		return err
	}
	fmt.Println("Creating Project at", projectSrcPath)

	err = createDirectory(projectSrcPath, projectIncludePath)
	if err != nil {
		return err
	}
	fmt.Println("Project folders created successfully.")

	err = createTemplateFiles(projectRootPath, projectSrcPath)
	if err != nil {
		return err
	}
	fmt.Println("Project initialization successful.")

	return nil
}
