// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"
)

// Function represents an instance of a function. This is the core interface
// that all functions must implement.
//
// NOTE: Provider-defined function support is in technical preview and offered
// without compatibility promises until Terraform 1.8 is generally available.
type Function interface {
	// Metadata should return the name of the function, such as parse_xyz.
	Metadata(context.Context, MetadataRequest, *MetadataResponse)

	// Definition should return the definition for the function.
	Definition(context.Context, DefinitionRequest, *DefinitionResponse)

	// Run should return the result of the function logic. It is called when
	// Terraform reaches a function call in the configuration. Argument data
	// values should be read from the [RunRequest] and the result value set in
	// the [RunResponse].
	Run(context.Context, RunRequest, *RunResponse)
}
