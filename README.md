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
- Optional virtio support for disk and network devices

## Notes

- This tool requires sudo privileges to copy the image file and interact with libvirt.
- The tool does not start the VM after registration. You can start it manually using `virsh start your-vm-name`.
- Make sure you have enough disk space in `/var/lib/libvirt/images` before running the tool.
- The VM's hostname is automatically set to the VM name specified as an argument, unless the `-no-hostname-change` option is used.
- Virtio drivers can significantly improve VM performance, but ensure your guest OS supports them before enabling.

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
- `-virtio-disk`: Use virtio for disk device
- `-virtio-network`: Use virtio for network device

### Examples

Standard VM creation:
```
$ ./kvm-vm-register -image ../debvirt-image-kit/output/debian-12.7.0-amd64 -memory 2048 -vcpus 2 bookworm64
```

VM creation with virtio for both disk and network:
```
$ ./kvm-vm-register -image ../debvirt-image-kit/output/debian-12.7.0-amd64 -memory 2048 -vcpus 2 -virtio-disk -virtio-network bookworm64
```

VM creation for FreeBSD (skipping hostname change):
```
$ ./kvm-vm-register -image /path/to/freebsd/image.qcow2 -memory 2048 -vcpus 2 -no-hostname-change freebsd-vm
```

These commands will:
1. Copy the specified image to `/var/lib/libvirt/images/<vm_name>.qcow2`
2. Set the VM's hostname to the specified name (unless `-no-hostname-change` is used)
3. Generate an XML configuration for the VM with the specified parameters
4. Register the VM without starting it

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
