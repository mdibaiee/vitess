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

package mysql

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mdibaiee/vitess/oracle/go/mysql/replication"
)

// Sample event data for MySQL 5.6.
var (
	mysql56FormatEvent = NewMysql56BinlogEvent([]byte{0x78, 0x4e, 0x49, 0x55, 0xf, 0x64, 0x0, 0x0, 0x0, 0x74, 0x0, 0x0, 0x0, 0x78, 0x0, 0x0, 0x0, 0x1, 0x0, 0x4, 0x0, 0x35, 0x2e, 0x36, 0x2e, 0x32, 0x34, 0x2d, 0x6c, 0x6f, 0x67, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x78, 0x4e, 0x49, 0x55, 0x13, 0x38, 0xd, 0x0, 0x8, 0x0, 0x12, 0x0, 0x4, 0x4, 0x4, 0x4, 0x12, 0x0, 0x0, 0x5c, 0x0, 0x4, 0x1a, 0x8, 0x0, 0x0, 0x0, 0x8, 0x8, 0x8, 0x2, 0x0, 0x0, 0x0, 0xa, 0xa, 0xa, 0x19, 0x19, 0x0, 0x1, 0x18, 0x4a, 0xf, 0xca})
	mysql56GTIDEvent   = NewMysql56BinlogEvent([]byte{0xff, 0x4e, 0x49, 0x55, 0x21, 0x64, 0x0, 0x0, 0x0, 0x30, 0x0, 0x0, 0x0, 0xf5, 0x2, 0x0, 0x0, 0x0, 0x0, 0x1, 0x43, 0x91, 0x92, 0xbd, 0xf3, 0x7c, 0x11, 0xe4, 0xbb, 0xeb, 0x2, 0x42, 0xac, 0x11, 0x3, 0x5a, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x48, 0x45, 0x82, 0x27})
	// This is the result of: begin; insert into customer values (1, "mlord@planetscale.com"), (2, "sup@planetscale.com"); commit;
	mysql56TransactionPayloadEvent = NewMysql56BinlogEvent([]byte{0xc7, 0xe1, 0x4b, 0x64, 0x28, 0x5b, 0xd2, 0xc7, 0x19, 0xdb, 0x00, 0x00, 0x00, 0x3a, 0x50, 0x00, 0x00, 0x00, 0x00, 0x02, 0x01, 0x00, 0x03, 0x03, 0xfc, 0xfe, 0x00, 0x01, 0x01, 0xb8, 0x00, 0x28, 0xb5, 0x2f, 0xfd, 0x00, 0x58, 0x64, 0x05, 0x00, 0xf2, 0x49, 0x23, 0x2a, 0xa0, 0x27, 0x69, 0x0c, 0xff, 0xe8, 0x06, 0xeb, 0xfe, 0xc3, 0xab, 0x8a, 0x7b, 0xc0, 0x36, 0x42, 0x5c, 0x6f, 0x1b, 0x2f, 0xfb, 0x6e, 0xc4, 0x9a, 0xe6, 0x6e, 0x6b, 0xda, 0x08, 0xf1, 0x37, 0x7e, 0xff, 0xb8, 0x6c, 0xbc, 0x27, 0x3c, 0xb7, 0x4f, 0xee, 0x14, 0xff, 0xaf, 0x09, 0x06, 0x69, 0xe3, 0x12, 0x68, 0x4a, 0x6e, 0xc3, 0xe1, 0x28, 0xaf, 0x3f, 0xc8, 0x14, 0x1c, 0xc3, 0x60, 0xce, 0xe3, 0x1e, 0x18, 0x4c, 0x63, 0xa1, 0x35, 0x90, 0x79, 0x04, 0xe8, 0xa9, 0xeb, 0x4a, 0x1b, 0xd7, 0x41, 0x53, 0x72, 0x17, 0xa4, 0x23, 0xa4, 0x47, 0x68, 0x00, 0xa2, 0x37, 0xee, 0xc1, 0xc7, 0x71, 0x30, 0x24, 0x19, 0xfd, 0x78, 0x49, 0x1b, 0x97, 0xd2, 0x94, 0xdc, 0x85, 0xa2, 0x21, 0xc1, 0xb0, 0x63, 0x8d, 0x7b, 0x0f, 0x32, 0x87, 0x07, 0xe2, 0x39, 0xf0, 0x7c, 0x3e, 0x01, 0xfe, 0x13, 0x8f, 0x11, 0xd0, 0x05, 0x9f, 0xbc, 0x18, 0x59, 0x91, 0x36, 0x2e, 0x6d, 0x4a, 0x6e, 0x0b, 0x00, 0x5e, 0x28, 0x10, 0xc0, 0x02, 0x50, 0x77, 0xe0, 0x64, 0x30, 0x02, 0x9e, 0x09, 0x54, 0xec, 0x80, 0x6d, 0x07, 0xa4, 0xc1, 0x7d, 0x60, 0xe4, 0x01, 0x78, 0x01, 0x01, 0x00, 0x00})
	mysql56QueryEvent              = NewMysql56BinlogEvent([]byte{0xff, 0x4e, 0x49, 0x55, 0x2, 0x64, 0x0, 0x0, 0x0, 0x77, 0x0, 0x0, 0x0, 0xdb, 0x3, 0x0, 0x0, 0x0, 0x0, 0x3d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x21, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6, 0x3, 0x73, 0x74, 0x64, 0x4, 0x8, 0x0, 0x8, 0x0, 0x21, 0x0, 0xc, 0x1, 0x74, 0x65, 0x73, 0x74, 0x0, 0x74, 0x65, 0x73, 0x74, 0x0, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x20, 0x69, 0x6e, 0x74, 0x6f, 0x20, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x20, 0x28, 0x6d, 0x73, 0x67, 0x29, 0x20, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x20, 0x28, 0x27, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x27, 0x29, 0x92, 0x12, 0x79, 0xc3})
	mysql56SemiSyncNoAckQueryEvent = NewMysql56BinlogEvent([]byte{0xef, 0x00, 0xff, 0x4e, 0x49, 0x55, 0x2, 0x64, 0x0, 0x0, 0x0, 0x77, 0x0, 0x0, 0x0, 0xdb, 0x3, 0x0, 0x0, 0x0, 0x0, 0x3d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x21, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6, 0x3, 0x73, 0x74, 0x64, 0x4, 0x8, 0x0, 0x8, 0x0, 0x21, 0x0, 0xc, 0x1, 0x74, 0x65, 0x73, 0x74, 0x0, 0x74, 0x65, 0x73, 0x74, 0x0, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x20, 0x69, 0x6e, 0x74, 0x6f, 0x20, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x20, 0x28, 0x6d, 0x73, 0x67, 0x29, 0x20, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x20, 0x28, 0x27, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x27, 0x29, 0x92, 0x12, 0x79, 0xc3})
	mysql56SemiSyncAckQueryEvent   = NewMysql56BinlogEvent([]byte{0xef, 0x01, 0xff, 0x4e, 0x49, 0x55, 0x2, 0x64, 0x0, 0x0, 0x0, 0x77, 0x0, 0x0, 0x0, 0xdb, 0x3, 0x0, 0x0, 0x0, 0x0, 0x3d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x21, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6, 0x3, 0x73, 0x74, 0x64, 0x4, 0x8, 0x0, 0x8, 0x0, 0x21, 0x0, 0xc, 0x1, 0x74, 0x65, 0x73, 0x74, 0x0, 0x74, 0x65, 0x73, 0x74, 0x0, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x20, 0x69, 0x6e, 0x74, 0x6f, 0x20, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x20, 0x28, 0x6d, 0x73, 0x67, 0x29, 0x20, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x20, 0x28, 0x27, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x27, 0x29, 0x92, 0x12, 0x79, 0xc3})
)

func TestMysql56IsGTID(t *testing.T) {
	if got, want := mysql56FormatEvent.IsGTID(), false; got != want {
		t.Errorf("%#v.IsGTID() = %#v, want %#v", mysql56FormatEvent, got, want)
	}
	if got, want := mysql56QueryEvent.IsGTID(), false; got != want {
		t.Errorf("%#v.IsGTID() = %#v, want %#v", mysql56QueryEvent, got, want)
	}
	if got, want := mysql56GTIDEvent.IsGTID(), true; got != want {
		t.Errorf("%#v.IsGTID() = %#v, want %#v", mysql56GTIDEvent, got, want)
	}
}

func TestMysql56StripChecksum(t *testing.T) {
	format, err := mysql56FormatEvent.Format()
	require.NoError(t, err, "Format() error: %v", err)

	stripped, gotChecksum, err := mysql56QueryEvent.StripChecksum(format)
	require.NoError(t, err, "StripChecksum() error: %v", err)

	// Check checksum.
	if want := []byte{0x92, 0x12, 0x79, 0xc3}; !reflect.DeepEqual(gotChecksum, want) {
		t.Errorf("checksum = %#v, want %#v", gotChecksum, want)
	}

	// Check query, to make sure checksum was stripped properly.
	// Query length is defined as "the rest of the bytes after offset X",
	// so the query will be wrong if the checksum is not stripped.
	gotQuery, err := stripped.Query(format)
	require.NoError(t, err, "Query() error: %v", err)

	if want := "insert into test_table (msg) values ('hello')"; string(gotQuery.SQL) != want {
		t.Errorf("query = %#v, want %#v", string(gotQuery.SQL), want)
	}
}

func TestMysql56GTID(t *testing.T) {
	format, err := mysql56FormatEvent.Format()
	require.NoError(t, err, "Format() error: %v", err)

	input, _, err := mysql56GTIDEvent.StripChecksum(format)
	require.NoError(t, err, "StripChecksum() error: %v", err)
	require.True(t, input.IsGTID(), "IsGTID() = false, want true")

	want := replication.Mysql56GTID{
		Server:   replication.SID{0x43, 0x91, 0x92, 0xbd, 0xf3, 0x7c, 0x11, 0xe4, 0xbb, 0xeb, 0x2, 0x42, 0xac, 0x11, 0x3, 0x5a},
		Sequence: 4,
	}
	got, hasBegin, err := input.GTID(format)
	require.NoError(t, err, "GTID() error: %v", err)
	assert.False(t, hasBegin, "GTID() returned hasBegin")
	assert.Equal(t, want, got, "GTID() = %#v, want %#v", got, want)
}

func TestMysql56DecodeTransactionPayload(t *testing.T) {
	format := NewMySQL56BinlogFormat()
	tableMap := &TableMap{}
	require.True(t, mysql56TransactionPayloadEvent.IsTransactionPayload())

	// The generated event is the result of the following SQL being executed in vtgate
	// against the commerce keyspace:
	// begin; insert into customer values (1, "mlord@planetscale.com"), (2, "sup@planetscale.com"); commit;
	// All of these below internal events are encoded in the compressed transaction
	// payload event.
	want := []string{
		"BEGIN",                     // Query event
		"vt_commerce.customer",      // TableMap event
		"[1 mlord@planetscale.com]", // WriteRows event
		"[2 sup@planetscale.com]",   // WriteRows event
		"COMMIT",                    // XID event
	}
	internalEvents, err := mysql56TransactionPayloadEvent.TransactionPayload(format)
	require.NoError(t, err)
	eventStrs := []string{}
	for _, ev := range internalEvents {
		switch {
		case ev.IsTableMap():
			tableMap, err = ev.TableMap(format)
			require.NoError(t, err)
			eventStrs = append(eventStrs, fmt.Sprintf("%s.%s", tableMap.Database, tableMap.Name))
		case ev.IsQuery():
			query, err := ev.Query(format)
			require.NoError(t, err)
			eventStrs = append(eventStrs, query.SQL)
		case ev.IsWriteRows():
			rows, err := ev.Rows(format, tableMap)
			require.NoError(t, err)
			for i := range rows.Rows {
				rowStr, err := rows.StringValuesForTests(tableMap, i)
				require.NoError(t, err)
				eventStrs = append(eventStrs, fmt.Sprintf("%v", rowStr))
			}
		case ev.IsXID():
			eventStrs = append(eventStrs, "COMMIT")
		}
	}
	require.Equal(t, want, eventStrs)
}

func TestMysql56ParsePosition(t *testing.T) {
	input := "00010203-0405-0607-0809-0a0b0c0d0e0f:1-2"

	sid := replication.SID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	var set replication.GTIDSet = replication.Mysql56GTIDSet{}
	set = set.AddGTID(replication.Mysql56GTID{Server: sid, Sequence: 1})
	set = set.AddGTID(replication.Mysql56GTID{Server: sid, Sequence: 2})
	want := replication.Position{GTIDSet: set}

	got, err := replication.ParsePosition(replication.Mysql56FlavorID, input)
	assert.NoError(t, err, "unexpected error: %v", err)
	assert.True(t, got.Equal(want), "(&mysql56{}).ParsePosition(%#v) = %#v, want %#v", input, got, want)

}

func TestMysql56SemiSyncAck(t *testing.T) {
	{
		c := Conn{ExpectSemiSyncIndicator: false}
		buf, semiSyncAckRequested, err := c.AnalyzeSemiSyncAckRequest(mysql56QueryEvent.Bytes())
		assert.NoError(t, err)
		e := NewMysql56BinlogEventWithSemiSyncInfo(buf, semiSyncAckRequested)

		assert.False(t, e.IsSemiSyncAckRequested())
		assert.True(t, e.IsQuery())
	}
	{
		c := Conn{ExpectSemiSyncIndicator: true}
		buf, semiSyncAckRequested, err := c.AnalyzeSemiSyncAckRequest(mysql56SemiSyncNoAckQueryEvent.Bytes())
		assert.NoError(t, err)
		e := NewMysql56BinlogEventWithSemiSyncInfo(buf, semiSyncAckRequested)

		assert.False(t, e.IsSemiSyncAckRequested())
		assert.True(t, e.IsQuery())
	}
	{
		c := Conn{ExpectSemiSyncIndicator: true}
		buf, semiSyncAckRequested, err := c.AnalyzeSemiSyncAckRequest(mysql56SemiSyncAckQueryEvent.Bytes())
		assert.NoError(t, err)
		e := NewMysql56BinlogEventWithSemiSyncInfo(buf, semiSyncAckRequested)

		assert.True(t, e.IsSemiSyncAckRequested())
		assert.True(t, e.IsQuery())
	}
}
