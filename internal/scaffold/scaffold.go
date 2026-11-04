package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const GITIGNORE = `build/
compiler_commands.json
.cache/
`

const CODE = `#include <iostream>

int main() {
	std::cout << "Hello from ViteCpp" << std::endl;
	return 0;
}
`

const CMAKE = `cmake_minimum_required(VERSION 3.20)

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

func projectAlreadyExist(name string) error {
	entries, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read contents of current directorys %w", err)
	}

	for _, entry := range entries {
		if entry.Name() == "CMakeLists.txt" && name == "." {
			return fmt.Errorf("project already initialized in current directory")
		}
		if entry.Name() == name {
			return fmt.Errorf("directory with project Name already exists")
		}
	}

	return nil
}

func constructProjectPaths(name string) (string, string, string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get current directory name %w", err)
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
		return fmt.Errorf("failed to create new directory %w", err)
	}

	err = os.MkdirAll(includePath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create new directory %w", err)
	}
	return nil
}

func createTemplateFiles(rootPath string, srcPath string) error {
	cMakeFilePath := filepath.Join(rootPath, "CMakeLists.txt")
	projectName := filepath.Base(rootPath)
	updatedCmake := strings.ReplaceAll(CMAKE, "{{PROJECT_NAME}}", projectName)
	cMakeContent := []byte(updatedCmake)
	err := os.WriteFile(cMakeFilePath, cMakeContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to create CMakeLists.txt %w", err)
	}

	gitIgnoreFilePath := filepath.Join(rootPath, ".gitignore")
	gitIgnoreContent := []byte(GITIGNORE)
	err = os.WriteFile(gitIgnoreFilePath, gitIgnoreContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to create .gitignore %w", err)
	}

	codePath := filepath.Join(srcPath, "main.cpp")
	codeContent := []byte(CODE)
	err = os.WriteFile(codePath, codeContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to create main.cpp %w", err)
	}

	return nil
}

func CreateProject(name string) error {
	err := projectAlreadyExist(name)
	if err != nil {
		return err
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
