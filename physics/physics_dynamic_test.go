package physics

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCalcTractionForce(t *testing.T) {
	dir := &Vector{1, 0}
	ftr := calcTractionForce(dir, 10000)
	assert.Equal(t, 10000., ftr.X)
	assert.Equal(t, 0., ftr.Y)
}
func TestCalcAirResistanceForce(t *testing.T) {
	velocity := &Vector{0, 0}
	fair := calcAirResistForce(velocity)
	assert.Equal(t, 0., fair.X)
	assert.Equal(t, 0., fair.Y)
}

func TestCalcFrictionForce(t *testing.T) {
	weight := 1000.
	dir := &Vector{1, 0}
	ffr := calcFrictionForce(dir, weight)
	assert.Equal(t, -6000., ffr.X)
	assert.Equal(t, 0., ffr.Y)
}

func TestSumForces(t *testing.T) {
	weight := 1000.
	dir := &Vector{1, 0}
	velocity := &Vector{0, 0}
	ftr := calcTractionForce(dir, 10000)
	fair := calcAirResistForce(velocity)
	ffr := calcFrictionForce(dir, weight)
	f := ftr.Add(fair).Add(ffr)

	assert.Equal(t, 4000., f.X)
	assert.Equal(t, 0., f.Y)
}

func TestCalcAccelerate(t *testing.T) {
	weight := 1000.
	dir := &Vector{1, 0}
	velocity := &Vector{0, 0}
	ftr := calcTractionForce(dir, 10000)
	fair := calcAirResistForce(velocity)
	ffr := calcFrictionForce(dir, weight)
	f := ftr.Add(fair).Add(ffr)
	a := calcAccelerate(f, weight)
	assert.Equal(t, 4., a.X)
	assert.Equal(t, 0., a.Y)
}

func TestApplyAccelerateToVelocity(t *testing.T) {
	weight := 1000.
	dir := &Vector{1, 0}
	velocity := &Vector{0, 0}
	ftr := calcTractionForce(dir, 10000)
	fair := calcAirResistForce(velocity)
	ffr := calcFrictionForce(dir, weight)
	f := ftr.Add(fair).Add(ffr)
	a := calcAccelerate(f, weight)
	dt := time.Second
	newV := applyAccelerateToVelocity(velocity, a, dt)
	assert.Equal(t, 4., newV.X)
	assert.Equal(t, 0., newV.Y)
}

func TestApplyAccelerateToVelocityNegative(t *testing.T) {
	weight := 1000.
	dir := &Vector{1, 0}
	velocity := &Vector{0, 0}
	ftr := calcTractionForce(dir, 2000)
	fair := calcAirResistForce(velocity)
	ffr := calcFrictionForce(dir, weight)
	f := ftr.Add(fair).Add(ffr)
	a := calcAccelerate(f, weight)
	dt := time.Second
	newV := applyAccelerateToVelocity(velocity, a, dt)
	multiply := dir.MultiplyOnVector(newV)
	require.Less(t, multiply, 0.)
}

func TestApplyVelocityToPosition(t *testing.T) {
	weight := 1000.
	dir := &Vector{1, 0}
	velocity := &Vector{0, 0}
	ftr := calcTractionForce(dir, 10000)
	fair := calcAirResistForce(velocity)
	ffr := calcFrictionForce(dir, weight)
	f := ftr.Add(fair).Add(ffr)
	a := calcAccelerate(f, weight)
	dt := time.Second
	newV := applyAccelerateToVelocity(velocity, a, dt)
	newP := ApplyVelocityToPosition(&Point{1000, 1000}, newV, dt)
	assert.Equal(t, 1004., newP.X)
	assert.Equal(t, 1000., newP.Y)
}
