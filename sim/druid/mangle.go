package druid

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

var MangleActionID = core.ActionID{SpellID: 33983}

func (druid *Druid) registerMangleSpell(sim *core.Simulation) {
	druid.MangleAura = core.MangleAura(sim.GetPrimaryTarget())

	if druid.Rotation.MangleBot {
		druid.MangleAura = core.MakePermanent(druid.MangleAura)
	}

	energyCost := 45.0 - float64(druid.Talents.Ferocity) - core.TernaryFloat64(ItemSetThunderheartFeral.CharacterHasSetBonus(&druid.Character, 2), 5.0, 0)
	refundAmount := energyCost * 0.8

	druid.Mangle = druid.RegisterSpell(core.SpellConfig{
		ActionID:    MangleActionID,
		SpellSchool: core.SpellSchoolPhysical,
		SpellExtras: core.SpellExtrasMeleeMetrics,

		ResourceType: stats.Energy,
		BaseCost:     energyCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: energyCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1 + 0.1*float64(druid.Talents.SavageFury),
			ThreatMultiplier: 1,
			BaseDamage: core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 264.0, 1.6, true),
			OutcomeApplier: core.OutcomeFuncMeleeSpecialHitAndCrit(druid.critMultiplier),
			OnSpellHit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.AddComboPoints(sim, 1, MangleActionID)
					druid.MangleAura.Activate()
				} else {
					druid.AddEnergy(sim, refundAmount, core.ActionID{OtherID: proto.OtherAction_OtherActionRefund})
				}
			},
		}),
	})
}
