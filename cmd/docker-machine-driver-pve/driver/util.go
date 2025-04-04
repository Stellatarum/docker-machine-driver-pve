package driver

import (
	"slices"
	"strings"
)

func getMACFromPveNetworkDevice(device string) string {
	models := []string{
		"e1000",
		"e1000e",
		"rtl8139",
		"virtio",
		"vmxnet3",
	}

	for _, param := range strings.Split(device, ",") {
		//nolint:mnd
		values := strings.SplitN(param, "=", 2)

		//nolint:mnd
		if len(values) != 2 {
			continue
		}

		if slices.Contains(models, values[0]) {
			return values[1]
		}
	}

	return ""
}
