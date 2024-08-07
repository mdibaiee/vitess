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

package grpcbinlogplayer

import (
	"context"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	"github.com/mdibaiee/vitess/go/netutil"
	"github.com/mdibaiee/vitess/go/vt/binlog/binlogplayer"
	"github.com/mdibaiee/vitess/go/vt/grpcclient"
	binlogdatapb "github.com/mdibaiee/vitess/go/vt/proto/binlogdata"
	binlogservicepb "github.com/mdibaiee/vitess/go/vt/proto/binlogservice"
	topodatapb "github.com/mdibaiee/vitess/go/vt/proto/topodata"
	"github.com/mdibaiee/vitess/go/vt/servenv"
)

var cert, key, ca, crl, name string

func init() {
	servenv.OnParseFor("vtcombo", registerFlags)
	servenv.OnParseFor("vttablet", registerFlags)
}

func registerFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cert, "binlog_player_grpc_cert", cert, "the cert to use to connect")
	fs.StringVar(&key, "binlog_player_grpc_key", key, "the key to use to connect")
	fs.StringVar(&ca, "binlog_player_grpc_ca", ca, "the server ca to use to validate servers when connecting")
	fs.StringVar(&crl, "binlog_player_grpc_crl", crl, "the server crl to use to validate server certificates when connecting")
	fs.StringVar(&name, "binlog_player_grpc_server_name", name, "the server name to use to validate server certificate")
}

// client implements a Client over go rpc
type client struct {
	cc *grpc.ClientConn
	c  binlogservicepb.UpdateStreamClient
}

func (client *client) Dial(ctx context.Context, tablet *topodatapb.Tablet) error {
	addr := netutil.JoinHostPort(tablet.Hostname, tablet.PortMap["grpc"])
	var err error
	opt, err := grpcclient.SecureDialOption(cert, key, ca, crl, name)
	if err != nil {
		return err
	}
	client.cc, err = grpcclient.DialContext(ctx, addr, grpcclient.FailFast(true), opt)
	if err != nil {
		return err
	}
	client.c = binlogservicepb.NewUpdateStreamClient(client.cc)
	return nil
}

func (client *client) Close() {
	client.cc.Close()
}

type serveStreamKeyRangeAdapter struct {
	stream binlogservicepb.UpdateStream_StreamKeyRangeClient
}

func (s *serveStreamKeyRangeAdapter) Recv() (*binlogdatapb.BinlogTransaction, error) {
	r, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	return r.BinlogTransaction, nil
}

func (client *client) StreamKeyRange(ctx context.Context, position string, keyRange *topodatapb.KeyRange, charset *binlogdatapb.Charset) (binlogplayer.BinlogTransactionStream, error) {
	query := &binlogdatapb.StreamKeyRangeRequest{
		Position: position,
		KeyRange: keyRange,
		Charset:  charset,
	}
	stream, err := client.c.StreamKeyRange(ctx, query)
	if err != nil {
		return nil, err
	}
	return &serveStreamKeyRangeAdapter{stream}, nil
}

type serveStreamTablesAdapter struct {
	stream binlogservicepb.UpdateStream_StreamTablesClient
}

func (s *serveStreamTablesAdapter) Recv() (*binlogdatapb.BinlogTransaction, error) {
	r, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	return r.BinlogTransaction, nil
}

func (client *client) StreamTables(ctx context.Context, position string, tables []string, charset *binlogdatapb.Charset) (binlogplayer.BinlogTransactionStream, error) {
	query := &binlogdatapb.StreamTablesRequest{
		Position: position,
		Tables:   tables,
		Charset:  charset,
	}
	stream, err := client.c.StreamTables(ctx, query)
	if err != nil {
		return nil, err
	}
	return &serveStreamTablesAdapter{stream}, nil
}

// Registration as a factory
func init() {
	binlogplayer.RegisterClientFactory("grpc", func() binlogplayer.Client {
		return &client{}
	})
}
