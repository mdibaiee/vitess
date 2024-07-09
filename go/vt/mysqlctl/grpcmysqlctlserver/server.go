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
Package grpcmysqlctlserver contains the gRPC implementation of the server
side of the remote execution of mysqlctl commands.
*/
package grpcmysqlctlserver

import (
	"context"

	"google.golang.org/grpc"

	"github.com/mdibaiee/vitess/go/protoutil"
	"github.com/mdibaiee/vitess/go/vt/mysqlctl"
	mysqlctlpb "github.com/mdibaiee/vitess/go/vt/proto/mysqlctl"
)

// server is our gRPC server.
type server struct {
	mysqlctlpb.UnimplementedMysqlCtlServer
	cnf    *mysqlctl.Mycnf
	mysqld *mysqlctl.Mysqld
}

// Start implements the server side of the MysqlctlClient interface.
func (s *server) Start(ctx context.Context, request *mysqlctlpb.StartRequest) (*mysqlctlpb.StartResponse, error) {
	return &mysqlctlpb.StartResponse{}, s.mysqld.Start(ctx, s.cnf, request.MysqldArgs...)
}

// Shutdown implements the server side of the MysqlctlClient interface.
func (s *server) Shutdown(ctx context.Context, request *mysqlctlpb.ShutdownRequest) (*mysqlctlpb.ShutdownResponse, error) {
	timeout, ok, err := protoutil.DurationFromProto(request.MysqlShutdownTimeout)
	if err != nil {
		return nil, err
	}
	if !ok {
		timeout = mysqlctl.DefaultShutdownTimeout
	}
	return &mysqlctlpb.ShutdownResponse{}, s.mysqld.Shutdown(ctx, s.cnf, request.WaitForMysqld, timeout)
}

// RunMysqlUpgrade implements the server side of the MysqlctlClient interface.
func (s *server) RunMysqlUpgrade(ctx context.Context, _ *mysqlctlpb.RunMysqlUpgradeRequest) (*mysqlctlpb.RunMysqlUpgradeResponse, error) {
	return &mysqlctlpb.RunMysqlUpgradeResponse{}, s.mysqld.RunMysqlUpgrade(ctx)
}

// RunMysqlUpgrade implements the server side of the MysqlctlClient interface.
func (s *server) ApplyBinlogFile(ctx context.Context, request *mysqlctlpb.ApplyBinlogFileRequest) (*mysqlctlpb.ApplyBinlogFileResponse, error) {
	return &mysqlctlpb.ApplyBinlogFileResponse{}, s.mysqld.ApplyBinlogFile(ctx, request)
}

// ReadBinlogFilesTimestamps implements the server side of the MysqlctlClient interface.
func (s *server) ReadBinlogFilesTimestamps(ctx context.Context, request *mysqlctlpb.ReadBinlogFilesTimestampsRequest) (*mysqlctlpb.ReadBinlogFilesTimestampsResponse, error) {
	return s.mysqld.ReadBinlogFilesTimestamps(ctx, request)
}

// ReinitConfig implements the server side of the MysqlctlClient interface.
func (s *server) ReinitConfig(ctx context.Context, request *mysqlctlpb.ReinitConfigRequest) (*mysqlctlpb.ReinitConfigResponse, error) {
	return &mysqlctlpb.ReinitConfigResponse{}, s.mysqld.ReinitConfig(ctx, s.cnf)
}

// RefreshConfig implements the server side of the MysqlctlClient interface.
func (s *server) RefreshConfig(ctx context.Context, request *mysqlctlpb.RefreshConfigRequest) (*mysqlctlpb.RefreshConfigResponse, error) {
	return &mysqlctlpb.RefreshConfigResponse{}, s.mysqld.RefreshConfig(ctx, s.cnf)
}

// VersionString registers the Server for RPCs.
func (s *server) VersionString(ctx context.Context, request *mysqlctlpb.VersionStringRequest) (*mysqlctlpb.VersionStringResponse, error) {
	version, err := s.mysqld.GetVersionString(ctx)
	if err != nil {
		return nil, err
	}
	return &mysqlctlpb.VersionStringResponse{Version: version}, nil
}

// StartServer registers the Server for RPCs.
func StartServer(s *grpc.Server, cnf *mysqlctl.Mycnf, mysqld *mysqlctl.Mysqld) {
	mysqlctlpb.RegisterMysqlCtlServer(s, &server{cnf: cnf, mysqld: mysqld})
}
