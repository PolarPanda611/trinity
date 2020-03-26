package trinity

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

const (
	consulScheme = "consul"
)

type consulResolverBuilder struct {
	address     string
	client      *consulapi.Client
	serviceName string
}

func (crb *consulResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	consulServiceName := crb.serviceName
	addrs, err := crb.resolveServiceFromConsul()
	if err != nil {
		return &consulResolver{}, nil
	}
	ch := make(chan struct{}, 1)
	r := &consulResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]resolver.Address{
			consulServiceName: addrs,
		},
		rnCh: ch,
	}
	r.start()
	go crb.csMonitor(r)

	return r, nil
}
func (*consulResolverBuilder) Scheme() string { return consulScheme }

func (crb *consulResolverBuilder) resolveServiceFromConsul() ([]resolver.Address, error) {
	serviceEntries, _, err := crb.client.Health().Service(crb.serviceName, "", true, &consulapi.QueryOptions{})
	if err != nil {
		fmt.Println("call consul Health API failed, ", err)
		return nil, err
	}

	addrs := make([]resolver.Address, 0)
	for _, serviceEntry := range serviceEntries {
		address := resolver.Address{Addr: fmt.Sprintf("%s:%d", serviceEntry.Service.Address, serviceEntry.Service.Port)}
		addrs = append(addrs, address)
	}
	return addrs, nil
}

func (crb *consulResolverBuilder) csMonitor(cr *consulResolver) {
	t := time.NewTicker(500 * time.Millisecond)
	//Get service addresses from consul every 500 Millisecond and update them to gRPC
	for {
		select {
		case <-t.C:
		//resolve now
		case <-cr.rnCh:
			//fmt.Println("resolve service adress now!")
		}

		addrs, err := crb.resolveServiceFromConsul()
		if err != nil {
			fmt.Println("resolveServiceFromConsul failed: ", err.Error())
			continue
		} else {
			//fmt.Println("resolveServiceFromConsul success: ", addrs)
		}
		cr.cc.UpdateState(resolver.State{Addresses: addrs})
	}
}

type consulResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]resolver.Address
	rnCh       chan struct{}
}

func (r *consulResolver) start() {
	addrs := r.addrsStore[r.target.Endpoint]
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

// ResolveNow gRPC will call it when connection is broken
func (r *consulResolver) ResolveNow(o resolver.ResolveNowOption) {
	r.rnCh <- struct{}{}
}
func (*consulResolver) Close() {}

// ConsulResolverInit init coonsul
func ConsulResolverInit(address string, serviceName string) error {
	config := consulapi.DefaultConfig()
	config.Address = address
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("new consul client failed: ", err.Error())
		return err
	}
	crb := &consulResolverBuilder{address: address, client: client, serviceName: serviceName}
	resolver.Register(crb)

	return nil
}

// NewEtcdClientConn new etcd client connection
func NewEtcdClientConn(address string, port int, serviceName string) (*grpc.ClientConn, error) {
	cli, err := clientv3.NewFromURL(fmt.Sprintf("http://%v:%v", address, port))
	if err != nil {
		return nil, fmt.Errorf("failed to conn etcd client , %v", err)
	}
	r := &etcdnaming.GRPCResolver{Client: cli}
	b := grpc.RoundRobin(r)

	ctx1, cel := context.WithTimeout(context.Background(), time.Second*3)
	defer cel()
	conn, err := grpc.DialContext(ctx1, serviceName, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBalancer(b))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
