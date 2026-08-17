// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheusremotewrite // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite"

import (
	"errors"

	"go.opentelemetry.io/collector/featuregate"
)

// Bare and namespaced resource attribute names used to look up preserved
// Prometheus scrape identity, per the job/instance identity option feature
// gates below.
const (
	bareJobAttr            = "job"
	bareInstanceAttr       = "instance"
	namespacedJobAttr      = "prometheus.job"
	namespacedInstanceAttr = "prometheus.instance"
)

// These feature gates are a POC for the three options described in
// "Preserving Prometheus Job and Instance in OTLP Translation" for how the
// OTLP -> Prometheus direction should derive job/instance labels from
// resource attributes. They are mutually exclusive; see
// ValidateJobInstanceOptionGates.
var (
	// JobInstanceOptionAFeatureGate makes job/instance label derivation look
	// up the bare `job`/`instance` resource attributes first, falling back
	// to today's service.name/service.namespace/service.instance.id
	// derivation when absent.
	JobInstanceOptionAFeatureGate = featuregate.GlobalRegistry().MustRegister(
		"exporter.prometheusremotewrite.jobInstanceOptionA",
		featuregate.StageAlpha,
		featuregate.WithRegisterDescription("POC (Option A): derive job/instance labels from bare `job`/`instance` resource attributes first, falling back to service.name/service.namespace/service.instance.id."),
	)

	// JobInstanceOptionBFeatureGate makes job/instance label derivation look
	// up the namespaced `prometheus.job`/`prometheus.instance` resource
	// attributes first, falling back to today's
	// service.name/service.namespace/service.instance.id derivation when
	// absent.
	JobInstanceOptionBFeatureGate = featuregate.GlobalRegistry().MustRegister(
		"exporter.prometheusremotewrite.jobInstanceOptionB",
		featuregate.StageAlpha,
		featuregate.WithRegisterDescription("POC (Option B): derive job/instance labels from namespaced `prometheus.job`/`prometheus.instance` resource attributes first, falling back to service.name/service.namespace/service.instance.id."),
	)

	// JobInstanceOptionCFeatureGate makes job/instance label derivation
	// prefer a Resource's declared identity (service.name and/or
	// service.instance.id); the namespaced `prometheus.job`/
	// `prometheus.instance` resource attributes are only used as an
	// identity fallback when the Resource declares neither.
	JobInstanceOptionCFeatureGate = featuregate.GlobalRegistry().MustRegister(
		"exporter.prometheusremotewrite.jobInstanceOptionC",
		featuregate.StageAlpha,
		featuregate.WithRegisterDescription("POC (Option C): prefer declared service.name/service.instance.id identity; fall back to namespaced `prometheus.job`/`prometheus.instance` resource attributes only when neither is declared."),
	)
)

// ValidateJobInstanceOptionGates returns an error if more than one of the
// mutually exclusive job/instance identity option feature gates (Option
// A/B/C) is enabled at the same time.
func ValidateJobInstanceOptionGates() error {
	enabled := 0
	for _, g := range []*featuregate.Gate{
		JobInstanceOptionAFeatureGate,
		JobInstanceOptionBFeatureGate,
		JobInstanceOptionCFeatureGate,
	} {
		if g.IsEnabled() {
			enabled++
		}
	}
	if enabled > 1 {
		return errors.New("only one of the exporter.prometheusremotewrite.jobInstanceOption{A,B,C} feature gates may be enabled at a time")
	}
	return nil
}
