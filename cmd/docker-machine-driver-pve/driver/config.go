package driver

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/mcnflag"
)

// Available flags.
const (
	flagURL              = "pve-url"
	flagInsecureTLS      = "pve-insecure-tls"
	flagTokenID          = "pve-token-id" //nolint:gosec // False-positive
	flagTokenSecret      = "pve-token-secret"
	flagResourcePool     = "pve-resource-pool"
	flagTemplateID       = "pve-template"
	flagISODevice        = "pve-iso-device"
	flagNetworkInterface = "pve-network-interface"
	flagSSHUser          = "pve-ssh-user"
	flagSSHPort          = "pve-ssh-port"
	flagProcessorSockets = "pve-processor-sockets"
	flagProcessorCores   = "pve-processor-cores"
	flagMemory           = "pve-memory"
	flagMemoryBalloon    = "pve-memory-balloon"
)

// Default values for flags.
const (
	defaultSSHUser = "service"
	defaultSSHPort = 22
)

// Driver's configuration.
type config struct {
	// Proxmox VE URL (e.g. 'https://<PROXMOX VE ADDRESS>:8006').
	URL string

	// Disables Proxmox VE TLS certificate verification.
	InsecureTLS bool

	// Proxmox VE API Token ID (including username and realm, e.g. 'root@pam!rancher').
	TokenID string

	// Proxmox VE API Token secret.
	TokenSecret string

	// Proxmox VE Resource Pool name.
	ResourcePoolName string

	// ID of the Proxmox VE template.
	TemplateID int

	// Bus/Device of the CD/DVD Drive to mount cloud-init ISO to (e.g. 'scsi1').
	ISODeviceName string

	// Bus/Device of the network interface to read machine's IP address from (e.g. 'net0').
	NetworkInterfaceName string

	// If set, number of processor sockets to configure for the machine.
	ProcessorSockets *int

	// If set, number of processor cores to configure for the machine.
	ProcessorCores *int

	// If set, amount of memory in MiB to configure for the machine.
	Memory *int

	// If set, minimum amount of memory in MiB to configure for the machine.
	// If set to 0, disables memory ballooning.
	MemoryBalloon *int
}

// GetCreateFlags implements drivers.Driver.
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			Name:   flagURL,
			EnvVar: flagEnvVarFromFlagName(flagURL),
			Usage:  "Proxmox VE URL (e.g. 'https://<PROXMOX VE ADDRESS>:8006')",
		},
		mcnflag.BoolFlag{
			Name:   flagInsecureTLS,
			EnvVar: flagEnvVarFromFlagName(flagInsecureTLS),
			Usage:  "Disables Proxmox VE TLS certificate verification",
		},
		mcnflag.StringFlag{
			Name:   flagTokenID,
			EnvVar: flagEnvVarFromFlagName(flagTokenID),
			Usage:  "Proxmox VE API Token ID (including username and realm, e.g. 'root@pam!rancher')",
		},
		mcnflag.StringFlag{
			Name:   flagTokenSecret,
			EnvVar: flagEnvVarFromFlagName(flagTokenSecret),
			Usage:  "Proxmox VE API Token secret",
		},
		mcnflag.StringFlag{
			Name:   flagResourcePool,
			EnvVar: flagEnvVarFromFlagName(flagResourcePool),
			Usage:  "Proxmox VE Resource Pool name",
		},
		mcnflag.IntFlag{
			Name:   flagTemplateID,
			EnvVar: flagEnvVarFromFlagName(flagTemplateID),
			Usage:  "ID of the Proxmox VE template",
		},
		mcnflag.StringFlag{
			Name:   flagISODevice,
			EnvVar: flagEnvVarFromFlagName(flagISODevice),
			Usage:  "Bus/Device of the CD/DVD Drive to mount cloud-init ISO to (e.g. 'scsi1')",
		},
		mcnflag.StringFlag{
			Name:   flagNetworkInterface,
			EnvVar: flagEnvVarFromFlagName(flagNetworkInterface),
			Usage:  "Bus/Device of the network interface to read machine's IP address from (e.g. 'net0')",
		},
		mcnflag.StringFlag{
			Name:   flagSSHUser,
			EnvVar: flagEnvVarFromFlagName(flagSSHUser),
			Usage:  fmt.Sprintf("Username for the SSH user that will be created via cloud-init, defaults to '%s'", defaultSSHUser),
		},
		mcnflag.IntFlag{
			Name:   flagSSHPort,
			EnvVar: flagEnvVarFromFlagName(flagSSHPort),
			Usage:  fmt.Sprintf("Port to use when connecting to the machine via SSH, defaults to '%d'", defaultSSHPort),
		},
		mcnflag.StringFlag{
			Name:   flagProcessorSockets,
			EnvVar: flagEnvVarFromFlagName(flagProcessorSockets),
			Usage:  "If set, number of processor sockets to configure for the machine.",
		},
		mcnflag.StringFlag{
			Name:   flagProcessorCores,
			EnvVar: flagEnvVarFromFlagName(flagProcessorCores),
			Usage:  "If set, number of processor cores to configure for the machine.",
		},
		mcnflag.StringFlag{
			Name:   flagMemory,
			EnvVar: flagEnvVarFromFlagName(flagMemory),
			Usage:  "If set, amount of memory in MiB to configure for the machine.",
		},
		mcnflag.StringFlag{
			Name:   flagMemoryBalloon,
			EnvVar: flagEnvVarFromFlagName(flagMemoryBalloon),
			Usage:  "If set, minimum amount of memory in MiB to configure for the machine. If set to 0, disables memory ballooning.",
		},
	}
}

// SetConfigFromFlags implements drivers.Driver.
//
//nolint:cyclop,gocyclo
func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	d.URL = opts.String(flagURL)
	if d.URL == "" {
		return fmt.Errorf("flag '--%s' is required", flagURL)
	}

	if _, err := url.Parse(d.URL); err != nil {
		return fmt.Errorf("failed to parse Proxmox VE URL (flag '--%s'): %w", flagURL, err)
	}

	d.InsecureTLS = opts.Bool(flagInsecureTLS)

	d.TokenID = opts.String(flagTokenID)
	if d.TokenID == "" {
		return fmt.Errorf("flag '--%s' is required", flagTokenID)
	}

	d.TokenSecret = opts.String(flagTokenSecret)
	if d.TokenSecret == "" {
		return fmt.Errorf("flag '--%s' is required", flagTokenSecret)
	}

	d.ResourcePoolName = opts.String(flagResourcePool)
	if d.ResourcePoolName == "" {
		return fmt.Errorf("flag '--%s' is required", flagResourcePool)
	}

	d.TemplateID = opts.Int(flagTemplateID)
	if d.TemplateID <= 0 {
		return fmt.Errorf("flag '--%s' is required and must be >= 0", flagTemplateID)
	}

	d.ISODeviceName = strings.ToLower(opts.String(flagISODevice))
	if d.ISODeviceName == "" {
		return fmt.Errorf("flag '--%s' is required", flagISODevice)
	}

	d.NetworkInterfaceName = opts.String(flagNetworkInterface)
	if d.NetworkInterfaceName == "" {
		return fmt.Errorf("flag '--%s' is required", flagNetworkInterface)
	}

	d.SSHUser = opts.String(flagSSHUser)
	if d.SSHUser == "" {
		d.SSHUser = defaultSSHUser
	}

	d.SSHPort = opts.Int(flagSSHPort)
	if d.SSHPort == 0 {
		d.SSHPort = defaultSSHPort
	} else if d.SSHPort < 0 {
		return fmt.Errorf("flag '--%s' must be > 0", flagSSHPort)
	}

	var err error

	if d.ProcessorSockets, err = parseStringFlagToInt(opts.String(flagProcessorSockets)); err != nil {
		return fmt.Errorf("failed to parse '--%s': %w", flagProcessorSockets, err)
	} else if d.ProcessorSockets != nil && *d.ProcessorSockets < 1 {
		return fmt.Errorf("flag '--%s' must be >= 1", flagProcessorSockets)
	}

	if d.ProcessorCores, err = parseStringFlagToInt(opts.String(flagProcessorCores)); err != nil {
		return fmt.Errorf("failed to parse '--%s': %w", flagProcessorCores, err)
	} else if d.ProcessorCores != nil && *d.ProcessorCores < 1 {
		return fmt.Errorf("flag '--%s' must be >= 1", flagProcessorCores)
	}

	if d.Memory, err = parseStringFlagToInt(opts.String(flagMemory)); err != nil {
		return fmt.Errorf("failed to parse '--%s': %w", flagMemory, err)
	} else if d.Memory != nil && *d.Memory < 1 {
		return fmt.Errorf("flag '--%s' must be >= 1", flagMemory)
	}

	if d.MemoryBalloon, err = parseStringFlagToInt(opts.String(flagMemoryBalloon)); err != nil {
		return fmt.Errorf("failed to parse '--%s': %w", flagMemoryBalloon, err)
	} else if d.MemoryBalloon != nil && *d.MemoryBalloon < 0 {
		return fmt.Errorf("flag '--%s' must be >= 1; set to 0 to disable", flagMemoryBalloon)
	}

	// Default memory/memory ballon to the other one if it's set
	if d.Memory != nil && d.MemoryBalloon == nil {
		d.MemoryBalloon = d.Memory
	}

	if d.MemoryBalloon != nil && *d.MemoryBalloon != 0 && d.Memory == nil {
		d.Memory = d.MemoryBalloon
	}

	// Balloon target can not be higher than total memory.
	if d.Memory != nil && d.MemoryBalloon != nil && *d.MemoryBalloon > *d.Memory {
		return fmt.Errorf("flag '--%s' must be <= than flag '--%s'", flagMemoryBalloon, flagMemory)
	}

	return nil
}

// Creates flag's EnvVar from it's name.
func flagEnvVarFromFlagName(name string) string {
	return strings.ToUpper(
		strings.ReplaceAll(
			name,
			"-",
			"_",
		),
	)
}

// Parses string flag to integer. Returns nil if the flag was unset/empty.
func parseStringFlagToInt(value string) (*int, error) {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		//nolint:nilnil
		return nil, nil
	}

	numberValue, err := strconv.Atoi(trimmedValue)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to int: %w", err)
	}

	return &numberValue, nil
}
