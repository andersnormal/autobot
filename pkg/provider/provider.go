package provider

// AbstractProvider is the base provider from which every provider inherits.
type AbstractProvider struct {
	Enable bool
}

// Provider defines the interface to a data provider (e.g. etcd, or bolt)
type Provider interface {
}
