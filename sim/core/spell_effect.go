package core

import (
	"fmt"
	"math"

	"github.com/wowsims/tbc/sim/core/stats"
)

// Callback for after a spell hits the target and after damage is calculated. Use it for proc effects
// or anything that comes from the final result of the spell.
type OnSpellHit func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect)

// OnPeriodicDamage is called when dots tick, after damage is calculated. Use it for proc effects
// or anything that comes from the final result of a tick.
type OnPeriodicDamage func(sim *Simulation, spell *Spell, spellEffect *SpellEffect)

type SpellEffect struct {
	// Target of the spell.
	Target *Target

	BaseDamage     BaseDamageConfig
	OutcomeApplier OutcomeApplier

	// Bonus stats to be added to the spell.
	BonusSpellHitRating  float64
	BonusSpellPower      float64
	BonusSpellCritRating float64

	BonusAttackPower float64
	BonusCritRating  float64

	// Additional multiplier that is always applied.
	DamageMultiplier float64

	// Multiplier for all threat generated by this effect.
	ThreatMultiplier float64

	// Adds a fixed amount of threat to this spell, before multipliers.
	FlatThreatBonus float64

	// TODO: Should be able to remove this after refactoring is done.
	IsPeriodic bool

	// Whether this is a phantom cast. Phantom casts are usually casts triggered by some effect,
	// like The Lightning Capacitor or Shaman Flametongue Weapon. Many on-hit effects do not
	// proc from phantom casts, only regular casts.
	IsPhantom bool

	// Controls which effects can proc from this effect.
	ProcMask ProcMask

	// Callbacks for providing additional custom behavior.
	OnInit           func(sim *Simulation, spell *Spell, spellEffect *SpellEffect)
	OnSpellHit       func(sim *Simulation, spell *Spell, spellEffect *SpellEffect)
	OnPeriodicDamage func(sim *Simulation, spell *Spell, spellEffect *SpellEffect)

	// Results
	Outcome HitOutcome
	Damage  float64 // Damage done by this cast.
	Threat  float64

	PreoutcomeDamage float64 // Damage done by this cast.
}

func (spellEffect *SpellEffect) Landed() bool {
	return spellEffect.Outcome.Matches(OutcomeLanded)
}

func (spellEffect *SpellEffect) TotalThreatMultiplier(character *Character) float64 {
	return spellEffect.ThreatMultiplier * character.PseudoStats.ThreatMultiplier
}

func (spellEffect *SpellEffect) MeleeAttackPower(character *Character) float64 {
	return character.stats[stats.AttackPower] + character.PseudoStats.MobTypeAttackPower + spellEffect.BonusAttackPower
}

func (spellEffect *SpellEffect) MeleeAttackPowerOnTarget() float64 {
	return spellEffect.Target.PseudoStats.BonusMeleeAttackPower
}

func (spellEffect *SpellEffect) RangedAttackPower(character *Character) float64 {
	return character.stats[stats.RangedAttackPower] + character.PseudoStats.MobTypeAttackPower + spellEffect.BonusAttackPower
}

func (spellEffect *SpellEffect) RangedAttackPowerOnTarget() float64 {
	return spellEffect.Target.PseudoStats.BonusRangedAttackPower
}

func (spellEffect *SpellEffect) BonusWeaponDamage(character *Character) float64 {
	return character.PseudoStats.BonusDamage
}

func (spellEffect *SpellEffect) ExpertisePercentage(character *Character) float64 {
	expertiseRating := character.stats[stats.Expertise]
	if spellEffect.ProcMask.Matches(ProcMaskMeleeMH) {
		expertiseRating += character.PseudoStats.BonusMHExpertiseRating
	} else if spellEffect.ProcMask.Matches(ProcMaskMeleeOH) {
		expertiseRating += character.PseudoStats.BonusOHExpertiseRating
	}
	return math.Floor(expertiseRating/ExpertisePerQuarterPercentReduction) / 400
}

func (spellEffect *SpellEffect) PhysicalHitChance(character *Character) float64 {
	hitRating := character.stats[stats.MeleeHit] + spellEffect.Target.PseudoStats.BonusMeleeHitRating

	if spellEffect.ProcMask.Matches(ProcMaskRanged) {
		hitRating += character.PseudoStats.BonusRangedHitRating
	}

	return (hitRating / (MeleeHitRatingPerHitChance * 100)) - spellEffect.Target.HitSuppression
}

func (spellEffect *SpellEffect) PhysicalCritChance(character *Character, spell *Spell) float64 {
	critRating := character.stats[stats.MeleeCrit] + spellEffect.BonusCritRating + spellEffect.Target.PseudoStats.BonusCritRating

	if spellEffect.ProcMask.Matches(ProcMaskRanged) {
		critRating += character.PseudoStats.BonusRangedCritRating
	} else {
		critRating += character.PseudoStats.BonusMeleeCritRating
	}
	if spell.SpellExtras.Matches(SpellExtrasAgentReserved1) {
		critRating += character.PseudoStats.BonusCritRatingAgentReserved1
	}
	if spellEffect.ProcMask.Matches(ProcMaskMeleeMH) {
		critRating += character.PseudoStats.BonusMHCritRating
	} else if spellEffect.ProcMask.Matches(ProcMaskMeleeOH) {
		critRating += character.PseudoStats.BonusOHCritRating
	}

	return (critRating / (MeleeCritRatingPerCritChance * 100)) - spellEffect.Target.CritSuppression
}

func (spellEffect *SpellEffect) SpellPower(character *Character, spell *Spell) float64 {
	return character.GetStat(stats.SpellPower) + character.GetStat(spell.SpellSchool.Stat()) + character.PseudoStats.MobTypeSpellPower + spellEffect.BonusSpellPower
}

func (spellEffect *SpellEffect) SpellCritChance(character *Character, spell *Spell) float64 {
	critRating := (character.GetStat(stats.SpellCrit) + spellEffect.BonusSpellCritRating + spellEffect.Target.PseudoStats.BonusCritRating)
	if spell.SpellSchool.Matches(SpellSchoolFire) {
		critRating += character.PseudoStats.BonusFireCritRating
	} else if spell.SpellSchool.Matches(SpellSchoolFrost) {
		critRating += spellEffect.Target.PseudoStats.BonusFrostCritRating
	}
	return critRating / (SpellCritRatingPerCritChance * 100)
}

func (spellEffect *SpellEffect) init(sim *Simulation, spell *Spell) {
	if spellEffect.OnInit != nil {
		spellEffect.OnInit(sim, spell, spellEffect)
	}
}

func (spellEffect *SpellEffect) calculateBaseDamage(sim *Simulation, spell *Spell) float64 {
	if spellEffect.BaseDamage.Calculator == nil {
		return 0
	} else {
		return spellEffect.BaseDamage.Calculator(sim, spellEffect, spell)
	}
}

func (spellEffect *SpellEffect) calcDamageSingle(sim *Simulation, spell *Spell, damage float64) {
	if !spell.SpellExtras.Matches(SpellExtrasIgnoreModifiers) {
		spellEffect.applyAttackerModifiers(sim, spell, &damage)
		spellEffect.applyTargetModifiers(sim, spell, spellEffect.BaseDamage.TargetSpellCoefficient, &damage)
		spellEffect.applyResistances(sim, spell, &damage)
		spellEffect.PreoutcomeDamage = damage
		spellEffect.OutcomeApplier(sim, spell, spellEffect, &damage)
	}
	spellEffect.Damage = damage
}
func (spellEffect *SpellEffect) calcDamageTargetOnly(sim *Simulation, spell *Spell, damage float64) {
	spellEffect.applyTargetModifiers(sim, spell, spellEffect.BaseDamage.TargetSpellCoefficient, &damage)
	spellEffect.applyResistances(sim, spell, &damage)
	spellEffect.OutcomeApplier(sim, spell, spellEffect, &damage)
	spellEffect.Damage = damage
}

func (spellEffect *SpellEffect) calcThreat(character *Character) float64 {
	if spellEffect.Landed() {
		return (spellEffect.Damage + spellEffect.FlatThreatBonus) * spellEffect.TotalThreatMultiplier(character)
	} else {
		return 0
	}
}

func (spellEffect *SpellEffect) finalize(sim *Simulation, spell *Spell) {
	spell.TotalDamage += spellEffect.Damage
	spell.TotalThreat += spellEffect.calcThreat(spell.Character)
	spellEffect.triggerProcs(sim, spell)
}

func (spellEffect *SpellEffect) triggerProcs(sim *Simulation, spell *Spell) {
	if sim.Log != nil {
		if spellEffect.IsPeriodic {
			spell.Character.Log(sim, "%s tick %s. (Threat: %0.3f)", spell.ActionID, spellEffect, spellEffect.calcThreat(spell.Character))
		} else {
			spell.Character.Log(sim, "%s %s. (Threat: %0.3f)", spell.ActionID, spellEffect, spellEffect.calcThreat(spell.Character))
		}
	}

	if !spellEffect.IsPeriodic {
		if spellEffect.OnSpellHit != nil {
			spellEffect.OnSpellHit(sim, spell, spellEffect)
		}
		spell.Character.OnSpellHit(sim, spell, spellEffect)
		spellEffect.Target.OnSpellHit(sim, spell, spellEffect)
	} else {
		if spellEffect.OnPeriodicDamage != nil {
			spellEffect.OnPeriodicDamage(sim, spell, spellEffect)
		}
		spell.Character.OnPeriodicDamage(sim, spell, spellEffect)
		spellEffect.Target.OnPeriodicDamage(sim, spell, spellEffect)
	}
}

func (spellEffect *SpellEffect) String() string {
	outcomeStr := spellEffect.Outcome.String()
	if !spellEffect.Landed() {
		return outcomeStr
	}
	return fmt.Sprintf("%s for %0.3f damage", outcomeStr, spellEffect.Damage)
}

func (spellEffect *SpellEffect) applyAttackerModifiers(sim *Simulation, spell *Spell, damage *float64) {
	attacker := spell.Character

	if spellEffect.ProcMask.Matches(ProcMaskRanged) {
		*damage *= attacker.PseudoStats.RangedDamageDealtMultiplier
	}
	if spell.SpellExtras.Matches(SpellExtrasAgentReserved1) {
		*damage *= attacker.PseudoStats.AgentReserved1DamageDealtMultiplier
	}

	*damage *= attacker.PseudoStats.DamageDealtMultiplier
	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		*damage *= attacker.PseudoStats.PhysicalDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolArcane) {
		*damage *= attacker.PseudoStats.ArcaneDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolFire) {
		*damage *= attacker.PseudoStats.FireDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolFrost) {
		*damage *= attacker.PseudoStats.FrostDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolHoly) {
		*damage *= attacker.PseudoStats.HolyDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolNature) {
		*damage *= attacker.PseudoStats.NatureDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolShadow) {
		*damage *= attacker.PseudoStats.ShadowDamageDealtMultiplier
	}
}

func (spellEffect *SpellEffect) applyTargetModifiers(sim *Simulation, spell *Spell, targetCoeff float64, damage *float64) {
	target := spellEffect.Target

	*damage *= target.PseudoStats.DamageTakenMultiplier
	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		if targetCoeff > 0 {
			*damage += target.PseudoStats.BonusPhysicalDamageTaken
		}
		*damage *= target.PseudoStats.PhysicalDamageTakenMultiplier
		if spellEffect.IsPeriodic {
			*damage *= target.PseudoStats.PeriodicPhysicalDamageTakenMultiplier
		}
	} else if spell.SpellSchool.Matches(SpellSchoolArcane) {
		*damage *= target.PseudoStats.ArcaneDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolFire) {
		*damage *= target.PseudoStats.FireDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolFrost) {
		*damage *= target.PseudoStats.FrostDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolHoly) {
		*damage += target.PseudoStats.BonusHolyDamageTaken * targetCoeff
		*damage *= target.PseudoStats.HolyDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolNature) {
		*damage *= target.PseudoStats.NatureDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolShadow) {
		*damage *= target.PseudoStats.ShadowDamageTakenMultiplier
	}
}

// Modifies damage based on Armor or Magic resistances, depending on the damage type.
func (spellEffect *SpellEffect) applyResistances(sim *Simulation, spell *Spell, damage *float64) {
	if spell.SpellExtras.Matches(SpellExtrasIgnoreResists) {
		return
	}

	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		// Physical resistance (armor).
		*damage *= spellEffect.Target.ArmorDamageReduction(spell.Character.stats[stats.ArmorPenetration])
	} else if !spell.SpellExtras.Matches(SpellExtrasBinary) {
		// Magical resistance.
		// https://royalgiraffe.github.io/resist-guide

		resistanceRoll := sim.RandomFloat("Partial Resist")
		if resistanceRoll > 0.18 { // 13% chance for 25% resist, 4% for 50%, 1% for 75%
			// No partial resist.
		} else if resistanceRoll > 0.05 {
			spellEffect.Outcome |= OutcomePartial1_4
			*damage *= 0.75
		} else if resistanceRoll > 0.01 {
			spellEffect.Outcome |= OutcomePartial2_4
			*damage *= 0.5
		} else {
			spellEffect.Outcome |= OutcomePartial3_4
			*damage *= 0.25
		}
	}
}
