package druid

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var PowershiftActionID = core.ActionID{SpellID: 768}

func (druid *Druid) registerPowershiftSpell(sim *core.Simulation) {
	baseCost := 830.0
	finalEnergy := 40.0

	if druid.Wolfshead {
		finalEnergy += 20.0
	}

	druid.Powershift = druid.RegisterSpell(core.SpellConfig{
		ActionID: PowershiftActionID,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.1*float64(druid.Talents.NaturalShapeshifter)),
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Target, spell *core.Spell) {
			currentEnergy := druid.CurrentEnergy()
			druid.AddEnergy(sim, finalEnergy - currentEnergy, PowershiftActionID)
			druid.CatForm = true
		},
	})
}
