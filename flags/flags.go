// Copyright 2021 Aapeli <aapeli.nian@gmail.com> All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flags

import "github.com/spf13/pflag"

// This package provides protobuf extensions for generating AddFlags methods
// using the protoc-gen-flags plugin.

// Defaulter is an interface that types can implement to provide default values
// for their fields. The SetDefaults method should be called to initialize
// all fields with their configured default values.
type Defaulter interface {
	SetDefaults()
}

// Flagger is an interface that protobuf-generated structs implement to expose
// their fields as command-line flags. The AddFlags method binds the struct's
// fields to a pflag.FlagSet, allowing automatic generation of CLI flags
// from protobuf message definitions.
//
// Parameters:
//   - fs: The pflag.FlagSet to which flags will be added
//   - prefix: Optional prefix strings that will be prepended to flag names
//     for hierarchical flag organization (e.g., "server", "database")
type Flagger interface {
	AddFlags(fs *pflag.FlagSet, prefix ...string)
}
