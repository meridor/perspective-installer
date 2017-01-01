//go:generate stringer -type=Cloud enums.go
package main

// Supported clouds
type Cloud int

const (
	DIGITAL_OCEAN Cloud = iota
	DOCKER
	OPENSTACK
)
