// Copyright 2015, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

// Package ewma implements exponentially weighted moving averages(EWMA).
// See https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average
// for formal definitions.
package ewma

import (
	"math"

	log "github.com/golang/glog"
)

const (
	// DefaultWeightingFactor is the default smoothing factor α
	// The value 0.875 is used in TCP RTT estimation
	DefaultWeightingFactor = 0.875
)

// EWMA is the class to calculate exponentially weighted moving average of a series of data
type EWMA struct {
	// The weighting factor α
	weightingFactor float64
	// The current value of the average
	currAverage float64
}

// NewEWMA returns a EWMA object which can calculate
// EWMA of a series of data gradually added to it
func NewEWMA(wf float64) *EWMA {
	if wf < 0 || wf > 1.0 {
		log.Infof(
			"Invalid weighting factor: %v, using default value(%v) instead.",
			wf,
			DefaultWeightingFactor,
		)
		wf = DefaultWeightingFactor
	}
	return &EWMA{
		weightingFactor: wf,
		currAverage:     math.NaN(),
	}
}

// AddValue adds a new data point into the EWMA calculation,
// by which the EWMA is automatically updated
func (e *EWMA) AddValue(value float64) {
	// the first value added
	if math.IsNaN(e.currAverage) {
		e.currAverage = value
	} else {
		e.currAverage = e.weightingFactor*e.currAverage + (1.0-e.weightingFactor)*value
	}
}

// GetEWMA returns the EWMA calculated from historical data
func (e *EWMA) GetEWMA() float64 {
	if math.IsNaN(e.currAverage) {
		// No value has been added yet, returning 0
		return 0
	}
	return e.currAverage
}

// Size makes EWMA to satisfy cache.Value interface
func (e *EWMA) Size() int {
	return 1
}
