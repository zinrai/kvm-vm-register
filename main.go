package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func checkCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func checkRequiredCommands() error {
	requiredCommands := []string{"sudo", "virt-install", "virsh", "qemu-img", "virt-customize"}
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

func createIndependentImage(src, dst string) error {
	args := []string{"qemu-img", "convert", "-O", "qcow2", src, dst}

	cmd := exec.Command("sudo", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create independent image: %v, output: %s", err, output)
	}

	return nil
}

func checkVMExists(name string) bool {
	cmd := exec.Command("sudo", "virsh", "dominfo", name)
	return cmd.Run() == nil
}

func setVMHostname(imagePath, hostname string) error {
	cmd := exec.Command("sudo", "virt-customize", "-a", imagePath, "--hostname", hostname)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set VM hostname: %v, output: %s", err, output)
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] vm_name\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	imagePath := flag.String("image", "", "Path to the VM image")
	memory := flag.Int("memory", 1024, "Memory size in MB")
	vcpus := flag.Int("vcpus", 1, "Number of virtual CPUs")
	network := flag.String("network", "network=default", "Network configuration for virt-install")
	noHostnameChange := flag.Bool("no-hostname-change", false, "Skip hostname change (for FreeBSD or other unsupported OSes)")
	virtioDisk := flag.Bool("virtio-disk", false, "Use virtio for disk device")
	virtioNetwork := flag.Bool("virtio-network", false, "Use virtio for network device")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 || *imagePath == "" {
		fmt.Println("Error: VM name and image path are required")
		flag.Usage()
		os.Exit(1)
	}

	vmName := args[0]

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

	newDiskPath := filepath.Join("/var/lib/libvirt/images", vmName+".qcow2")

	// Check if the image file already exists
	if _, err := os.Stat(newDiskPath); err == nil {
		log.Fatalf("Image file already exists at %s. Please choose a different VM name or remove the existing image.", newDiskPath)
	}

	// Check if VM with the same name already exists
	if checkVMExists(vmName) {
		log.Fatalf("A VM with the name '%s' already exists. Please choose a different name.", vmName)
	}

	fmt.Printf("Creating independent image file at %s...\n", newDiskPath)
	if err := createIndependentImage(absImagePath, newDiskPath); err != nil {
		log.Fatalf("Failed to create independent image file: %v", err)
	}
	fmt.Println("Independent image file created successfully.")

	if !*noHostnameChange {
		fmt.Printf("Setting VM hostname to %s...\n", vmName)
		if err := setVMHostname(newDiskPath, vmName); err != nil {
			log.Fatalf("Failed to set VM hostname: %v", err)
		}
		fmt.Println("VM hostname set successfully.")
	} else {
		fmt.Println("Skipping hostname change as requested.")
	}

	diskOption := newDiskPath
	networkOption := *network

	if *virtioDisk {
		diskOption = fmt.Sprintf("%s,bus=virtio", newDiskPath)
	}

	if *virtioNetwork {
		if networkOption == "network=default" {
			networkOption = "network=default,model=virtio"
		} else {
			networkOption += ",model=virtio"
		}
	}

	virtInstallArgs := []string{
		"virt-install",
		"--name", vmName,
		"--memory", strconv.Itoa(*memory),
		"--vcpus", strconv.Itoa(*vcpus),
		"--disk", diskOption,
		"--import",
		"--os-variant", "generic",
		"--network", networkOption,
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

	fmt.Printf("VM '%s' registered successfully\n", vmName)
}
