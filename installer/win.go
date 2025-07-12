//go:build windows
// +build windows

package main

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
)

func installWindows() error {
	fmt.Println("Installing application for Windows...")
	programFiles, err := windows.KnownFolderPath(windows.FOLDERID_ProgramFiles, 0)
	if err != nil {
		return fmt.Errorf("failed to get Program Files path: %v", err)
	}
	appDir := filepath.Join(programFiles, AppName)
	fmt.Printf("Creating application directory at: %s\n", appDir)
	if err := os.MkdirAll(appDir, defaultInstallPerms); err != nil {
		return fmt.Errorf("failed to create app directory: %v", err)
	}
	binaryDest := filepath.Join(appDir, BinaryName)
	fmt.Printf("Copying executable to: %s\n", binaryDest)
	if err := copyEmbeddedFile(installFiles, ("install" + "/" + BinaryName), binaryDest, defaultInstallPerms); err != nil {
		return fmt.Errorf("failed to copy binary: %v", err)
	}
	libDir := filepath.Join(appDir, "lib")
	if entries, err := installFiles.ReadDir("install/lib"); err == nil && len(entries) > 0 {
		fmt.Printf("Copying application libraries to: %s\n", libDir)
		if err := copyEmbeddedDir(installFiles, "install/lib", libDir, defaultInstallPerms); err != nil {
			return fmt.Errorf("failed to copy libraries: %v", err)
		}
	}
	fmt.Println("Creating Start Menu shortcut...")
	startMenu, err := windows.KnownFolderPath(windows.FOLDERID_StartMenu, 0)
	if err != nil {
		return fmt.Errorf("failed to get Start Menu path: %v", err)
	}
	shortcutPath := filepath.Join(startMenu, "Programs", AppName, AppName+".lnk")
	fmt.Printf("Creating shortcut at: %s\n", shortcutPath)
	if err := os.MkdirAll(filepath.Dir(shortcutPath), defaultInstallPerms); err != nil {
		return fmt.Errorf("failed to create Start Menu directory: %v", err)
	}
	if err := createWindowsShortcut(binaryDest, shortcutPath); err != nil {
		return fmt.Errorf("failed to create shortcut: %v", err)
	}
	return nil
}

func createWindowsShortcut(targetPath, shortcutPath string) error {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()
	unknown, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return fmt.Errorf("failed to create WScript.Shell object: %v", err)
	}
	defer unknown.Release()
	shell, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("failed to query WScript.Shell interface: %v", err)
	}
	defer shell.Release()
	cs, err := oleutil.CallMethod(shell, "CreateShortcut", shortcutPath)
	if err != nil {
		return fmt.Errorf("failed to create shortcut: %v", err)
	}
	shortcut := cs.ToIDispatch()
	defer shortcut.Release()
	if _, err := oleutil.PutProperty(shortcut, "TargetPath", targetPath); err != nil {
		return fmt.Errorf("failed to set TargetPath: %v", err)
	}
	workingDir := filepath.Dir(targetPath)
	if _, err := oleutil.PutProperty(shortcut, "WorkingDirectory", workingDir); err != nil {
		return fmt.Errorf("failed to set WorkingDirectory: %v", err)
	}
	if _, err := oleutil.PutProperty(shortcut, "Description", "Shortcut for "+AppName); err != nil {
		return fmt.Errorf("failed to set Description: %v", err)
	}
	if _, err := oleutil.PutProperty(shortcut, "IconLocation", targetPath+",0"); err != nil {
		return fmt.Errorf("failed to set IconLocation: %v", err)
	}
	if _, err := oleutil.CallMethod(shortcut, "Save"); err != nil {
		return fmt.Errorf("failed to save shortcut: %v", err)
	}
	return nil
}

func uninstallWindows() error {
	fmt.Println("Removing Windows application...")
	programFiles, err := windows.KnownFolderPath(windows.FOLDERID_ProgramFiles, 0)
	if err != nil {
		return fmt.Errorf("failed to get Program Files path: %v", err)
	}
	appDir := filepath.Join(programFiles, AppName)
	fmt.Printf("Removing application directory: %s\n", appDir)
	if err := os.RemoveAll(appDir); err != nil {
		return fmt.Errorf("failed to remove app directory: %v", err)
	}
	startMenu, err := windows.KnownFolderPath(windows.FOLDERID_StartMenu, 0)
	if err != nil {
		return fmt.Errorf("failed to get Start Menu path: %v", err)
	}
	shortcutDir := filepath.Join(startMenu, "Programs", AppName)
	fmt.Printf("Removing Start Menu folder: %s\n", shortcutDir)
	if err := os.RemoveAll(shortcutDir); err != nil {
		return fmt.Errorf("failed to remove Start Menu folder: %v", err)
	}
	return nil
}

func isAdmin() bool {
	token, err := windows.OpenCurrentProcessToken()
	if err != nil {
		return false
	}
	defer token.Close()

	elevated := token.IsElevated()
	return elevated
}
