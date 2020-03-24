package trinity

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFreePorts(t *testing.T) {
	port, err := GetFreePort()
	assert.Equal(t, nil, err, "func exec error , test failed")
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%v", port))
	_, err = net.ListenTCP("tcp", addr)
	assert.Equal(t, nil, err, "func exec error , test failed")
}
