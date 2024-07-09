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

/*
Package grpcvtctlserver contains the gRPC implementation of the server side
of the remote execution of vtctl commands.
*/
package grpcvtctlserver

import (
	"sync"

	"google.golang.org/grpc"

	"github.com/mdibaiee/vitess/go/vt/vtenv"

	"github.com/mdibaiee/vitess/go/vt/logutil"
	"github.com/mdibaiee/vitess/go/vt/servenv"
	"github.com/mdibaiee/vitess/go/vt/topo"
	"github.com/mdibaiee/vitess/go/vt/vtctl"
	"github.com/mdibaiee/vitess/go/vt/vttablet/tmclient"
	"github.com/mdibaiee/vitess/go/vt/wrangler"

	logutilpb "github.com/mdibaiee/vitess/go/vt/proto/logutil"
	vtctldatapb "github.com/mdibaiee/vitess/go/vt/proto/vtctldata"
	vtctlservicepb "github.com/mdibaiee/vitess/go/vt/proto/vtctlservice"
)

// VtctlServer is our RPC server
type VtctlServer struct {
	vtctlservicepb.UnimplementedVtctlServer
	ts  *topo.Server
	env *vtenv.Environment
}

// NewVtctlServer returns a new Vtctl Server for the topo server.
func NewVtctlServer(env *vtenv.Environment, ts *topo.Server) *VtctlServer {
	return &VtctlServer{env: env, ts: ts}
}

// ExecuteVtctlCommand is part of the vtctldatapb.VtctlServer interface
func (s *VtctlServer) ExecuteVtctlCommand(args *vtctldatapb.ExecuteVtctlCommandRequest, stream vtctlservicepb.Vtctl_ExecuteVtctlCommandServer) (err error) {
	defer servenv.HandlePanic("vtctl", &err)

	// Create a logger, send the result back to the caller.
	// We may execute this in parallel (inside multiple go routines),
	// but the stream.Send() method is not thread safe in gRPC.
	// So use a mutex to protect it.
	mu := sync.Mutex{}
	logstream := logutil.NewCallbackLogger(func(e *logutilpb.Event) {
		// If the client disconnects, we will just fail
		// to send the log events, but won't interrupt
		// the command.
		mu.Lock()
		stream.Send(&vtctldatapb.ExecuteVtctlCommandResponse{
			Event: e,
		})
		mu.Unlock()
	})
	logger := logutil.NewTeeLogger(logstream, logutil.NewConsoleLogger())

	// create the wrangler
	tmc := tmclient.NewTabletManagerClient()
	defer tmc.Close()
	wr := wrangler.New(s.env, logger, s.ts, tmc)

	// execute the command
	return vtctl.RunCommand(stream.Context(), wr, args.Args)
}

// StartServer registers the VtctlServer for RPCs
func StartServer(s *grpc.Server, env *vtenv.Environment, ts *topo.Server) {
	vtctlservicepb.RegisterVtctlServer(s, NewVtctlServer(env, ts))
}
