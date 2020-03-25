package trinity

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/naming"
)

// ServiceMesh interface
type ServiceMesh interface {
	GetClient() interface{}
	RegService(projectName string, projectVersion string, serviceIP string, servicePort int, Tags []string, deregisterSecondAfterCritical int, interval int, tlsEnabled bool) error
	DeRegService(projectName string, projectVersion string, serviceIP string, servicePort int) error
}

// ServiceMeshConsulImpl consul register
type ServiceMeshConsulImpl struct {
	// config
	ConsulAddress string // consul address
	ConsulPort    int

	// runtime
	consulClient *consulapi.Client
}

// NewConsulRegister New consul register
func NewConsulRegister(consulAddress string, consulPort int) (ServiceMesh, error) {
	c := &ServiceMeshConsulImpl{
		ConsulAddress: consulAddress, //localhost:8500
		ConsulPort:    consulPort,
	}
	config := consulapi.DefaultConfig()
	config.Address = fmt.Sprintf("%v:%v", c.ConsulAddress, c.ConsulPort)
	// config.TLSConfig = consulapi.TLSConfig{
	// 	CAFile:             "/Users/daniel/Documents/workspace/SolutionDelivery/conf/ca.pem",
	// 	CAPath:             "/Users/daniel/Documents/workspace/SolutionDelivery/conf/ca.pem",
	// 	CertFile:           "/Users/daniel/Documents/workspace/SolutionDelivery/conf/client/client.pem",
	// 	KeyFile:            "/Users/daniel/Documents/workspace/SolutionDelivery/conf/client/client.key",
	// 	Address:            "SolutionDelivery",
	// 	InsecureSkipVerify: true,
	// }
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	c.consulClient = client
	return c, nil
}

// GetClient get consul client
func (c *ServiceMeshConsulImpl) GetClient() interface{} {
	return c.consulClient
}

// RegService  consul register service
func (c *ServiceMeshConsulImpl) RegService(projectName string, projectVersion string, serviceIP string, servicePort int, Tags []string, deregisterSecondAfterCritical int, interval int, tlsEnabled bool) error {
	reg := consulapi.AgentServiceRegistration{
		ID:      GetServiceID(projectName, projectVersion, serviceIP, servicePort),
		Name:    GetServiceName(projectName),
		Tags:    Tags,
		Port:    servicePort,
		Address: serviceIP,
	}
	// temperolly has heath check issue when using tls
	if !tlsEnabled {
		reg.Check = &consulapi.AgentServiceCheck{
			// 健康检查间隔
			Interval: (time.Duration(interval) * time.Second).String(),
			//grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			GRPC: fmt.Sprintf("%v:%v/%v", serviceIP, servicePort, GetServiceName(projectName)),
			// 注销时间，相当于过期时间
			DeregisterCriticalServiceAfter: (time.Duration(deregisterSecondAfterCritical) * time.Second).String(),
			GRPCUseTLS:                     tlsEnabled,
		}
	}
	if err := c.consulClient.Agent().ServiceRegister(&reg); err != nil {
		return err
	}
	fmt.Println("service register successfully ")
	return nil
}

// DeRegService  deregister service
func (c *ServiceMeshConsulImpl) DeRegService(projectName string, projectVersion string, serviceIP string, servicePort int) error {
	if err := c.consulClient.Agent().ServiceDeregister(GetServiceID(projectName, projectVersion, serviceIP, servicePort)); err != nil {
		return err
	}
	fmt.Println("service deregister successfully ")
	return nil
}

// ServiceMeshEtcdImpl consul register
type ServiceMeshEtcdImpl struct {
	// config
	Address string // consul address
	Port    int

	// runtime
	client *clientv3.Client
}

// NewEtcdRegister New consul register
func NewEtcdRegister(address string, port int) (ServiceMesh, error) {
	s := &ServiceMeshEtcdImpl{
		Address: address,
		Port:    port,
	}

	cli, err := clientv3.NewFromURL(fmt.Sprintf("http://%v:%v", s.Address, s.Port))

	if err != nil {
		return nil, err
	}
	s.client = cli
	return s, nil
}

func (s *ServiceMeshEtcdImpl) GetClient() interface{} {
	return s.client
}

func (s *ServiceMeshEtcdImpl) RegService(projectName string, projectVersion string, serviceIP string, servicePort int, Tags []string, deregisterSecondAfterCritical int, interval int, tlsEnabled bool) error {
	r := &etcdnaming.GRPCResolver{Client: s.client}
	err := r.Update(context.TODO(), GetServiceName(projectName), naming.Update{Op: naming.Add, Addr: fmt.Sprintf("%v:%v", serviceIP, servicePort), Metadata: fmt.Sprintf("%v", Tags)})
	if err != nil {
		return err
	}
	return nil
}
func (s *ServiceMeshEtcdImpl) DeRegService(projectName string, projectVersion string, serviceIP string, servicePort int) error {
	r := &etcdnaming.GRPCResolver{Client: s.client}
	err := r.Update(context.TODO(), GetServiceName(projectName), naming.Update{Op: naming.Delete, Addr: fmt.Sprintf("%v:%v", serviceIP, servicePort)})
	if err != nil {
		return err
	}
	return nil
}
