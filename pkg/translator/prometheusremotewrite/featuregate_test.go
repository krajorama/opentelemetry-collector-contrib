// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheusremotewrite

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/featuregate"
)

func TestValidateJobInstanceOptionGates(t *testing.T) {
	require.NoError(t, ValidateJobInstanceOptionGates())

	require.NoError(t, featuregate.GlobalRegistry().Set(JobInstanceOptionAFeatureGate.ID(), true))
	t.Cleanup(func() {
		require.NoError(t, featuregate.GlobalRegistry().Set(JobInstanceOptionAFeatureGate.ID(), false))
	})
	require.NoError(t, ValidateJobInstanceOptionGates())

	require.NoError(t, featuregate.GlobalRegistry().Set(JobInstanceOptionCFeatureGate.ID(), true))
	t.Cleanup(func() {
		require.NoError(t, featuregate.GlobalRegistry().Set(JobInstanceOptionCFeatureGate.ID(), false))
	})
	require.Error(t, ValidateJobInstanceOptionGates())

	require.NoError(t, featuregate.GlobalRegistry().Set(JobInstanceOptionCFeatureGate.ID(), false))
	require.NoError(t, featuregate.GlobalRegistry().Set(JobInstanceOptionC1FeatureGate.ID(), true))
	t.Cleanup(func() {
		require.NoError(t, featuregate.GlobalRegistry().Set(JobInstanceOptionC1FeatureGate.ID(), false))
	})
	require.Error(t, ValidateJobInstanceOptionGates(), "option A and option C1 together must be rejected")
}
