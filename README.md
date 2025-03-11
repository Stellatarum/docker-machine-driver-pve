# Docker/Rancher Machine driver for Proxmox VE

This is a Docker/[Rancher Machine](https://github.com/rancher/machine) driver for [Proxmox VE](https://www.proxmox.com/en/proxmox-virtual-environment/overview).

## Installation

```sh
go install github.com/stellatarum/docker-machine-driver-pve/cmd/docker-machine-driver-pve
```

## Template requirements

This driver requires a Proxmox VE template with:

* [`qemu-guest-agent`](https://pve.proxmox.com/wiki/Qemu-guest-agent),
* cloud-init initialization enabled,
* empty CD/DVD drive (**NOT** PVE's CloudInit Drive) on IDE, SATA or SCSI bus,
* DHCP enabled network interface.

The template must be placed in the same resource pool where the machines will be deployed (i.e. `--pve-resource-pool`).

You can use [sample Ubuntu Server template](deploy/templates/ubuntu-server) for development and testing.

## Configuration

| Flag                      | Environment variable    | Default value                      | Description                                                                     |
| ------------------------- | ----------------------- | ---------------------------------- | ------------------------------------------------------------------------------- |
| `--pve-url`               | `PVE_URL`               | N/A (required)                     | Proxmox VE URL (e.g. `https://<PROXMOX VE ADDRESS>:8006`)                       |
| `--pve-insecure-tls`      | `PVE_INSECURE_TLS`      | `false`                            | Disables Proxmox VE TLS certificate verification                                |
| `--pve-token-id`          | `PVE_TOKEN_ID`          | N/A (required)                     | Proxmox VE API Token ID (including username and realm, e.g. `root@pam!rancher`) |
| `--pve-token-secret`      | `PVE_TOKEN_SECRET`      | N/A (required)                     | Proxmox VE API Token secret                                                     |
| `--pve-resource-pool`     | `PVE_RESOURCE_POOL`     | N/A (required)                     | Proxmox VE Resource Pool name                                                   |
| `--pve-template`          | `PVE_TEMPLATE`          | N/A (required)                     | ID of the Proxmox VE template                                                   |
| `--pve-iso-device`        | `PVE_ISO_DEVICE`        | N/A (required)                     | Bus/Device of the CD/DVD Drive to mount cloud-init ISO to (e.g. `scsi1`)        |
| `--pve-network-interface` | `PVE_NETWORK_INTERFACE` | N/A (required)                     | Network interface to read machine's IP address form                             |
| `--pve-ssh-user`          | `PVE_SSH_USER`          | `service`                          | Username for the SSH user that will be created via cloud-init                   |
| `--pve-ssh-port`          | `PVE_SSH_PORT`          | `22`                               | Port to use when connecting to the machine via SSH                              |

## Contributing

See [DEVELOPMENT.md](./docs/DEVELOPMENT.md) for development guidelines.

## License

Distributed under the Apache License 2.0. See [LICENSE](./LICENSE) for more information.

## Contact

You can contact us by e-mail: [contact@stellatarum.com](mailto:contact@stellatarum.com)
