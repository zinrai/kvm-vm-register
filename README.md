# kvm-vm-register

This tool allows you to easily register a new KVM (Kernel-based Virtual Machine) virtual machine using an existing image file. It copies the image to the libvirt images directory and registers the VM without starting it.

Use [debvirt-image-kit](https://github.com/zinrai/debvirt-image-kit) to create virtual machine images.

## Features

- Copies the specified VM image to `/var/lib/libvirt/images`
- Generates XML configuration for the new VM
- Registers the VM with libvirt without starting it
- Customizable VM parameters (memory, vCPUs, network)
- Creates an independent qcow2 image
- Automatically sets the VM's hostname to match the VM name (can be skipped for unsupported OSes)

## Notes

- This tool requires sudo privileges to copy the image file and interact with libvirt.
- The tool does not start the VM after registration. You can start it manually using `virsh start your-vm-name`.
- Make sure you have enough disk space in `/var/lib/libvirt/images` before running the tool.
- The VM's hostname is automatically set to the VM name specified as an argument, unless the `-no-hostname-change` option is used.

## Requirements

The following commands must be available in the system PATH:

- `virt-install`
- `virsh`
- `qemu-img`
- `virt-customize`
- `sudo`

## Installation

Build the tool:

```
$ go build
```

## Usage

Run the tool with the following command:

```
$ ./kvm-vm-register [options] -image /path/to/your/image.qcow2 vm_name
```

### Options

- `-image`: Path to the VM image (required)
- `-memory`: Memory size in MB (default: 1024)
- `-vcpus`: Number of virtual CPUs (default: 1)
- `-network`: Network configuration for virt-install (default: "network=default")
- `-no-hostname-change`: Skip hostname change (for FreeBSD or other unsupported OSes)

### Example

For a standard Linux VM:
```
$ ./kvm-vm-register -image ../debvirt-image-kit/output/debian-12.7.0-amd64 -memory 2048 -vcpus 2 bookworm64
```

This command will:
1. Copy the Debian image to `/var/lib/libvirt/images/bookworm64.qcow2`
2. Set the VM's hostname to "bookworm64"
3. Generate an XML configuration for a VM named "bookworm64" with 2048MB of RAM and 2 vCPUs
4. Register the VM without starting it

For a FreeBSD or other unsupported OS:
```
$ ./kvm-vm-register -image /path/to/freebsd/image.qcow2 -memory 2048 -vcpus 2 -no-hostname-change freebsd-vm
```

This command will perform the same steps as above, but skip the hostname change step.

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
