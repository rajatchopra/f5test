package f5

import (
	"fmt"
	"github.com/golang/glog"
)

// F5Plugin holds state for the f5 plugin.
type F5Plugin struct {
	// F5Client is the object that represents the F5 BIG-IP host, holds state,
	// and provides an interface to manipulate F5 BIG-IP.
	F5Client *f5LTM
}

// F5PluginConfig holds configuration for the f5 plugin.
type F5PluginConfig struct {
	// Host specifies the hostname or IP address of the F5 BIG-IP host.
	Host string

	// Username specifies the username with the plugin should authenticate
	// with the F5 BIG-IP host.
	Username string

	// Password specifies the password with which the plugin should
	// authenticate with F5 BIG-IP.
	Password string

	// HttpVserver specifies the name of the vserver object in F5 BIG-IP that the
	// plugin will configure for HTTP connections.
	HttpVserver string

	// HttpsVserver specifies the name of the vserver object in F5 BIG-IP that the
	// plugin will configure for HTTPS connections.
	HttpsVserver string

	// PrivateKey specifies the path to the SSH private-key file for
	// authenticating with F5 BIG-IP.  The file must exist with this pathname
	// inside the F5 router's filesystem namespace.  The F5 router requires this
	// key to copy certificates and keys to the F5 BIG-IP host.
	PrivateKey string

	// Insecure specifies whether the F5 plugin should perform strict certificate
	// validation for connections to the F5 BIG-IP host.
	Insecure bool

	// PartitionPath specifies the F5 partition path to use. This is used
	// to create an access control boundary for users and applications.
	PartitionPath string
}

// NewF5Plugin makes a new f5 router plugin.
func NewF5Plugin(cfg F5PluginConfig) (*F5Plugin, error) {
	f5LTMCfg := f5LTMCfg{
		host:          cfg.Host,
		username:      cfg.Username,
		password:      cfg.Password,
		httpVserver:   cfg.HttpVserver,
		httpsVserver:  cfg.HttpsVserver,
		privkey:       cfg.PrivateKey,
		insecure:      cfg.Insecure,
		partitionPath: cfg.PartitionPath,
	}
	f5, err := newF5LTM(f5LTMCfg)
	if err != nil {
		return nil, err
	}
	return &F5Plugin{f5}, f5.Initialize()
}

// ensurePoolExists checks whether the named pool already exists in F5 BIG-IP
// and creates it if it does not.
func (p *F5Plugin) ensurePoolExists(poolname string) error {
	poolExists, err := p.F5Client.PoolExists(poolname)
	if err != nil {
		glog.V(4).Infof("F5Client.PoolExists failed: %v", err)
		return err
	}

	if !poolExists {
		err = p.F5Client.CreatePool(poolname)
		if err != nil {
			glog.V(4).Infof("Error creating pool %s: %v", poolname, err)
			return err
		}
	}

	return nil
}

// deletePool delete the named pool from F5 BIG-IP.
func (p *F5Plugin) deletePool(poolname string) error {
	poolExists, err := p.F5Client.PoolExists(poolname)
	if err != nil {
		glog.V(4).Infof("F5Client.PoolExists failed: %v", err)
		return err
	}

	if poolExists {
		err = p.F5Client.DeletePool(poolname)
		if err != nil {
			glog.V(4).Infof("Error deleting pool %s: %v", poolname, err)
			return err
		}
	}

	return nil
}

// deletePoolIfEmpty deletes the named pool from F5 BIG-IP if, and only if, it
// has no members.
func (p *F5Plugin) deletePoolIfEmpty(poolname string) error {
	poolExists, err := p.F5Client.PoolExists(poolname)
	if err != nil {
		glog.V(4).Infof("F5Client.PoolExists failed: %v", err)
		return err
	}

	if poolExists {
		members, err := p.F5Client.GetPoolMembers(poolname)
		if err != nil {
			glog.V(4).Infof("F5Client.GetPoolMembers failed: %v", err)
			return err
		}

		// We only delete the pool if the pool is empty, which it may not be
		// if a service has been added and has not (yet) been deleted.
		if len(members) == 0 {
			err = p.F5Client.DeletePool(poolname)
			if err != nil {
				glog.V(4).Infof("Error deleting pool %s: %v", poolname, err)
				return err
			}
		}
	}

	return nil
}

// poolName returns a string that can be used as a poolname in F5 BIG-IP and
// is distinct for the given endpoints namespace and name.
func poolName(endpointsNamespace, endpointsName string) string {
	return fmt.Sprintf("openshift_%s_%s", endpointsNamespace, endpointsName)
}


// No-op since f5 configuration can be updated piecemeal
func (p *F5Plugin) SetLastSyncProcessed(processed bool) error {
	return nil
}
