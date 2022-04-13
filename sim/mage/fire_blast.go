package mage

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

const SpellIDFireBlast int32 = 27079

var FireBlastCooldownID = core.NewCooldownID()
var FireBlastActionID = core.ActionID{SpellID: SpellIDFireBlast, CooldownID: FireBlastCooldownID}

func (mage *Mage) registerFireBlastSpell(sim *core.Simulation) {
	baseCost := 465.0

	mage.FireBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:    FireBlastActionID,
		SpellSchool: core.SpellSchoolFire,
		SpellExtras: SpellFlagMage,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.01*float64(mage.Talents.Pyromaniac)) *
					(1 - 0.01*float64(mage.Talents.ElementalPrecision)),

				GCD: core.GCDDefault,
			},
			Cooldown: time.Second*8 - time.Millisecond*500*time.Duration(mage.Talents.ImprovedFireBlast),
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BonusSpellHitRating: float64(mage.Talents.ElementalPrecision) * 1 * core.SpellHitRatingPerHitChance,

			BonusSpellCritRating: 0 +
				float64(mage.Talents.CriticalMass)*2*core.SpellCritRatingPerCritChance +
				float64(mage.Talents.Pyromaniac)*1*core.SpellCritRatingPerCritChance,

			DamageMultiplier: mage.spellDamageMultiplier * (1 + 0.02*float64(mage.Talents.FirePower)),

			ThreatMultiplier: 1 - 0.05*float64(mage.Talents.BurningSoul),

			BaseDamage:     core.BaseDamageConfigMagic(664, 786, 1.5/3.5),
			OutcomeApplier: core.OutcomeFuncMagicHitAndCrit(mage.SpellCritMultiplier(1, 0.25*float64(mage.Talents.SpellPower))),
		}),
	})
}
