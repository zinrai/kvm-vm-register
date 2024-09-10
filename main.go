package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func checkCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func checkRequiredCommands() error {
	requiredCommands := []string{"sudo", "virt-install", "virsh", "dd"}
	missingCommands := []string{}

	for _, cmd := range requiredCommands {
		if !checkCommand(cmd) {
			missingCommands = append(missingCommands, cmd)
		}
	}

	if len(missingCommands) > 0 {
		return fmt.Errorf("Missing required commands: %v", missingCommands)
	}

	return nil
}

func copyImageWithDD(src, dst string) error {
	cmd := exec.Command("sudo", "dd", fmt.Sprintf("if=%s", src), fmt.Sprintf("of=%s", dst), "bs=4M", "status=progress")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func checkVMExists(name string) bool {
	cmd := exec.Command("sudo", "virsh", "dominfo", name)
	return cmd.Run() == nil
}

func main() {
	imagePath := flag.String("image", "", "Path to the VM image")
	vmName := flag.String("name", "", "Name of the new VM")
	memory := flag.Int("memory", 1024, "Memory size in MB")
	vcpus := flag.Int("vcpus", 1, "Number of virtual CPUs")
	network := flag.String("network", "network=default", "Network configuration for virt-install")

	flag.Parse()

	if *imagePath == "" || *vmName == "" {
		fmt.Println("Error: Image path and VM name are required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := checkRequiredCommands(); err != nil {
		log.Fatalf("Error: %v\nPlease install the missing commands and try again.", err)
	}

	absImagePath, err := filepath.Abs(*imagePath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}
	if _, err := os.Stat(absImagePath); os.IsNotExist(err) {
		log.Fatalf("Image file does not exist: %s", absImagePath)
	}

	newDiskPath := filepath.Join("/var/lib/libvirt/images", *vmName+".qcow2")

	// Check if the image file already exists
	if _, err := os.Stat(newDiskPath); err == nil {
		log.Fatalf("Image file already exists at %s. Please choose a different VM name or remove the existing image.", newDiskPath)
	}

	// Check if VM with the same name already exists
	if checkVMExists(*vmName) {
		log.Fatalf("A VM with the name '%s' already exists. Please choose a different name.", *vmName)
	}

	fmt.Printf("Copying image file to %s...\n", newDiskPath)
	if err := copyImageWithDD(absImagePath, newDiskPath); err != nil {
		log.Fatalf("Failed to copy image file: %v", err)
	}
	fmt.Println("Image file copied successfully.")

	virtInstallArgs := []string{
		"virt-install",
		"--name", *vmName,
		"--memory", fmt.Sprintf("%d", *memory),
		"--vcpus", fmt.Sprintf("%d", *vcpus),
		"--disk", newDiskPath,
		"--import",
		"--os-variant", "generic",
		"--network", *network,
		"--print-xml",
	}

	virtInstallCmd := exec.Command("sudo", virtInstallArgs...)
	xmlOutput, err := virtInstallCmd.Output()
	if err != nil {
		log.Fatalf("Failed to generate VM XML: %v\nCommand: %s", err, virtInstallCmd.String())
	}

	tmpfile, err := os.CreateTemp("", "vm-*.xml")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(xmlOutput); err != nil {
		log.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatalf("Failed to close temporary file: %v", err)
	}

	defineCmd := exec.Command("sudo", "virsh", "define", tmpfile.Name())
	if output, err := defineCmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to define VM: %v\nCommand: %s\nOutput: %s", err, defineCmd.String(), output)
	}

	fmt.Printf("VM '%s' registered successfully\n", *vmName)
}
