package druid

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func init() {
	core.AddItemSet(ItemSetMalorne)
}

var Malorne2PcAuraID = core.NewAuraID()

var ItemSetMalorne = core.ItemSet{
	Name:  "Malorne Rainment",
	Items: map[int32]struct{}{29093: {}, 29094: {}, 29091: {}, 29092: {}, 29095: {}},
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddPermanentAura(func(sim *core.Simulation) core.Aura {
				return core.Aura{
					ID:   Malorne2PcAuraID,
					Name: "Malorne 2pc Bonus",
					OnSpellHit: func(sim *core.Simulation, spellCast *core.SpellCast, spellEffect *core.SpellEffect) {
						if sim.RandomFloat("malorne 2p") < 0.05 {
							spellCast.Character.AddStat(stats.Mana, 120)
						}
					},
				}
			})
		},
		4: func(agent core.Agent) {
			// Currently this is handled in druid.go (reducing CD of innervate)
		},
	},
}

var ItemSetNordrassil = core.ItemSet{
	Name:  "Nordrassil Regalia",
	Items: map[int32]struct{}{30231: {}, 30232: {}, 30233: {}, 30234: {}, 30235: {}},
	Bonuses: map[int32]core.ApplyEffect{
		4: func(agent core.Agent) {
			// handled in druid.go on spell hit
		},
	},
}

var ItemSetThunderheart = core.ItemSet{
	Name:  "Thunderheart Regalia",
	Items: map[int32]struct{}{31043: {}, 31035: {}, 31040: {}, 31046: {}, 31049: {}, 34572: {}, 34446: {}, 34555: {}},
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// handled in moonfire.go in template construction
		},
		4: func(agent core.Agent) {
			// handled in starfire.go in template construction
		},
	},
}
