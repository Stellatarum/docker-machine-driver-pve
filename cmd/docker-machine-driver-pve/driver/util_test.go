package driver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getMACFromPveNetworkDevice(t *testing.T) {
	tests := map[string]string{
		"":                                       "",
		",":                                      "",
		"bridge=vmbr1":                           "",
		"e1000=BC:24:11:45:CD:E8,bridge=vmbr0":   "BC:24:11:45:CD:E8",
		"e1000e=BC:24:11:E1:F0:71,bridge=vmbr3":  "BC:24:11:E1:F0:71",
		"rtl8139=BC:24:11:18:BB:08,bridge=vmbr1": "BC:24:11:18:BB:08",
		"virtio=BC:24:11:87:63:EC,bridge=vmbr1":  "BC:24:11:87:63:EC",
		"vmxnet3=BC:24:11:EB:05:E9,bridge=vmbr4": "BC:24:11:EB:05:E9",
	}

	for deviceConfiguration, expectedAddress := range tests {
		t.Run(
			deviceConfiguration,
			func(t *testing.T) {
				require.Equal(
					t,
					expectedAddress,
					getMACFromPveNetworkDevice(deviceConfiguration),
				)
			},
		)
	}
}
