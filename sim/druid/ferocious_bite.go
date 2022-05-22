package druid

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var FerociousBiteActionID = core.ActionID{SpellID: 24248}
var FerociousBiteEnergyCost = 35.0

func (druid *Druid) registerFerociousBiteSpell(sim *core.Simulation) {
	druid.FerociousBite = druid.RegisterSpell(core.SpellConfig{
		ActionID:    FerociousBiteActionID,
		SpellSchool: core.SpellSchoolPhysical,
		SpellExtras: core.SpellExtrasMeleeMetrics,

		ResourceType: stats.Energy,
		BaseCost:     FerociousBiteEnergyCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: FerociousBiteEnergyCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1 + 0.03*float64(druid.Talents.FeralAggression),
			ThreatMultiplier: 1,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					comboPoints := float64(druid.ComboPoints())
					excessEnergy := druid.CurrentEnergy() - FerociousBiteEnergyCost
					base := 57.0 + 169.0*comboPoints + 4.1*excessEnergy
					roll := sim.RandomFloat("Ferocious Bite") * 66.0
					return base + roll + hitEffect.MeleeAttackPower(spell.Character)*0.05*comboPoints
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: core.OutcomeFuncMeleeSpecialHitAndCrit(druid.critMultiplier),
			OnSpellHit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.SpendComboPoints(sim, spell.ActionID)
				}
			},
		}),
	})
}

func (druid *Druid) ShouldCastBite(sim *core.Simulation, rotation proto.FeralDruid_Rotation) bool {
}
