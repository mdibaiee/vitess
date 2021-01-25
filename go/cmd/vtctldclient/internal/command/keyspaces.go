/*
Copyright 2021 The Vitess Authors.

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

package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"vitess.io/vitess/go/cmd/vtctldclient/cli"
	"vitess.io/vitess/go/vt/logutil"
	"vitess.io/vitess/go/vt/topo"

	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
	vtctldatapb "vitess.io/vitess/go/vt/proto/vtctldata"
	"vitess.io/vitess/go/vt/proto/vttime"
)

var (
	// CreateKeyspace makes a CreateKeyspace gRPC call to a vtctld.
	CreateKeyspace = &cobra.Command{
		Use:  "CreateKeyspace KEYSPACE_NAME [--force] [--sharding-column-name NAME --sharding-column-type TYPE] [--base-keyspace KEYSPACE --snapshot-timestamp TIME] [--served-from DB_TYPE:KEYSPACE ...]",
		Args: cobra.ExactArgs(1),
		RunE: commandCreateKeyspace,
	}
	// FindAllShardsInKeyspace makes a FindAllShardsInKeyspace gRPC call to a vtctld.
	FindAllShardsInKeyspace = &cobra.Command{
		Use:     "FindAllShardsInKeyspace keyspace",
		Aliases: []string{"findallshardsinkeyspace"},
		Args:    cobra.ExactArgs(1),
		RunE:    commandFindAllShardsInKeyspace,
	}
	// GetKeyspace makes a GetKeyspace gRPC call to a vtctld.
	GetKeyspace = &cobra.Command{
		Use:     "GetKeyspace keyspace",
		Aliases: []string{"getkeyspace"},
		Args:    cobra.ExactArgs(1),
		RunE:    commandGetKeyspace,
	}
	// GetKeyspaces makes a GetKeyspaces gRPC call to a vtctld.
	GetKeyspaces = &cobra.Command{
		Use:     "GetKeyspaces",
		Aliases: []string{"getkeyspaces"},
		Args:    cobra.NoArgs,
		RunE:    commandGetKeyspaces,
	}
)

var createKeyspaceOptions = struct {
	Force             bool
	AllowEmptyVSchema bool

	ShardingColumnName string
	ShardingColumnType cli.KeyspaceIDTypeFlag

	ServedFromsMap cli.StringMapValue

	KeyspaceType      cli.KeyspaceTypeFlag
	BaseKeyspace      string
	SnapshotTimestamp string
}{
	KeyspaceType: cli.KeyspaceTypeFlag(topodatapb.KeyspaceType_NORMAL),
}

func commandCreateKeyspace(cmd *cobra.Command, args []string) error {
	name := cmd.Flags().Arg(0)

	switch topodatapb.KeyspaceType(createKeyspaceOptions.KeyspaceType) {
	case topodatapb.KeyspaceType_NORMAL, topodatapb.KeyspaceType_SNAPSHOT:
	default:
		return errors.New("invalid keyspace type passed to --type")
	}

	var snapshotTime *vttime.Time
	if topodatapb.KeyspaceType(createKeyspaceOptions.KeyspaceType) == topodatapb.KeyspaceType_SNAPSHOT {
		if createKeyspaceOptions.BaseKeyspace == "" {
			return errors.New("--base-keyspace is required for a snapshot keyspace")
		}

		if createKeyspaceOptions.SnapshotTimestamp == "" {
			return errors.New("--snapshot-timestamp is required for a snapshot keyspace")
		}

		t, err := time.Parse(time.RFC3339, createKeyspaceOptions.SnapshotTimestamp)
		if err != nil {
			return err
		}

		if now := time.Now(); t.After(now) {
			return fmt.Errorf("--snapshot-time cannot be in the future; snapshot = %v, now = %v", t, now)
		}

		snapshotTime = logutil.TimeToProto(t)
	}

	req := &vtctldatapb.CreateKeyspaceRequest{
		Name:               name,
		Force:              createKeyspaceOptions.Force,
		AllowEmptyVSchema:  createKeyspaceOptions.AllowEmptyVSchema,
		ShardingColumnName: createKeyspaceOptions.ShardingColumnName,
		ShardingColumnType: topodatapb.KeyspaceIdType(createKeyspaceOptions.ShardingColumnType),
		Type:               topodatapb.KeyspaceType(createKeyspaceOptions.KeyspaceType),
		BaseKeyspace:       createKeyspaceOptions.BaseKeyspace,
		SnapshotTime:       snapshotTime,
	}

	for n, v := range createKeyspaceOptions.ServedFromsMap.StringMapValue {
		tt, err := topo.ParseServingTabletType(n)
		if err != nil {
			return err
		}

		req.ServedFroms = append(req.ServedFroms, &topodatapb.Keyspace_ServedFrom{
			TabletType: tt,
			Keyspace:   v,
		})
	}

	resp, err := client.CreateKeyspace(commandCtx, req)
	if err != nil {
		return err
	}

	data, err := cli.MarshalJSON(resp.Keyspace)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully created keyspace %s. Result:\n%s\n", name, data)

	return nil
}

func commandFindAllShardsInKeyspace(cmd *cobra.Command, args []string) error {
	ks := cmd.Flags().Arg(0)
	resp, err := client.FindAllShardsInKeyspace(commandCtx, &vtctldatapb.FindAllShardsInKeyspaceRequest{
		Keyspace: ks,
	})

	if err != nil {
		return err
	}

	data, err := cli.MarshalJSON(resp)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", data)
	return nil
}

func commandGetKeyspace(cmd *cobra.Command, args []string) error {
	ks := cmd.Flags().Arg(0)
	resp, err := client.GetKeyspace(commandCtx, &vtctldatapb.GetKeyspaceRequest{
		Keyspace: ks,
	})

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp.Keyspace)

	return nil
}

func commandGetKeyspaces(cmd *cobra.Command, args []string) error {
	resp, err := client.GetKeyspaces(commandCtx, &vtctldatapb.GetKeyspacesRequest{})
	if err != nil {
		return err
	}

	data, err := cli.MarshalJSON(resp.Keyspaces)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", data)

	return nil
}

func init() {
	CreateKeyspace.Flags().BoolVarP(&createKeyspaceOptions.Force, "force", "f", false, "Proceeds even if the keyspace already exists. Does not overwrite the existing keyspace record")
	CreateKeyspace.Flags().BoolVarP(&createKeyspaceOptions.AllowEmptyVSchema, "allow-empty-vschema", "e", false, "Allows a new keyspace to have no vschema")
	CreateKeyspace.Flags().StringVar(&createKeyspaceOptions.ShardingColumnName, "sharding-column-name", "", "The column name to use for sharding operations")
	CreateKeyspace.Flags().Var(&createKeyspaceOptions.ShardingColumnType, "sharding-column-type", "The type of the column to use for sharding operations")
	CreateKeyspace.Flags().Var(&createKeyspaceOptions.ServedFromsMap, "served-from", "TODO")
	CreateKeyspace.Flags().Var(&createKeyspaceOptions.KeyspaceType, "type", "The type of the keyspace")
	CreateKeyspace.Flags().StringVar(&createKeyspaceOptions.BaseKeyspace, "base-keyspace", "", "The base keyspace for a snapshot keyspace.")
	CreateKeyspace.Flags().StringVar(&createKeyspaceOptions.SnapshotTimestamp, "snapshot-timestamp", "", "The snapshot time for a snapshot keyspace, as a timestamp in RFC3339 format.")
	Root.AddCommand(CreateKeyspace)

	Root.AddCommand(FindAllShardsInKeyspace)
	Root.AddCommand(GetKeyspace)
	Root.AddCommand(GetKeyspaces)
}
