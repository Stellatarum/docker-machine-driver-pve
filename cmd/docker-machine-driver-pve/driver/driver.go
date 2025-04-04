package driver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"slices"
	"strings"

	"github.com/luthermonson/go-proxmox"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/log"
	"github.com/rancher/machine/libmachine/ssh"
	"github.com/rancher/machine/libmachine/state"
)

var _ drivers.Driver = (*Driver)(nil)

// Driver is the implementation of drivers.Driver interface.
type Driver struct {
	*drivers.BaseDriver
	config

	// Cached client for the Proxmox VE.
	pveClient *proxmox.Client

	// Proxmox VE ID of the current machine.
	PVEMachineID *int
}

// Creates a new driver.
func NewDriver(machineName, storePath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: machineName,
			StorePath:   storePath,
		},
		config: config{},
	}
}

// PreCreateCheck implements drivers.Driver.
func (d *Driver) PreCreateCheck() error {
	// Check resource pool
	resourcePool, err := d.getCurrentPVEResourcePool(context.TODO())
	if err != nil {
		return err
	}

	// Check template
	template, err := d.getPVETemplate(context.TODO())
	if err != nil {
		return err
	}

	// Check ISO device
	var (
		isoDeviceConfig string
		isoDeviceFound  bool
	)

	switch {
	case strings.HasPrefix(d.ISODeviceName, "ide"):
		isoDeviceConfig, isoDeviceFound = template.VirtualMachineConfig.MergeIDEs()[d.ISODeviceName]
	case strings.HasPrefix(d.ISODeviceName, "sata"):
		isoDeviceConfig, isoDeviceFound = template.VirtualMachineConfig.MergeSATAs()[d.ISODeviceName]
	case strings.HasPrefix(d.ISODeviceName, "scsi"):
		isoDeviceConfig, isoDeviceFound = template.VirtualMachineConfig.MergeSCSIs()[d.ISODeviceName]
	default:
		return errors.New("only 'ide', 'sata' and 'scsi' devices can be used for cloud-init ISO")
	}

	if !isoDeviceFound {
		return fmt.Errorf("cloud-init ISO device '%s' not found on the template", d.ISODeviceName)
	}

	if !strings.Contains(isoDeviceConfig, "media=cdrom") {
		return errors.New("cloud-init ISO device must be of type media=cdrom")
	}

	// Check network interface
	_, networkInterfaceFound := template.VirtualMachineConfig.MergeNets()[d.NetworkInterfaceName]
	if !networkInterfaceFound {
		return fmt.Errorf("network interface '%s' not found on the template", d.NetworkInterfaceName)
	}

	log.Debugf("Using resource pool '%s'", resourcePool.PoolID)
	log.Debugf("Using template name '%s' on node '%s'", template.Name, template.Node)
	log.Debugf("Using device '%s' for cloud-init ISO", d.ISODeviceName)
	log.Debugf("Using network interface '%s' for IP address", d.NetworkInterfaceName)

	return nil
}

// Create implements drivers.Driver.
func (d *Driver) Create() error {
	log.Info("Generating SSH keys...")

	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return fmt.Errorf("failed to generate SSH key pair: %w", err)
	}

	log.Info("Creating the machine...")

	vmid, err := d.createPVEVirtualMachine(context.TODO())
	if err != nil {
		if vmid > 0 {
			log.Warn("Machine might have been created with ID='%d'", vmid)
		}

		return fmt.Errorf("failed to create machine: %w", err)
	}

	d.PVEMachineID = &vmid

	log.Info("Tagging the machine...")

	machine, err := d.getPVEVirtualMachine(context.TODO(), *d.PVEMachineID)
	if err != nil {
		return fmt.Errorf("failed to retrieve newly created Proxmox VE virtual machine ID='%d': %w", *d.PVEMachineID, err)
	}

	tagTask, err := machine.AddTag(context.TODO(), pveMachineTag)

	if err == nil {
		err = d.waitForPVETaskToSucceed(context.TODO(), tagTask)
	}

	if err != nil {
		return fmt.Errorf("failed to add tag '%s' to Proxmox VE virtual machine ID='%d': %w", pveMachineTag, *d.PVEMachineID, err)
	}

	log.Info("Configuring cloud-init...")

	if err := d.setupCloudinit(context.TODO()); err != nil {
		return err
	}

	log.Info("Starting the machine...")

	if err := d.Start(); err != nil {
		return fmt.Errorf("failed to start the machine: %w", err)
	}

	log.Info("Waiting for cloud-init to finish...")

	if err := d.waitForCloudinit(); err != nil {
		return fmt.Errorf("failed waiting for cloud-init to finish: %w", err)
	}

	log.Info("Cleaning up...")

	if err := d.cleanupCloudinit(context.TODO()); err != nil {
		return err
	}

	return nil
}

// GetState implements drivers.Driver.
func (d *Driver) GetState() (state.State, error) {
	machine, err := d.getCurrentMachine(context.TODO())
	if err != nil {
		return state.Error, err
	}

	if !machine.IsRunning() {
		return state.Stopped, nil
	}

	_, err = machine.AgentOsInfo(context.TODO())
	if err == nil {
		return state.Running, nil
	}

	if !strings.Contains(err.Error(), "500 QEMU guest agent is not running") {
		return state.Error, fmt.Errorf("failed to retrieve Proxmox machine ID='%d' agent status: %w", *d.PVEMachineID, err)
	}

	return state.Starting, nil
}

// GetIP implements drivers.Driver.
func (d *Driver) GetIP() (string, error) {
	machine, err := d.getCurrentMachine(context.TODO())
	if err != nil {
		return "", err
	}

	if !machine.IsRunning() {
		return "", errors.New("machine is powered off")
	}

	networkInterfaceConfiguration, networkInterfaceFound := machine.VirtualMachineConfig.MergeNets()[d.NetworkInterfaceName]
	if !networkInterfaceFound {
		return "", fmt.Errorf("failed to retrieve Proxmox VE machine's ID='%d' network interface configuration for device '%s'", *d.PVEMachineID, d.NetworkInterfaceName)
	}

	networkInterfaceMAC := getMACFromPveNetworkDevice(networkInterfaceConfiguration)
	if networkInterfaceMAC == "" {
		return "", fmt.Errorf(
			"failed to retrieve Proxmox VE machine's ID='%d' network interface MAC for device '%s' with configuration '%s'",
			*d.PVEMachineID,
			d.NetworkInterfaceName,
			networkInterfaceConfiguration,
		)
	}

	networkInterfaceMAC = strings.ToLower(networkInterfaceMAC)

	osNetworkInterfaces, err := machine.AgentGetNetworkIFaces(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to retrieve Proxmox VE machine's ID='%d' network interfaces: %w", d.PVEMachineID, err)
	}

	possibleIPv4s := []string{}
	possibleIPv6s := []string{}

	for _, osNetworkInterface := range osNetworkInterfaces {
		if strings.ToLower(osNetworkInterface.HardwareAddress) != networkInterfaceMAC {
			continue
		}

		for _, address := range osNetworkInterface.IPAddresses {
			parsedAddress := net.ParseIP(address.IPAddress)

			if parsedAddress == nil {
				continue
			}

			if parsedAddress.IsLoopback() || parsedAddress.IsUnspecified() {
				continue
			}

			if address.IPAddressType == "ipv4" {
				possibleIPv4s = append(possibleIPv4s, parsedAddress.String())
			} else {
				possibleIPv6s = append(possibleIPv6s, parsedAddress.String())
			}
		}
	}

	if len(possibleIPv4s) > 0 {
		slices.Sort(possibleIPv4s)
		return possibleIPv4s[0], nil
	}

	if len(possibleIPv6s) > 0 {
		slices.Sort(possibleIPv6s)
		return possibleIPv6s[0], nil
	}

	return "", fmt.Errorf("failed to find Proxmox VE machine's ID='%d' address on interface '%s'", *d.PVEMachineID, d.NetworkInterfaceName)
}

// GetSSHHostname implements drivers.Driver.
func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// Returns a public key path for use with SSH.
func (d *Driver) GetSSHPublicKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}

// GetURL implements drivers.Driver.
func (d *Driver) GetURL() (string, error) {
	address, err := d.GetIP()
	return "tcp://" + net.JoinHostPort(address, "2376"), err
}

// DriverName implements drivers.Driver.
func (d *Driver) DriverName() string {
	return "pve"
}

// Start implements drivers.Driver.
func (d *Driver) Start() error {
	err := d.runTaskOnCurrentMachine(context.TODO(), func(ctx context.Context, vm *proxmox.VirtualMachine) (*proxmox.Task, error) {
		return vm.Start(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed to start the machine: %w", err)
	}

	return nil
}

// Restart implements drivers.Driver.
func (d *Driver) Restart() error {
	err := d.runTaskOnCurrentMachine(context.TODO(), func(ctx context.Context, vm *proxmox.VirtualMachine) (*proxmox.Task, error) {
		return vm.Reboot(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed to restart the machine: %w", err)
	}

	return nil
}

// Stop implements drivers.Driver.
func (d *Driver) Stop() error {
	err := d.runTaskOnCurrentMachine(context.TODO(), func(ctx context.Context, vm *proxmox.VirtualMachine) (*proxmox.Task, error) {
		return vm.Shutdown(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed to stop the machine: %w", err)
	}

	return nil
}

// Kill implements drivers.Driver.
func (d *Driver) Kill() error {
	err := d.runTaskOnCurrentMachine(context.TODO(), func(ctx context.Context, vm *proxmox.VirtualMachine) (*proxmox.Task, error) {
		return vm.Stop(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed to kill the machine: %w", err)
	}

	return nil
}

// Remove implements drivers.Driver.
func (d *Driver) Remove() error {
	err := d.Kill()
	if err != nil {
		return err
	}

	err = d.runTaskOnCurrentMachine(context.TODO(), func(ctx context.Context, vm *proxmox.VirtualMachine) (*proxmox.Task, error) {
		return vm.Delete(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed to remove the machine: %w", err)
	}

	return nil
}
