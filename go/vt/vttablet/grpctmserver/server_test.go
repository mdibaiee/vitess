/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package grpctmserver_test

import (
	"net"
	"testing"

	"google.golang.org/grpc"

	"github.com/mdibaiee/vitess/go/vt/vttablet/grpctmclient"
	"github.com/mdibaiee/vitess/go/vt/vttablet/grpctmserver"
	"github.com/mdibaiee/vitess/go/vt/vttablet/tmrpctest"

	topodatapb "github.com/mdibaiee/vitess/go/vt/proto/topodata"
)

// TestGRPCTMServer creates a fake server implementation, a fake client
// implementation, and runs the test suite against the setup.
func TestGRPCTMServer(t *testing.T) {
	// Listen on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Cannot listen: %v", err)
	}
	host := listener.Addr().(*net.TCPAddr).IP.String()
	port := int32(listener.Addr().(*net.TCPAddr).Port)

	// Create a gRPC server and listen on the port.
	s := grpc.NewServer()
	fakeTM := tmrpctest.NewFakeRPCTM(t)
	grpctmserver.RegisterForTest(s, fakeTM)
	go s.Serve(listener)

	// Create a gRPC client to talk to the fake tablet.
	client := grpctmclient.NewClient()
	tablet := &topodatapb.Tablet{
		Alias: &topodatapb.TabletAlias{
			Cell: "test",
			Uid:  123,
		},
		Hostname: host,
		PortMap: map[string]int32{
			"grpc": port,
		},
	}

	// and run the test suite
	tmrpctest.Run(t, client, tablet, fakeTM)
}
