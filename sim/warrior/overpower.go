package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warrior *Warrior) registerOverpowerSpell() {
	actionID := core.ActionID{SpellID: 11585}

	warrior.RegisterAura(core.Aura{
		Label:    "Overpower Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHit: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeDodge) {
				warrior.overpowerValidUntil = sim.CurrentTime + time.Second*5
			}
		},
	})

	cost := 5 - float64(warrior.Talents.FocusedRage)
	refundAmount := cost * 0.8

	damageEffect := core.ApplyEffectFuncDirectDamage(core.SpellEffect{
		ProcMask: core.ProcMaskMeleeMHSpecial,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCritRating:  25 * core.MeleeCritRatingPerCritChance * float64(warrior.Talents.ImprovedOverpower),

		BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 35, 1, true),
		OutcomeApplier: core.OutcomeFuncMeleeSpecialNoBlockDodgeParry(warrior.critMultiplier(true)),

		OnSpellHit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				warrior.AddRage(sim, refundAmount, core.ActionID{OtherID: proto.OtherAction_OtherActionRefund})
			}
		},
	})

	warrior.Overpower = warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Target, spell *core.Spell) {
			warrior.overpowerValidUntil = 0
			damageEffect(sim, target, spell)
		},
	})
}

func (warrior *Warrior) ShouldOverpower(sim *core.Simulation) bool {
	return sim.CurrentTime < warrior.overpowerValidUntil &&
		warrior.Overpower.IsReady(sim) &&
		warrior.CurrentRage() >= warrior.Overpower.DefaultCast.Cost
}
