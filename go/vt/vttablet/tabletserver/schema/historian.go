/*
Copyright 2020 The Vitess Authors.

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

package schema

import (
	"context"
	"sort"
	"sync"
	"time"

	"mdibaiee/vitess/go/constants/sidecar"
	"mdibaiee/vitess/go/mysql/replication"
	"mdibaiee/vitess/go/sqltypes"
	"mdibaiee/vitess/go/vt/log"
	"mdibaiee/vitess/go/vt/sqlparser"
	"mdibaiee/vitess/go/vt/vttablet/tabletserver/connpool"
	"mdibaiee/vitess/go/vt/vttablet/tabletserver/tabletenv"

	binlogdatapb "mdibaiee/vitess/go/vt/proto/binlogdata"
)

const getInitialSchemaVersions = "select id, pos, ddl, time_updated, schemax from %s.schema_version where time_updated > %d order by id asc"
const getNextSchemaVersions = "select id, pos, ddl, time_updated, schemax from %s.schema_version where id > %d order by id asc"

// vl defines the glog verbosity level for the package
const vl = 10

// trackedSchema has the snapshot of the table at a given pos (reached by ddl)
type trackedSchema struct {
	schema      map[string]*binlogdatapb.MinimalTable
	pos         replication.Position
	ddl         string
	timeUpdated int64
}

// historian implements the Historian interface by calling schema.Engine for the underlying schema
// and supplying a schema for a specific version by loading the cached values from the schema_version table
// The schema version table is populated by the Tracker
type historian struct {
	conns               *connpool.Pool
	lastID              int64
	schemas             []*trackedSchema
	mu                  sync.Mutex
	enabled             bool
	isOpen              bool
	schemaMaxAgeSeconds int64
}

// newHistorian creates a new historian. It expects a schema.Engine instance
func newHistorian(enabled bool, schemaMaxAgeSeconds int64, conns *connpool.Pool) *historian {
	sh := historian{
		conns:               conns,
		lastID:              0,
		enabled:             enabled,
		schemaMaxAgeSeconds: schemaMaxAgeSeconds,
	}
	return &sh
}

func (h *historian) Enable(enabled bool) error {
	h.mu.Lock()
	h.enabled = enabled
	h.mu.Unlock()
	if enabled {
		return h.Open()
	}
	h.Close()
	return nil
}

// Open opens the underlying schema Engine. Called directly by a user purely interested in schema.Engine functionality
func (h *historian) Open() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.enabled {
		return nil
	}
	if h.isOpen {
		return nil
	}
	log.Info("Historian: opening")

	ctx := tabletenv.LocalContext()
	if err := h.loadFromDB(ctx); err != nil {
		log.Errorf("Historian failed to open: %v", err)
		return err
	}

	h.isOpen = true
	return nil
}

// Close closes the underlying schema engine and empties the version cache
func (h *historian) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.isOpen {
		return
	}

	h.schemas = nil
	h.isOpen = false
	log.Info("Historian: closed")
}

// RegisterVersionEvent is called by the vstream when it encounters a version event (an
// insert into the schema_tracking table). It triggers the historian to load the newer
// rows from the database to update its cache.
func (h *historian) RegisterVersionEvent() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.isOpen {
		return nil
	}
	ctx := tabletenv.LocalContext()
	if err := h.loadFromDB(ctx); err != nil {
		return err
	}
	return nil
}

// GetTableForPos returns a best-effort schema for a specific gtid
func (h *historian) GetTableForPos(tableName sqlparser.IdentifierCS, gtid string) (*binlogdatapb.MinimalTable, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.isOpen {
		return nil, nil
	}

	log.V(2).Infof("GetTableForPos called for %s with pos %s", tableName, gtid)
	if gtid == "" {
		return nil, nil
	}
	pos, err := replication.DecodePosition(gtid)
	if err != nil {
		return nil, err
	}
	var t *binlogdatapb.MinimalTable
	if len(h.schemas) > 0 {
		t = h.getTableFromHistoryForPos(tableName, pos)
	}
	if t != nil {
		log.V(2).Infof("Returning table %s from history for pos %s, schema %s", tableName, gtid, t)
	}
	return t, nil
}

// loadFromDB loads all rows from the schema_version table that the historian does not have as yet
// caller should have locked h.mu
func (h *historian) loadFromDB(ctx context.Context) error {
	conn, err := h.conns.Get(ctx, nil)
	if err != nil {
		return err
	}
	defer conn.Recycle()

	var tableData *sqltypes.Result
	if h.lastID == 0 && h.schemaMaxAgeSeconds > 0 { // only at vttablet start
		schemaMaxAge := time.Now().UTC().Add(time.Duration(-h.schemaMaxAgeSeconds) * time.Second)
		tableData, err = conn.Conn.Exec(ctx, sqlparser.BuildParsedQuery(getInitialSchemaVersions, sidecar.GetIdentifier(),
			schemaMaxAge.Unix()).Query, 10000, true)
	} else {
		tableData, err = conn.Conn.Exec(ctx, sqlparser.BuildParsedQuery(getNextSchemaVersions, sidecar.GetIdentifier(),
			h.lastID).Query, 10000, true)
	}

	if err != nil {
		log.Infof("Error reading schema_tracking table %v, will operate with the latest available schema", err)
		return nil
	}
	for _, row := range tableData.Rows {
		trackedSchema, id, err := h.readRow(row)
		if err != nil {
			return err
		}
		h.schemas = append(h.schemas, trackedSchema)
		h.lastID = id
	}

	if h.lastID != 0 && h.schemaMaxAgeSeconds > 0 {
		// To avoid keeping old schemas in memory which can lead to an eventual memory leak
		// we purge any older than h.schemaMaxAgeSeconds. Only needs to be done when adding
		// new schema rows.
		h.purgeOldSchemas()
	}

	h.sortSchemas()
	return nil
}

// readRow converts a row from the schema_version table to a trackedSchema
func (h *historian) readRow(row []sqltypes.Value) (*trackedSchema, int64, error) {
	id, _ := row[0].ToCastInt64()
	rowBytes, err := row[1].ToBytes()
	if err != nil {
		return nil, 0, err
	}
	pos, err := replication.DecodePosition(string(rowBytes))
	if err != nil {
		return nil, 0, err
	}
	rowBytes, err = row[2].ToBytes()
	if err != nil {
		return nil, 0, err
	}
	ddl := string(rowBytes)
	timeUpdated, err := row[3].ToCastInt64()
	if err != nil {
		return nil, 0, err
	}
	sch := &binlogdatapb.MinimalSchema{}
	rowBytes, err = row[4].ToBytes()
	if err != nil {
		return nil, 0, err
	}
	if err := sch.UnmarshalVT(rowBytes); err != nil {
		return nil, 0, err
	}
	log.V(vl).Infof("Read tracked schema from db: id %d, pos %v, ddl %s, schema len %d, time_updated %d \n",
		id, replication.EncodePosition(pos), ddl, len(sch.Tables), timeUpdated)

	tables := map[string]*binlogdatapb.MinimalTable{}
	for _, t := range sch.Tables {
		tables[t.Name] = t
	}
	tSchema := &trackedSchema{
		schema:      tables,
		pos:         pos,
		ddl:         ddl,
		timeUpdated: timeUpdated,
	}
	return tSchema, id, nil
}

func (h *historian) purgeOldSchemas() {
	maxAgeDuration := time.Duration(h.schemaMaxAgeSeconds) * time.Second
	shouldPurge := false

	// check if we have any schemas we need to purge and only create the filtered
	// slice if necessary
	for _, s := range h.schemas {
		if time.Since(time.Unix(s.timeUpdated, 0)) > maxAgeDuration {
			shouldPurge = true
			break
		}
	}

	if !shouldPurge {
		return
	}

	filtered := make([]*trackedSchema, 0)
	for _, s := range h.schemas {
		if time.Since(time.Unix(s.timeUpdated, 0)) < maxAgeDuration {
			filtered = append(filtered, s)
		}
	}
	h.schemas = filtered
}

// sortSchemas sorts entries in ascending order of gtid, ex: 40,44,48
func (h *historian) sortSchemas() {
	sort.Slice(h.schemas, func(i int, j int) bool {
		return h.schemas[j].pos.AtLeast(h.schemas[i].pos)
	})
}

// getTableFromHistoryForPos looks in the cache for a schema for a specific gtid
func (h *historian) getTableFromHistoryForPos(tableName sqlparser.IdentifierCS, pos replication.Position) *binlogdatapb.MinimalTable {
	idx := sort.Search(len(h.schemas), func(i int) bool {
		return pos.Equal(h.schemas[i].pos) || !pos.AtLeast(h.schemas[i].pos)
	})
	if idx >= len(h.schemas) || idx == 0 && !pos.Equal(h.schemas[idx].pos) { // beyond the range of the cache
		log.Infof("Schema not found in cache for %s with pos %s", tableName, pos)
		return nil
	}
	if pos.Equal(h.schemas[idx].pos) { //exact match to a cache entry
		return h.schemas[idx].schema[tableName.String()]
	}
	//not an exact match, so based on our sort algo idx is one less than found: from 40,44,48 : 43 < 44 but we want 40
	return h.schemas[idx-1].schema[tableName.String()]
}
