// Licensed to Adam Shannon under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. The Moov Authors licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Files []*os.File

func (fs Files) Close() {
	for i := range fs {
		err := fs[i].Close()
		if err != nil {
			log.Printf("WARN: closing %s failed: %v", fs[i].Name(), err)
		}
	}
}

func OpenPaths(paths []string) (Files, error) {
	var out Files
	for i := range paths {
		path, err := filepath.Abs(paths[i])
		if err != nil {
			return nil, fmt.Errorf("expanding %s failed: %v", paths[i], err)
		}
		fd, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("opening %s failed: %v", path, err)
		}
		out = append(out, fd)
	}
	return out, nil
}
