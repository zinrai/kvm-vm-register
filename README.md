# kvm-vm-register

This tool allows you to easily register a new KVM (Kernel-based Virtual Machine) virtual machine using an existing image file. It copies the image to the libvirt images directory and registers the VM without starting it.

Use [debvirt-image-kit](https://github.com/zinrai/debvirt-image-kit) to create virtual machine images.

## Features

- Copies the specified VM image to `/var/lib/libvirt/images`
- Generates XML configuration for the new VM
- Registers the VM with libvirt without starting it
- Customizable VM parameters (name, memory, vCPUs, network)
- Creates an independent qcow2 image, allowing for disk size expansion
- Automatically sets the VM's hostname to match the VM name

## Notes

- This tool requires sudo privileges to copy the image file and interact with libvirt.
- The tool does not start the VM after registration. You can start it manually using `virsh start your-vm-name`.
- Make sure you have enough disk space in `/var/lib/libvirt/images` before running the tool.
- The VM's hostname is automatically set to the name specified with the `-name` option.

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
$ ./kvm-vm-register -image /path/to/your/image.qcow2 -name your-vm-name [options]
```

### Example

```
$ ./kvm-vm-register -image ../debvirt-image-kit/output/debian-12.7.0-amd64 -name bookworm64 -memory 2048 -vcpus 2 -disk-size 20
```

This command will:
1. Copy the Debian image to `/var/lib/libvirt/images/bookworm64.qcow2`
2. Expand the disk size to 20GB
3. Set the VM's hostname to "bookworm64"
4. Generate an XML configuration for a VM named "bookworm64" with 2048MB of RAM and 2 vCPUs
5. Register the VM without starting it

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
## License
