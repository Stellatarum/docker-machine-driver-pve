module github.com/stellatarum/docker-machine-driver-pve

go 1.24.1

replace github.com/docker/docker => github.com/moby/moby v1.4.2-0.20170731201646-1009e6a40b29

require (
	github.com/luthermonson/go-proxmox v0.2.1
	github.com/rancher/machine v0.15.0-rancher126
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.36.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/buger/goterm v1.0.4 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/diskfs/go-diskfs v1.2.0 // indirect
	github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jinzhu/copier v0.3.4 // indirect
	github.com/magefile/mage v1.14.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	gopkg.in/djherbis/times.v1 v1.2.0 // indirect
)
