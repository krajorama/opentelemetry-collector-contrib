// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver/internal"

import (
	"errors"

	"go.opentelemetry.io/collector/featuregate"
)

// These feature gates are a POC for the three options described in
// "Preserving Prometheus Job and Instance in OTLP Translation" for how the
// receiver should preserve Prometheus scrape identity (job/instance) as
// OTLP resource attributes alongside service.name/service.instance.id.
// They are mutually exclusive; see ValidateJobInstanceOptionGates.
var (
	// JobInstanceOptionAFeatureGate stores the scrape job/instance as the
	// bare `job`/`instance` resource attributes, in addition to today's
	// service.name/service.instance.id defaulting from job/instance.
	JobInstanceOptionAFeatureGate = featuregate.GlobalRegistry().MustRegister(
		"receiver.prometheusreceiver.jobInstanceOptionA",
		featuregate.StageAlpha,
		featuregate.WithRegisterDescription("POC (Option A): store scrape identity as bare `job`/`instance` OTLP resource attributes, in addition to today's service.name/service.instance.id defaulting."),
	)

	// JobInstanceOptionBFeatureGate stores the scrape job/instance as the
	// namespaced `prometheus.job`/`prometheus.instance` resource
	// attributes, in addition to today's service.name/service.instance.id
	// defaulting from job/instance.
	JobInstanceOptionBFeatureGate = featuregate.GlobalRegistry().MustRegister(
		"receiver.prometheusreceiver.jobInstanceOptionB",
		featuregate.StageAlpha,
		featuregate.WithRegisterDescription("POC (Option B): store scrape identity as namespaced `prometheus.job`/`prometheus.instance` OTLP resource attributes, in addition to today's service.name/service.instance.id defaulting."),
	)

	// JobInstanceOptionCFeatureGate stores the scrape job/instance as the
	// namespaced `prometheus.job`/`prometheus.instance` resource
	// attributes, and stops defaulting service.name/service.instance.id
	// from job/instance for targets that don't declare their own identity
	// via target_info (the "never-derive" behavior). It additionally
	// recognizes target_info's service_name/service_namespace/
	// service_instance_id (flattened, underscore) labels as declaring the
	// same identity as their dotted spelling ("covered-name recognition").
	JobInstanceOptionCFeatureGate = featuregate.GlobalRegistry().MustRegister(
		"receiver.prometheusreceiver.jobInstanceOptionC",
		featuregate.StageAlpha,
		featuregate.WithRegisterDescription("POC (Option C): store scrape identity as namespaced `prometheus.job`/`prometheus.instance` OTLP resource attributes, stop defaulting service.name/service.instance.id from job/instance for undeclared targets, and recognize flattened service_name/service_namespace/service_instance_id target_info labels as declared identity."),
	)

	// JobInstanceOptionC1FeatureGate is Option C without covered-name
	// recognition: target_info's service_name/service_namespace/
	// service_instance_id labels are treated as ordinary, unrecognized
	// Resource attributes rather than as declaring identity — only the
	// dotted spelling counts as a declared identity. Everything else
	// (namespaced pair storage, never-derive, the identity fallback) is
	// identical to Option C.
	JobInstanceOptionC1FeatureGate = featuregate.GlobalRegistry().MustRegister(
		"receiver.prometheusreceiver.jobInstanceOptionC1",
		featuregate.StageAlpha,
		featuregate.WithRegisterDescription("POC (Option C.1): Option C without covered-name recognition — only the dotted service.name/service.namespace/service.instance.id spelling on target_info counts as declared identity."),
	)
)

// ValidateJobInstanceOptionGates returns an error if more than one of the
// mutually exclusive job/instance identity option feature gates (Option
// A/B/C/C1) is enabled at the same time.
func ValidateJobInstanceOptionGates() error {
	enabled := 0
	for _, g := range []*featuregate.Gate{
		JobInstanceOptionAFeatureGate,
		JobInstanceOptionBFeatureGate,
		JobInstanceOptionCFeatureGate,
		JobInstanceOptionC1FeatureGate,
	} {
		if g.IsEnabled() {
			enabled++
		}
	}
	if enabled > 1 {
		return errors.New("only one of the receiver.prometheusreceiver.jobInstanceOption{A,B,C,C1} feature gates may be enabled at a time")
	}
	return nil
}
