/*
Copyright 2023 The Vitess Authors.

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

package evalengine

import (
	"mdibaiee/vitess/oracle/go/sqltypes"
)

type evalTuple struct {
	t []eval
}

var _ eval = (*evalTuple)(nil)

func (e *evalTuple) ToRawBytes() []byte {
	return nil
}

func (e *evalTuple) SQLType() sqltypes.Type {
	return sqltypes.Tuple
}

func (e *evalTuple) Size() int32 {
	return 0
}

func (e *evalTuple) Scale() int32 {
	return 0
}
