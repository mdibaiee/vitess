#!/bin/bash

# Copyright 2019 The Vitess Authors.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
source $DIR/shell_functions.inc

# Normal builds run directly against the git repo, but when packaging (for example with rpms)
# a tar ball might be used, which will prevent the git metadata from being available.
# Should this be the case then allow environment variables to be used to source
# this information instead.
DEFAULT_BUILD_GIT_REV=$(git rev-parse HEAD)
DEFAULT_BUILD_GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
DEFAULT_BUILD_TIME=$(LC_ALL=C date)

echo "\
  -X 'mdibaiee/vitess/oracle/go/vt/servenv.buildHost=$(hostname)' \
  -X 'mdibaiee/vitess/oracle/go/vt/servenv.buildUser=$(whoami)' \
  -X 'mdibaiee/vitess/oracle/go/vt/servenv.buildGitRev=${BUILD_GIT_REV:-$DEFAULT_BUILD_GIT_REV}' \
  -X 'mdibaiee/vitess/oracle/go/vt/servenv.buildGitBranch=${BUILD_GIT_BRANCH:-$DEFAULT_BUILD_GIT_BRANCH}' \
  -X 'mdibaiee/vitess/oracle/go/vt/servenv.buildTime=${BUILD_TIME:-$DEFAULT_BUILD_TIME}' \
  -X 'mdibaiee/vitess/oracle/go/vt/servenv.jenkinsBuildNumberStr=${BUILD_NUMBER}' \
"
