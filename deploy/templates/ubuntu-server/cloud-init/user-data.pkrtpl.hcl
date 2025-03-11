#cloud-config
autoinstall:
  version: 1
  refresh-installer:
    update: false
  source:
    id: "ubuntu-server-minimal"
    search_drivers: false
  timezone: "Etc/UTC"
  locale: "en_US.UTF-8"
  keyboard:
    layout: "us"
  storage:
    swap:
      size: 0
    config:
      - id: "disk-sda"
        type: "disk"
        path: "/dev/sda"
        ptable: "gpt"
        preserve: false
        grub_device: true
      - id: "sda1"
        type: "partition"
        device: "disk-sda"
        number: 1
        preserve: false
        grub_device: false
        size: 1048576
        flag: "bios_grub"
        offset: 1048576
      - id: "sda2"
        type: "partition"
        device: "disk-sda"
        number: 2
        preserve: false
        grub_device: false
        size: -1
        wipe: "superblock"
        offset: 2097152
      - id: "format-0"
        type: "format"
        volume: "sda2"
        fstype: "xfs"
        preserve: false
      - id: "mount-0"
        type: "mount"
        device: "format-0"
        path: /
  updates: "all"
  codecs:
    install: false
  drivers:
    install: false
  oem:
    install: false
  apt:
    sources:
      docker.list:
        source: "deb [arch=amd64] https://download.docker.com/linux/ubuntu $RELEASE stable"
        keyid: 9DC858229FC7DD38854AE2D88D81803C0EBFCD88
  packages:
    - "qemu-guest-agent"
    - "docker-ce"
    - "docker-ce-cli"
    - "containerd.io"
    - "docker-buildx-plugin"
    - "docker-compose-plugin"
  user-data:
    hostname: "${hostname}"
    users:
      - name: "${ssh_username}"
        groups: ["adm", "sudo"]
        lock_passwd: true
        sudo: "ALL=(ALL) NOPASSWD:ALL"
        shell: "/bin/bash"
        ssh_authorized_keys:
          - "${ssh_authorized_key}"
  ssh:
    install-server: true
    allow-pw: false
