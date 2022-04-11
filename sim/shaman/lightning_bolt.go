package shaman

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/items"
)

const SpellIDLB12 int32 = 25449

// newLightningBoltTemplate returns a cast generator for Lightning Bolt with as many fields precomputed as possible.
func (shaman *Shaman) newLightningBoltSpell(sim *core.Simulation, isLightningOverload bool) *core.Spell {
	baseCost := 300.0
	if shaman.Equip[items.ItemSlotRanged].ID == TotemOfThePulsingEarth {
		baseCost -= 27.0
	}

	spellConfig := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: SpellIDLB12},
		baseCost,
		time.Millisecond*2500,
		isLightningOverload)

	spellConfig.Cast.ModifyCast = func(sim *core.Simulation, spell *core.Spell, cast *core.NewCast) {
		shaman.applyElectricSpellCastInitModifiers(spell, cast)
		if shaman.NaturesSwiftnessAura != nil && shaman.NaturesSwiftnessAura.IsActive() {
			cast.CastTime = 0
		}
	}

	effect := shaman.newElectricSpellEffect(571, 652, 0.794, isLightningOverload)

	if ItemSetSkyshatterRegalia.CharacterHasSetBonus(&shaman.Character, 4) {
		effect.DamageMultiplier *= 1.05
	}

	if !isLightningOverload && shaman.Talents.LightningOverload > 0 {
		lightningOverloadChance := float64(shaman.Talents.LightningOverload) * 0.04
		effect.OnSpellHit = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			if shaman.Talents.ElementalFocus && spellEffect.Outcome.Matches(core.OutcomeCrit) {
				shaman.ElementalFocusStacks = 2
			}

			if sim.RandomFloat("LB Lightning Overload") > lightningOverloadChance {
				return
			}
			shaman.LightningBoltLO.Cast(sim, spellEffect.Target)
		}
	} else {
		effect.OnSpellHit = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if shaman.Talents.ElementalFocus && spellEffect.Outcome.Matches(core.OutcomeCrit) {
				shaman.ElementalFocusStacks = 2
			}
		}
	}

	spellConfig.ApplyEffects = core.ApplyEffectFuncDirectDamage(effect)
	return shaman.RegisterSpell(spellConfig)
}
