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
	// via target_info (the "never-derive" behavior).
	JobInstanceOptionCFeatureGate = featuregate.GlobalRegistry().MustRegister(
		"receiver.prometheusreceiver.jobInstanceOptionC",
		featuregate.StageAlpha,
		featuregate.WithRegisterDescription("POC (Option C): store scrape identity as namespaced `prometheus.job`/`prometheus.instance` OTLP resource attributes, and stop defaulting service.name/service.instance.id from job/instance for undeclared targets."),
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
		return errors.New("only one of the receiver.prometheusreceiver.jobInstanceOption{A,B,C} feature gates may be enabled at a time")
	}
	return nil
}
