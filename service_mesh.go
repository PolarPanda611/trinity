package trinity

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

// ServiceMesh interface
type ServiceMesh interface {
	RegService() error
	DeRegService() error
}

// ServiceMeshConsulImpl consul register
type ServiceMeshConsulImpl struct {
	// config
	ConsulAddress                  string // consul address
	ConsulPort                     int
	ProjectName                    string
	ProjectVersion                 string
	ServiceIP                      string
	Tags                           []string // consul tags
	ServicePort                    int      //service port
	DeregisterCriticalServiceAfter time.Duration
	ServiceCheckInterval           time.Duration

	// runtime
	consulClient *consulapi.Client
}

// NewConsulRegister New consul register
func NewConsulRegister(consulAddress string, consulPort int, projectName string, projectVersion string, serviceIP string, servicePort int, Tags []string, deregisterSecondAfterCritical int, interval int) (ServiceMesh, error) {
	c := &ServiceMeshConsulImpl{
		ConsulAddress:                  consulAddress, //localhost:8500
		ConsulPort:                     consulPort,
		ProjectName:                    projectName,
		ProjectVersion:                 projectVersion,
		ServiceIP:                      serviceIP,
		ServicePort:                    servicePort,
		Tags:                           Tags,
		DeregisterCriticalServiceAfter: time.Duration(deregisterSecondAfterCritical) * time.Second,
		ServiceCheckInterval:           time.Duration(interval) * time.Second,
	}
	if err := c.getConsulClient(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ServiceMeshConsulImpl) getConsulClient() error {
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
		return err
	}
	c.consulClient = client
	return nil
}

// RegService  consul register service
func (c *ServiceMeshConsulImpl) RegService() error {
	reg := consulapi.AgentServiceRegistration{
		ID:      GetServiceID(c.ProjectName, c.ProjectVersion, c.ServiceIP, c.ServicePort),
		Name:    GetServiceName(c.ProjectName, c.ProjectVersion),
		Tags:    c.Tags,
		Port:    c.ServicePort,
		Address: c.ServiceIP,
		// Check: &consulapi.AgentServiceCheck{
		// 	// 健康检查间隔
		// 	Interval: c.Interval.String(),
		// 	//grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
		// 	GRPC: fmt.Sprintf("%v:%v/%v", c.ServiceIP, c.ServicePort, c.ServiceName),
		// 	// 注销时间，相当于过期时间
		// 	DeregisterCriticalServiceAfter: c.DeregisterCriticalServiceAfter.String(),
		// 	GRPCUseTLS:                     true,
		// 	TLSSkipVerify:                  true,
		// },
	}
	if err := c.consulClient.Agent().ServiceRegister(&reg); err != nil {
		return err
	}
	fmt.Println("service register successfully ")
	return nil
}

// DeRegService  deregister service
func (c *ServiceMeshConsulImpl) DeRegService() error {
	if err := c.consulClient.Agent().ServiceDeregister(GetServiceID(c.ProjectName, c.ProjectVersion, c.ServiceIP, c.ServicePort)); err != nil {
		return err
	}
	fmt.Println("service deregister successfully ")
	return nil
}
