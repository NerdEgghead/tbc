package druid

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

type Druid struct {
	core.Character
	SelfBuffs
	Talents proto.DruidTalents

	RebirthUsed bool
	CatForm     bool
	Wolfshead   bool

	FaerieFire  *core.Spell
	Hurricane   *core.Spell
	InsectSwarm *core.Spell
	Moonfire    *core.Spell
	Rebirth     *core.Spell
	Starfire6   *core.Spell
	Starfire8   *core.Spell
	Wrath       *core.Spell
	Powershift  *core.Spell

	InsectSwarmDot *core.Dot
	MoonfireDot    *core.Dot

	FaerieFireAura       *core.Aura
	NaturesGraceProcAura *core.Aura
	NaturesSwiftnessAura *core.Aura
}

type SelfBuffs struct {
	Omen bool

	InnervateTarget proto.RaidTarget
}

func (druid *Druid) GetCharacter() *core.Character {
	return &druid.Character
}

func (druid *Druid) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.GiftOfTheWild = core.MaxTristate(raidBuffs.GiftOfTheWild, proto.TristateEffect_TristateEffectRegular)
	if druid.Talents.ImprovedMarkOfTheWild == 5 { // probably could work on actually calculating the fraction effect later if we care.
		raidBuffs.GiftOfTheWild = proto.TristateEffect_TristateEffectImproved
	}
}

const ravenGoddessItemID = 32387
const wolfsheadItemID = 8345

func (druid *Druid) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if druid.Talents.MoonkinForm { // assume if you have moonkin talent you are using it.
		partyBuffs.MoonkinAura = core.MaxTristate(partyBuffs.MoonkinAura, proto.TristateEffect_TristateEffectRegular)
		for _, e := range druid.Equip {
			if e.ID == ravenGoddessItemID {
				partyBuffs.MoonkinAura = proto.TristateEffect_TristateEffectImproved
				break
			}
		}
	}
	if druid.Talents.LeaderOfThePack {
		partyBuffs.LeaderOfThePack = core.MaxTristate(partyBuffs.LeaderOfThePack, proto.TristateEffect_TristateEffectRegular)
		for _, e := range druid.Equip {
			if e.ID == ravenGoddessItemID {
				partyBuffs.LeaderOfThePack = proto.TristateEffect_TristateEffectImproved
				break
			}
		}
	}
}

func (druid *Druid) Init(sim *core.Simulation) {
	if druid.CatForm {
		druid.registerPowershiftSpell(sim)
		druid.registerRipSpell(sim)
	} else {
		druid.registerFaerieFireSpell(sim)
		druid.registerHurricaneSpell(sim)
		druid.registerInsectSwarmSpell(sim)
		druid.registerMoonfireSpell(sim)
		druid.registerRebirthSpell(sim)
		druid.Starfire8 = druid.newStarfireSpell(sim, 8)
		druid.Starfire6 = druid.newStarfireSpell(sim, 6)
		druid.registerWrathSpell(sim)
	}

	if druid.CatForm {
		druid.DelayCooldownsForArmorDebuffs(sim)
	}
}

func (druid *Druid) Reset(sim *core.Simulation) {
	druid.RebirthUsed = false
}

func (druid *Druid) Act(sim *core.Simulation) time.Duration {
	return core.NeverExpires // does nothing
}

func New(char core.Character, selfBuffs SelfBuffs, talents proto.DruidTalents) *Druid {
	druid := &Druid{
		Character:   char,
		SelfBuffs:   selfBuffs,
		Talents:     talents,
		RebirthUsed: false,
		CatForm:     false,
		Wolfshead:   false,
	}
	druid.EnableManaBar()

	druid.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Intellect,
		ModifiedStat: stats.SpellCrit,
		Modifier: func(intellect float64, spellCrit float64) float64 {
			return spellCrit + (intellect/79.4)*core.SpellCritRatingPerCritChance
		},
	})

	druid.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Strength,
		ModifiedStat: stats.AttackPower,
		Modifier: func(strength float64, attackPower float64) float64 {
			return attackPower + strength*2
		},
	})

	druid.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.MeleeCrit,
		Modifier: func(agility float64, meleeCrit float64) float64 {
			return meleeCrit + (agility/25)*core.MeleeCritRatingPerCritChance
		},
	})

	// Check if Wolfshead Helm is equipped
	for _, e := range druid.Equip {
		if e.ID == wolfsheadItemID {
			druid.Wolfshead = true
			break
		}
	}

	return druid
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassDruid}] = stats.Stats{
		stats.Health:      3434, // 4498 health shown on naked character (would include tauren bonus)
		stats.Strength:    81,
		stats.Agility:     65,
		stats.Stamina:     85,
		stats.Intellect:   115,
		stats.Spirit:      135,
		stats.Mana:        2370,
		stats.SpellCrit:   40.66,                                    // 3.29% chance to crit shown on naked character screen
		stats.AttackPower: -20,                                      // accounts for the fact that the first 20 points in Str only provide 1 AP rather than 2
		stats.MeleeCrit:   0.96 * core.MeleeCritRatingPerCritChance, // 3.56% chance to crit shown on naked character screen
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassDruid}] = stats.Stats{
		stats.Health:      3434, // 4254 health shown on naked character
		stats.Strength:    73,
		stats.Agility:     75,
		stats.Stamina:     82,
		stats.Intellect:   120,
		stats.Spirit:      133,
		stats.Mana:        2370,
		stats.SpellCrit:   40.60,                                    // 3.35% chance to crit shown on naked character screen
		stats.AttackPower: -20,                                      // accounts for the fact that the first 20 points in Str only provide 1 AP rather than 2
		stats.MeleeCrit:   0.96 * core.MeleeCritRatingPerCritChance, // 3.96% chance to crit shown on naked character screen
	}
}

// Agent is a generic way to access underlying druid on any of the agents (for example balance druid.)
type Agent interface {
	GetDruid() *Druid
}
