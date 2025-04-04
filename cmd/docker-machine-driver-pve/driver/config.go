package driver

import (
	"fmt"
	"net/url"
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

	// Network interface to read machine's IP address form (e.g. 'net0').
	NetworkInterfaceName string
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
			Usage:  "Network interface to read machine's IP address form (e.g. 'net0')",
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
	}
}

// SetConfigFromFlags implements drivers.Driver.
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
