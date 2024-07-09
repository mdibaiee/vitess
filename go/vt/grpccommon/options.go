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

package grpccommon

import (
	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	"github.com/mdibaiee/vitess/go/stats"
)

var (
	// maxMessageSize is the maximum message size which the gRPC server will
	// accept. Larger messages will be rejected.
	// Note: We're using 16 MiB as default value because that's the default in MySQL
	maxMessageSize = 16 * 1024 * 1024
	// enablePrometheus sets a flag to enable grpc client/server grpc monitoring.
	enablePrometheus bool
)

// RegisterFlags installs grpccommon flags on the given FlagSet.
//
// `go/cmd/*` entrypoints should either use servenv.ParseFlags(WithArgs)? which
// calls this function, or call this function directly before parsing
// command-line arguments.
func RegisterFlags(fs *pflag.FlagSet) {
	fs.IntVar(&maxMessageSize, "grpc_max_message_size", maxMessageSize, "Maximum allowed RPC message size. Larger messages will be rejected by gRPC with the error 'exceeding the max size'.")
	fs.BoolVar(&grpc.EnableTracing, "grpc_enable_tracing", grpc.EnableTracing, "Enable gRPC tracing.")
	fs.BoolVar(&enablePrometheus, "grpc_prometheus", enablePrometheus, "Enable gRPC monitoring with Prometheus.")
}

// EnableGRPCPrometheus returns the value of the --grpc_prometheus flag.
func EnableGRPCPrometheus() bool {
	return enablePrometheus
}

// MaxMessageSize returns the value of the --grpc_max_message_size flag.
func MaxMessageSize() int {
	return maxMessageSize
}

func init() {
	stats.NewString("GrpcVersion").Set(grpc.Version)
}
