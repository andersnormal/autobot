package bolt

import (
	"github.com/andersnormal/autobot/pkg/provider"
	"github.com/andersnormal/autobot/pkg/provider/kv"

	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
)

var _ provider.Provider = (*Provider)(nil)

// Provider holds configuration for BoltDB.
type Provider struct {
	kv.Provider
}

// CreateStore creates the etcd store
func (p *Provider) CreateStore(bucket string) (store.Store, error) {
	p.SetStoreType(store.ETCD)
	etcd.Register()

	return p.Provider.CreateStore(bucket)
}
