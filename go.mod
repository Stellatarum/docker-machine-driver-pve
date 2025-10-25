module github.com/stellatarum/docker-machine-driver-pve

go 1.24.1

replace github.com/docker/docker => github.com/moby/moby v1.4.2-0.20170731201646-1009e6a40b29

require (
	github.com/luthermonson/go-proxmox v0.2.3
	github.com/rancher/machine v0.15.0-rancher134
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.40.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/buger/goterm v1.0.4 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/diskfs/go-diskfs v1.5.0 // indirect
	github.com/djherbis/times v1.6.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jinzhu/copier v0.3.4 // indirect
	github.com/magefile/mage v1.14.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/term v0.33.0 // indirect
)
