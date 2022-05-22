package feral

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

func (cat *FeralDruid) OnGCDReady(sim *core.Simulation) {
	cat.doRotation(sim)
}

func (cat *FeralDruid) doRotation(sim *core.Simulation) {
	// If we're out of form because we just cast Innervate, always shift
	if !cat.CatForm {
		 cat.Powershift.Cast(sim, nil)
		 return
	}

	// If we previously decided to shift, then execute the shift now once the input delay is over.
	if cat.readyToShift {
		cat.innervateOrShift(sim)
		return
	}

	// Get current Energy and CP
	energy := cat.CurrentEnergy()
	comboPoints := cat.ComboPoints()

	// Decide whether to cast Rip as our next special
	ripNow := cat.ShouldCastRip(sim, cat.Rotation)

	// Decide whether to cast Mangle as our next special
	mangleNow := !ripNow && !cat.MangleAura.IsActive()
}

func (cat *FeralDruid) innervateOrShift(sim *core.Simulation) {
	cat.waitingForTick = false

	// If we have just now decided to shift, then we do not execute the shift immediately, but instead trigger an input delay for realism.
	if !cat.readyToShift {
		cat.readyToShift = true
		return
	}

	cat.readyToShift = false

	// Logic for Innervate and Haste Pot usage will go here. For now we just execute simple powershifts without bundling any caster form CDs.
	cat.Powershift.Cast(sim, nil)
}
