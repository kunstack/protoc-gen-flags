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

import (
	"strings"

	"github.com/spf13/pflag"
)

const DelimiterDot = "."

type Options struct {
	Prefix    []string
	Delimiter string
	Renamer   func(string) string
}

type Option func(*Options)

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

func WithDelimiter(delimiter string) Option {
	return func(o *Options) {
		o.Delimiter = delimiter
	}
}

func WithRenamer(renamer func(name string) string) Option {
	return func(o *Options) {
		o.Renamer = renamer
	}
}

func WithPrefix(prefix ...string) Option {
	var nonEmpty []string
	for _, part := range prefix {
		// Remove leading and trailing dots
		trimmed := strings.Trim(part, ".")
		if trimmed != "" {
			nonEmpty = append(nonEmpty, trimmed)
		}
	}
	return func(o *Options) {
		o.Prefix = append(o.Prefix, nonEmpty...)
	}
}

type NameBuilder struct {
	options *Options
}

func (n NameBuilder) Build(name string) string {
	return n.options.Renamer(strings.Join(append(n.options.Prefix, name), "."))
}

func NewNameBuilder(opts ...Option) NameBuilder {
	options := &Options{
		Delimiter: DelimiterDot,
		Renamer:   func(s string) string { return s },
	}
	for _, opt := range opts {
		opt(options)
	}
	return NameBuilder{options: options}
}
