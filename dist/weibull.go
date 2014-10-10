// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dist

import (
	"math"
	"math/rand"
)

// Weibull distribution. Valid range for x is [0,+∞).
type Weibull struct {
	// Shape parameter of the distribution. A value of 1 represents
	// the exponential distribution. A value of 2 represents the
	// Rayleigh distribution. Valid range is (0,+∞).
	K float64
	// Scale parameter of the distribution. Valid range is (0,+∞).
	Lambda float64
	// Source of random numbers
	Source *rand.Rand
}

// CDF computes the value of the cumulative density function at x.
func (w Weibull) CDF(x float64) float64 {
	if x < 0 {
		return 0
	} else {
		return 1.0 - math.Exp(-math.Pow(x/w.Lambda, w.K))
	}
}

// ConjugateUpdate updates the parameters of the distribution from the sufficient
// statistics of a set of samples. The sufficient statistics, suffStat, have been
// observed with nSamples observations. The prior values of the distribution are those
// currently in the distribution, and have been observed with priorStrength samples.
/*func (w *Weibull) ConjugateUpdate(suffStat []float64, nSamples float64, priorStrength []float64) {
	// TODO: Implement
}*/

// DLogProbDX returns the derivative of the log of the probability with
// respect to the input x.
//
// Special cases are:
//  DLogProbDX(0) = NaN
func (w Weibull) DLogProbDX(x float64) float64 {
	if x > 0 {
		return -(w.K*math.Pow(x/w.Lambda, w.K) + w.K - 1.0) / x
	}
	if x < 0 {
		return 0
	}
	return math.NaN()
}

// DLogProbDParam returns the derivative of the log of the probability with
// respect to the parameters of the distribution. The deriv slice must have length
// equal to the number of parameters of the distribution.
//
// The order is ∂LogProb / ∂K and then ∂LogProb / ∂λ.
//
// Special cases are:
//  The derivative at 0 is NaN.
func (w Weibull) DLogProbDParam(x float64, deriv []float64) {
	if len(deriv) != w.NumParameters() {
		panic("dist Weibull: slice length mismatch")
	}
	if x > 0 {
		deriv[0] = (1.0 - w.K*(math.Pow(x/w.Lambda, w.K)-1.0)*math.Log(x/w.Lambda)) / w.K
		deriv[1] = (w.K * (math.Pow(x/w.Lambda, w.K) - 1.0)) / w.Lambda
		return
	}
	if x < 0 {
		deriv[0] = 0
		deriv[1] = 0
		return
	}
	deriv[0] = math.NaN()
	deriv[0] = math.NaN()
	return
}

// Entropy returns the entropy of the distribution.
func (w Weibull) Entropy() float64 {
	return eulerGamma*(1.0-1.0/w.K) + math.Log(w.Lambda/w.K) + 1.0
}

// ExKurtosis returns the excess kurtosis of the distribution.
func (w Weibull) ExKurtosis() float64 {
	return (-6*w.gammaIPow(1, 4) + 12*w.gammaIPow(1, 2)*w.gammaIPow(2, 1) - 3*w.gammaIPow(2, 2) - 4*w.gammaIPow(1, 1)*w.gammaIPow(3, 1) + w.gammaIPow(4, 1)) / math.Pow(w.gammaIPow(2, 1)-w.gammaIPow(1, 2), 2)
}

// Fit sets the parameters of the probability distribution from the
// data samples x with relative weights w. If weights is nil, then all the weights
// are 1. If weights is not nil, then the len(weights) must equal len(samples).
/*func (w *Weibull) Fit(samples []float64, weights []float64) {
	// TODO: Implement
}*/

// gammaIPow aids in readability for the ExKurtosis, Mean,
// Skewness, and Variance calculations.
func (w Weibull) gammaIPow(i, pow float64) float64 {
	return math.Pow(math.Gamma(1+i/w.K), pow)
}

// LogProb computes the natural logarithm of the value of the probability
// density function at x. Zero is returned if x is less than zero.
//
// Special cases when x is zero are dependent on the shape parameter:
//  If 0 < K < 1, then log of the probability at 0 is +Inf.
//  If K == 1, then log of the probability at 0 is 0.
//  If K > 1, then log of the probability at 0 is -Inf.
func (w Weibull) LogProb(x float64) float64 {
	if x < 0 {
		return 0
	} else {
		return math.Log(w.K/w.Lambda) + (w.K-1)*math.Log(x/w.Lambda) + -math.Pow(x/w.Lambda, w.K)
	}
}

// MarshalParameters implements the ParameterMarshaler interface
func (w Weibull) MarshalParameters(p []Parameter) {
	nParam := w.NumParameters()
	if len(p) != nParam {
		panic("weibull: improper parameter length")
	}
	p[0].Name = "K"
	p[0].Value = w.K
	p[1].Name = "λ"
	p[1].Value = w.Lambda
	return
}

// Mean returns the mean of the probability distribution.
func (w Weibull) Mean() float64 {
	return w.Lambda * w.gammaIPow(1, 1)
}

// Median returns the median of the normal distribution.
func (w Weibull) Median() float64 {
	return w.Lambda * math.Pow(ln2, 1.0/w.K)
}

// Mode returns the mode of the normal distribution.
//
// Special case is:
//  The mode is NaN if the K (shape) parameter is less than 1.
func (w Weibull) Mode() float64 {
	if w.K > 1.0 {
		return w.Lambda * math.Pow((w.K-1.0)/w.K, 1.0/w.K)
	} else if w.K == 1.0 {
		return 0
	} else {
		return math.NaN()
	}
}

// NumParameters returns the number of parameters in the distribution.
func (Weibull) NumParameters() int {
	return 2
}

/*func (Weibull) NumSuffStat() int {
	// TODO: Implement
}*/

// Prob computes the value of the probability density function at x.
func (w Weibull) Prob(x float64) float64 {
	if x < 0 {
		return 0
	} else {
		return math.Exp(w.LogProb(x))
	}
}

func (w Weibull) Quantile(p float64) float64 {
	if p < 0 || p > 1 {
		panic("weibull: percentile out of bounds")
	}
	return w.Lambda * math.Pow(-math.Log(1-p), 1.0/w.K)
}

// Rand returns a random sample drawn from the distribution.
func (w Weibull) Rand() float64 {
	var rnd float64
	if w.Source == nil {
		rnd = rand.NormFloat64()
	} else {
		rnd = w.Source.NormFloat64()
	}
	return w.Quantile(rnd)
}

// Skewness returns the skewness of the distribution.
func (w Weibull) Skewness() float64 {
	stdDev := w.StdDev()
	return (w.gammaIPow(3, 1)*math.Pow(w.Lambda, 3) - 3*w.Mean()*math.Pow(stdDev, 2) - math.Pow(w.Mean(), 3)) / math.Pow(stdDev, 3)
}

// StdDev returns the standard deviation of the probability distribution.
func (w Weibull) StdDev() float64 {
	return math.Sqrt(w.Variance())
}

// SuffStat computes the sufficient statistics of a set of samples to update
// the distribution. The sufficient statistics are stored in place, and the
// effective number of samples are returned.
//
// If weights is nil, the weights are assumed to be 1, otherwise panics if
// len(samples) != len(weights).
/*func (Weibull) SuffStat(samples, weights, suffStat []float64) (nSamples float64) {
	// TODO: Implement
}*/

// Survival returns the survival function (complementary CDF) at x.
func (w Weibull) Survival(x float64) float64 {
	if x < 0 {
		return 1
	} else {
		return math.Exp(-math.Pow(x/w.Lambda, w.K))
	}
}

// UnmarshalParameters implements the ParameterMarshaler interface
func (w *Weibull) UnmarshalParameters(p []Parameter) {
	if len(p) != w.NumParameters() {
		panic("weibull: incorrect number of parameters to set")
	}
	if p[0].Name != "K" {
		panic("weibull: " + panicNameMismatch)
	}
	if p[1].Name != "λ" {
		panic("weibull: " + panicNameMismatch)
	}
	w.K = p[0].Value
	w.Lambda = p[1].Value
}

// Variance returns the variance of the probability distribution.
func (w Weibull) Variance() float64 {
	return math.Pow(w.Lambda, 2) * (w.gammaIPow(2, 1) - w.gammaIPow(1, 2))
}