package hunter

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var ScorpidStingAuraID = core.NewAuraID()

func (hunter *Hunter) registerScorpidStingSpell(sim *core.Simulation) {
	actionID := core.ActionID{SpellID: 3043}
	cost := core.ResourceCost{Type: stats.Mana, Value: hunter.BaseMana() * 0.09}
	ama := core.SimpleSpell{
		SpellCast: core.SpellCast{
			Cast: core.Cast{
				ActionID:    actionID,
				Character:   &hunter.Character,
				SpellSchool: core.SpellSchoolNature,
				GCD:         core.GCDDefault,
				Cost:        cost,
				BaseCost:    cost,
				IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			},
		},
		Effect: core.SpellEffect{
			OutcomeRollCategory: core.OutcomeRollCategoryRanged,
			CritRollCategory:    core.CritRollCategoryPhysical,
			ProcMask:            core.ProcMaskRangedSpecial,
			OnSpellHit: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}

				spellEffect.Target.AddAura(sim, core.Aura{
					ID:       ScorpidStingAuraID,
					ActionID: actionID,
					Duration: time.Second * 20,
				})
			},
		},
	}

	ama.Cost.Value *= 1 - 0.02*float64(hunter.Talents.Efficiency)

	hunter.ScorpidSting = hunter.RegisterSpell(core.SpellConfig{
		Template:   ama,
		ModifyCast: core.ModifyCastAssignTarget,
	})
}
