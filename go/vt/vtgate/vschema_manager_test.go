package vtgate

import (
	"testing"

	"github.com/stretchr/testify/require"

	"vitess.io/vitess/go/test/utils"
	querypb "vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/sqlparser"

	vschemapb "vitess.io/vitess/go/vt/proto/vschema"
	"vitess.io/vitess/go/vt/vtgate/vindexes"
)

func TestWatchSrvVSchema(t *testing.T) {
	cols := []vindexes.Column{{
		Name: sqlparser.NewColIdent("id"),
		Type: querypb.Type_INT64,
	}}
	cols2 := []vindexes.Column{{
		Name: sqlparser.NewColIdent("id"),
		Type: querypb.Type_INT64,
	}, {
		Name: sqlparser.NewColIdent("name"),
		Type: querypb.Type_VARCHAR,
	}}
	ks := &vindexes.Keyspace{Name: "ks"}
	dual := &vindexes.Table{Type: vindexes.TypeReference, Name: sqlparser.NewTableIdent("dual"), Keyspace: ks}
	tcases := []struct {
		name       string
		srvVschema *vschemapb.SrvVSchema
		st         []SchemaTable
		expected   map[string]*vindexes.Table
	}{{
		name:       "Single table known by mysql schema and not by vschema",
		srvVschema: &vschemapb.SrvVSchema{Keyspaces: map[string]*vschemapb.Keyspace{"ks": {}}},
		st:         []SchemaTable{{Name: "tbl", Columns: cols}},
		expected: map[string]*vindexes.Table{
			"dual": dual,
			"tbl": {
				Name:                    sqlparser.NewTableIdent("tbl"),
				Keyspace:                ks,
				Columns:                 cols,
				ColumnListAuthoritative: true,
			},
		},
	}, {
		name: "Single table known by both - vschema is not authoritative",
		srvVschema: &vschemapb.SrvVSchema{Keyspaces: map[string]*vschemapb.Keyspace{"ks": {
			Tables: map[string]*vschemapb.Table{
				"tbl": {}, // we know of it, but nothing else
			},
		}}},
		st: []SchemaTable{{Name: "tbl", Columns: cols}},
		expected: map[string]*vindexes.Table{
			"dual": dual,
			"tbl": {
				Name:                    sqlparser.NewTableIdent("tbl"),
				Keyspace:                ks,
				Columns:                 cols,
				ColumnListAuthoritative: true,
			},
		},
	}, {
		name: "Single table known by both - vschema is authoritative",
		srvVschema: &vschemapb.SrvVSchema{Keyspaces: map[string]*vschemapb.Keyspace{"ks": {
			Tables: map[string]*vschemapb.Table{
				"tbl": {
					Columns: []*vschemapb.Column{
						{Name: "id", Type: querypb.Type_INT64},
						{Name: "name", Type: querypb.Type_VARCHAR},
					},
					ColumnListAuthoritative: true},
			},
		}}},
		st: []SchemaTable{{Name: "tbl", Columns: cols}},
		expected: map[string]*vindexes.Table{
			"dual": dual,
			"tbl": {
				Name:                    sqlparser.NewTableIdent("tbl"),
				Keyspace:                ks,
				Columns:                 cols2,
				ColumnListAuthoritative: true,
			},
		},
	}}

	vm := &VSchemaManager{}
	var vs *vindexes.VSchema
	vm.subscriber = func(vschema *vindexes.VSchema, _ *VSchemaStats) {
		vs = vschema
	}
	for _, tcase := range tcases {
		t.Run(tcase.name, func(t *testing.T) {
			vs = nil
			vm.schema = fakeSchema{t: tcase.st}
			vm.VSchemaUpdate(tcase.srvVschema, nil)

			require.NotNil(t, vs)
			ks := vs.Keyspaces["ks"]
			require.NotNil(t, ks, "keyspace was not found")
			utils.MustMatch(t, tcase.expected, ks.Tables)
		})
	}
}

type fakeSchema struct {
	t []SchemaTable
}

func (f fakeSchema) Tables(_ string) []SchemaTable {
	return f.t
}
