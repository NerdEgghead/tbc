package runner

import (
	"math"
	"math/rand"
	"sync"

	"github.com/wowsims/tbc/sim/core"
)

type StatWeightsResult struct {
	Weights       []float64
	WeightsStdev  []float64
	EpValues      []float64
	EpValuesStdev []float64
}

func CalcStatWeight(params IndividualParams) StatWeightsResult {
	baseSim := SetupIndividualSim(params)
	baseStats := baseSim.Raid.Parties[0].Players[0].Stats
	baselineResult := RunIndividualSim(baseSim)

	var waitGroup sync.WaitGroup
	result := StatWeightsResult{
		Weights:       make([]float64, core.StatLen),
		WeightsStdev:  make([]float64, core.StatLen),
		EpValues:      make([]float64, core.StatLen),
		EpValuesStdev: make([]float64, core.StatLen),
	}
	dpsHists := [core.StatLen]map[int32]int32{}

	doStat := func(stat core.Stat, value float64) {
		defer waitGroup.Done()

		newParams := params
		newParams.CustomStats = make([]float64, core.StatLen)
		newParams.CustomStats[stat] = value
		newSim := SetupIndividualSim(newParams)
		simResult := RunIndividualSim(newSim)
		result.Weights[stat] = (simResult.DpsAvg - baselineResult.DpsAvg) / value
		dpsHists[stat] = simResult.DpsHist
	}

	// Spell hit mod shouldn't go over hit cap.
	spellHitMod := math.Max(0, math.Min(10, 202-baseStats[core.StatSpellHit]))

	statMods := core.Stats{
		core.StatIntellect:  50,
		core.StatSpellPower: 50,
		core.StatSpellCrit:  50,
		core.StatSpellHit:   spellHitMod,
		core.StatSpellHaste: 50,
		core.StatMP5:        50,
	}

	for stat, mod := range statMods {
		if mod == 0 {
			continue
		}
		waitGroup.Add(1)
		go doStat(core.Stat(stat), mod)
	}

	waitGroup.Wait()

	for stat, mod := range statMods {
		if mod == 0 {
			continue
		}

		result.EpValues[stat] = result.Weights[stat] / result.Weights[core.StatSpellPower]
		result.WeightsStdev[stat] = computeStDevFromHists(params.Options.Iterations, mod, dpsHists[stat], baselineResult.DpsHist, nil, statMods[core.StatSpellPower])
		result.EpValuesStdev[stat] = computeStDevFromHists(params.Options.Iterations, mod, dpsHists[stat], baselineResult.DpsHist, dpsHists[core.StatSpellPower], statMods[core.StatSpellPower])
	}
	return result
}

func computeStDevFromHists(iters int, modValue float64, moddedStatDpsHist map[int32]int32, baselineDpsHist map[int32]int32, spellDmgDpsHist map[int32]int32, spellDmgModValue float64) float64 {
	sum := 0.0
	sumSquared := 0.0
	n := iters * 10
	for i := 0; i < n; {
		denominator := 1.0
		if spellDmgDpsHist != nil {
			denominator = float64(sampleFromDpsHist(spellDmgDpsHist, iters)-sampleFromDpsHist(baselineDpsHist, iters)) / spellDmgModValue
		}

		if denominator != 0 {
			ep := (float64(sampleFromDpsHist(moddedStatDpsHist, iters)-sampleFromDpsHist(baselineDpsHist, iters)) / modValue) / denominator
			sum += ep
			sumSquared += ep * ep
			i++
		}
	}
	epAvg := sum / float64(n)
	epStDev := math.Sqrt((sumSquared / float64(n)) - (epAvg * epAvg))
	return epStDev
}

func sampleFromDpsHist(hist map[int32]int32, histNumSamples int) int32 {
	r := rand.Float64()
	sampleIdx := int32(math.Floor(float64(histNumSamples) * r))

	var curSampleIdx int32
	for roundedDps, count := range hist {
		curSampleIdx += count
		if curSampleIdx >= sampleIdx {
			return roundedDps
		}
	}

	panic("Invalid dps histogram")
}
