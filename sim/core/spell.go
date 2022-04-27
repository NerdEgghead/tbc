package core

import (
	"fmt"
	"time"

	"github.com/wowsims/tbc/sim/core/stats"
)

type ApplySpellEffects func(*Simulation, *Target, *Spell)

type SpellConfig struct {
	// See definition of Spell (below) for comments on these.
	ActionID
	SpellSchool  SpellSchool
	SpellExtras  SpellExtras
	ResourceType stats.Stat
	BaseCost     float64

	Cast CastConfig

	ApplyEffects   ApplySpellEffects
	DisableMetrics bool
}

type SpellMetrics struct {
	// Metric totals for this spell, for the current iteration.
	Casts              int32
	Misses             int32
	Hits               int32
	Crits              int32
	Dodges             int32
	Glances            int32
	Parries            int32
	Blocks             int32
	PartialResists_1_4 int32   // 1/4 of the spell was resisted
	PartialResists_2_4 int32   // 2/4 of the spell was resisted
	PartialResists_3_4 int32   // 3/4 of the spell was resisted
	TotalDamage        float64 // Damage done by all casts of this spell.
	TotalThreat        float64 // Threat generated by all casts of this spell.
}

type Spell struct {
	// ID for this spell.
	ActionID

	// The unit who will perform this spell.
	Unit *Unit

	// Fire, Frost, Shadow, etc.
	SpellSchool SpellSchool

	// Flags
	SpellExtras SpellExtras

	// Should be stats.Mana, stats.Energy, stats.Rage, or unset.
	ResourceType stats.Stat

	// Base cost. Many effects in the game which 'reduce mana cost by X%'
	// are calculated using the base cost.
	BaseCost float64

	// Default cast parameters with all static effects applied.
	DefaultCast Cast

	CD       Cooldown
	SharedCD Cooldown

	// Performs a cast of this spell.
	castFn CastSuccessFunc

	SpellMetrics []SpellMetrics

	ApplyEffects ApplySpellEffects

	// The current or most recent cast data.
	CurCast Cast

	DisableMetrics bool
}

// Registers a new spell to the unit. Returns the newly created spell.
func (unit *Unit) RegisterSpell(config SpellConfig) *Spell {
	if len(unit.Spellbook) > 100 {
		panic(fmt.Sprintf("Over 100 registered spells when registering %s! There is probably a spell being registered every iteration.", config.ActionID))
	}

	spell := &Spell{
		ActionID:     config.ActionID,
		Unit:         unit,
		SpellSchool:  config.SpellSchool,
		SpellExtras:  config.SpellExtras,
		ResourceType: config.ResourceType,
		BaseCost:     config.BaseCost,

		DefaultCast: config.Cast.DefaultCast,
		CD:          config.Cast.CD,
		SharedCD:    config.Cast.SharedCD,

		ApplyEffects:   config.ApplyEffects,
		DisableMetrics: config.DisableMetrics,
	}

	spell.castFn = spell.makeCastFunc(config.Cast, spell.applyEffects)

	if spell.ApplyEffects == nil {
		spell.ApplyEffects = func(*Simulation, *Target, *Spell) {}
	}

	unit.Spellbook = append(unit.Spellbook, spell)

	return spell
}

// Returns the first registered spell with the given ID, or nil if there are none.
func (unit *Unit) GetSpell(actionID ActionID) *Spell {
	for _, spell := range unit.Spellbook {
		if spell.ActionID.SameAction(actionID) {
			return spell
		}
	}
	return nil
}

// Retrieves an existing spell with the same ID as the config uses, or registers it if there is none.
func (unit *Unit) GetOrRegisterSpell(config SpellConfig) *Spell {
	registered := unit.GetSpell(config.ActionID)
	if registered == nil {
		return unit.RegisterSpell(config)
	} else {
		return registered
	}
}

// Metrics for the current iteration
func (spell *Spell) CurDamagePerCast() float64 {
	if spell.SpellMetrics[0].Casts == 0 {
		return 0
	} else {
		casts := int32(0)
		damage := 0.0
		for _, targetMetrics := range spell.SpellMetrics {
			casts += targetMetrics.Casts
			damage += targetMetrics.TotalDamage
		}
		return damage / float64(casts)
	}
}

func (spell *Spell) reset(sim *Simulation) {
	spell.SpellMetrics = make([]SpellMetrics, sim.GetNumTargets())
}

func (spell *Spell) doneIteration() {
	if !spell.DisableMetrics {
		spell.Unit.Metrics.addSpell(spell)
	}
}

func (spell *Spell) IsReady(sim *Simulation) bool {
	return BothTimersReady(spell.CD.Timer, spell.SharedCD.Timer, sim)
}

func (spell *Spell) TimeToReady(sim *Simulation) time.Duration {
	return MaxTimeToReady(spell.CD.Timer, spell.SharedCD.Timer, sim)
}

func (spell *Spell) Cast(sim *Simulation, target *Target) bool {
	return spell.castFn(sim, target)
}

// Skips the actual cast and applies spell effects immediately.
func (spell *Spell) SkipCastAndApplyEffects(sim *Simulation, target *Target) {
	if sim.Log != nil {
		spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s)",
			spell.ActionID, spell.DefaultCast.Cost, 0)
		spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
	}
	spell.applyEffects(sim, target)
}

func (spell *Spell) applyEffects(sim *Simulation, target *Target) {
	if spell.SpellMetrics == nil {
		spell.reset(sim)
	}
	if target == nil {
		target = sim.GetPrimaryTarget()
	}
	spell.SpellMetrics[target.Index].Casts++
	spell.ApplyEffects(sim, target, spell)
}

func ApplyEffectFuncDirectDamage(baseEffect SpellEffect) ApplySpellEffects {
	if baseEffect.BaseDamage.Calculator == nil {
		// Just a hit check.
		return func(sim *Simulation, target *Target, spell *Spell) {
			effect := &baseEffect
			effect.Target = target
			effect.init(sim, spell)

			damage := 0.0
			effect.OutcomeApplier(sim, spell, effect, &damage)
			effect.triggerProcs(sim, spell)
		}
	} else {
		return func(sim *Simulation, target *Target, spell *Spell) {
			effect := &baseEffect
			effect.Target = target
			effect.init(sim, spell)

			damage := effect.calculateBaseDamage(sim, spell) * effect.DamageMultiplier
			effect.calcDamageSingle(sim, spell, damage)
			effect.finalize(sim, spell)
		}
	}
}

func ApplyEffectFuncDirectDamageTargetModifiersOnly(baseEffect SpellEffect) ApplySpellEffects {
	return func(sim *Simulation, target *Target, spell *Spell) {
		effect := &baseEffect
		effect.Target = target

		damage := effect.calculateBaseDamage(sim, spell) * effect.DamageMultiplier
		effect.calcDamageTargetOnly(sim, spell, damage)
		effect.finalize(sim, spell)
	}
}

func ApplyEffectFuncDamageMultiple(baseEffects []SpellEffect) ApplySpellEffects {
	if len(baseEffects) == 0 {
		panic("Multiple damage requires hits")
	} else if len(baseEffects) == 1 {
		return ApplyEffectFuncDirectDamage(baseEffects[0])
	}

	return func(sim *Simulation, _ *Target, spell *Spell) {
		for i := range baseEffects {
			effect := &baseEffects[i]
			effect.init(sim, spell)
			damage := effect.calculateBaseDamage(sim, spell) * effect.DamageMultiplier
			effect.calcDamageSingle(sim, spell, damage)
		}
		for i := range baseEffects {
			effect := &baseEffects[i]
			effect.finalize(sim, spell)
		}
	}
}
func ApplyEffectFuncDamageMultipleTargeted(baseEffects []SpellEffect) ApplySpellEffects {
	if len(baseEffects) == 0 {
		panic("Multiple damage requires hits")
	} else if len(baseEffects) == 1 {
		return ApplyEffectFuncDirectDamage(baseEffects[0])
	}

	return func(sim *Simulation, target *Target, spell *Spell) {
		for i := range baseEffects {
			effect := &baseEffects[i]
			effect.Target = target
			effect.init(sim, spell)
			damage := effect.calculateBaseDamage(sim, spell) * effect.DamageMultiplier
			effect.calcDamageSingle(sim, spell, damage)
		}
		for i := range baseEffects {
			effect := &baseEffects[i]
			effect.finalize(sim, spell)
		}
	}
}
func ApplyEffectFuncAOEDamage(sim *Simulation, baseEffect SpellEffect) ApplySpellEffects {
	numHits := sim.GetNumTargets()
	effects := make([]SpellEffect, numHits)
	for i := int32(0); i < numHits; i++ {
		effects[i] = baseEffect
		effects[i].Target = sim.GetTarget(i)
	}
	return ApplyEffectFuncDamageMultiple(effects)
}

func ApplyEffectFuncDot(dot *Dot) ApplySpellEffects {
	return func(sim *Simulation, _ *Target, _ *Spell) {
		dot.Apply(sim)
	}
}

// AOE Cap Mechanics:
// http://web.archive.org/web/20081023033855/http://elitistjerks.com/f47/t25902-aoe_spell_cap_mechanics/
func applyAOECap(effects []SpellEffect, outcomeMultipliers []float64, aoeCap float64) {
	// Increased damage from crits doesn't count towards the cap, so need to
	// tally pre-crit damage.
	totalTowardsCap := 0.0
	numHits := 0
	for i, _ := range effects {
		effect := &effects[i]
		if effect.Landed() {
			numHits++
			if effect.Outcome.Matches(OutcomeCrit) {
				totalTowardsCap += effect.Damage / outcomeMultipliers[i]
			} else {
				totalTowardsCap += effect.Damage
			}
		}
	}

	if totalTowardsCap <= aoeCap {
		return
	}

	maxDamagePerHit := aoeCap / float64(numHits)
	for i, _ := range effects {
		effect := &effects[i]
		if effect.Landed() {
			if effect.Outcome.Matches(OutcomeCrit) {
				preCritDamage := effect.Damage / outcomeMultipliers[i]
				capped := MinFloat(preCritDamage, maxDamagePerHit)
				effect.Damage = capped * outcomeMultipliers[i]
			} else {
				effect.Damage = MinFloat(effect.Damage, maxDamagePerHit)
			}
		}
	}
}
func ApplyEffectFuncDamageMultipleAOECapped(sim *Simulation, baseEffect SpellEffect, aoeCap float64) ApplySpellEffects {
	numHits := sim.GetNumTargets()
	if numHits == 0 {
		return nil
	} else if numHits == 1 {
		return ApplyEffectFuncDirectDamage(baseEffect)
	} else if numHits < 4 {
		// Just assume its impossible to hit AOE cap with <4 targets.
		return ApplyEffectFuncAOEDamage(sim, baseEffect)
	}

	baseEffects := make([]SpellEffect, numHits)
	outcomeMultipliers := make([]float64, numHits)
	for i := int32(0); i < numHits; i++ {
		baseEffects[i] = baseEffect
		baseEffects[i].Target = sim.GetTarget(i)
	}

	return func(sim *Simulation, _ *Target, spell *Spell) {
		for i := range baseEffects {
			effect := &baseEffects[i]
			effect.init(sim, spell)
			damage := effect.calculateBaseDamage(sim, spell) * effect.DamageMultiplier

			effect.applyAttackerModifiers(sim, spell, &damage)
			effect.applyResistances(sim, spell, &damage)
			damageBefore := damage
			effect.OutcomeApplier(sim, spell, effect, &damage)
			outcomeMultipliers[i] = damage / damageBefore
			effect.Damage = damage
		}
		applyAOECap(baseEffects, outcomeMultipliers, aoeCap)
		for i := range baseEffects {
			effect := &baseEffects[i]
			effect.applyTargetModifiers(sim, spell, effect.BaseDamage.TargetSpellCoefficient, &effect.Damage)
		}
		for i := range baseEffects {
			effect := &baseEffects[i]
			effect.finalize(sim, spell)
		}
	}
}
