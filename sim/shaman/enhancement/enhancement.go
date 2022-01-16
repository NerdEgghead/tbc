package enhancement

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/shaman"
)

func RegisterEnhancementShaman() {
	core.RegisterAgentFactory(
		proto.Player_EnhancementShaman{},
		func(character core.Character, options proto.Player) core.Agent {
			return NewEnhancementShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_EnhancementShaman)
			if !ok {
				panic("Invalid spec value for Enhancement Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewEnhancementShaman(character core.Character, options proto.Player) *EnhancementShaman {
	enhOptions := options.GetEnhancementShaman()

	selfBuffs := shaman.SelfBuffs{
		Bloodlust:   enhOptions.Options.Bloodlust,
		WaterShield: enhOptions.Options.WaterShield,
	}

	if enhOptions.Rotation.Totems != nil {
		selfBuffs.ManaSpring = enhOptions.Rotation.Totems.Water == proto.WaterTotem_ManaSpringTotem
		selfBuffs.EarthTotem = enhOptions.Rotation.Totems.Earth
		selfBuffs.AirTotem = enhOptions.Rotation.Totems.Air
		selfBuffs.NextTotemDropType[shaman.AirTotem] = int32(enhOptions.Rotation.Totems.Air)
		selfBuffs.FireTotem = enhOptions.Rotation.Totems.Fire
		selfBuffs.NextTotemDropType[shaman.FireTotem] = int32(enhOptions.Rotation.Totems.Fire)

		if enhOptions.Rotation.Totems.WindfuryTotemRank == 0 {
			// If rank is 0, disable windfury options.
			enhOptions.Rotation.Totems.TwistWindfury = false
			if enhOptions.Rotation.Totems.Air == proto.AirTotem_WindfuryTotem {
				enhOptions.Rotation.Totems.Air = proto.AirTotem_NoAirTotem
			}
		}
		if enhOptions.Rotation.Totems.Air == proto.AirTotem_WindfuryTotem {
			// No need to twist windfury if its already the default totem.
			enhOptions.Rotation.Totems.TwistWindfury = false
		} else if enhOptions.Rotation.Totems.Air == proto.AirTotem_NoAirTotem && enhOptions.Rotation.Totems.TwistWindfury {
			// If twisting windfury without a default air totem, make windfury the default instead.
			enhOptions.Rotation.Totems.Air = proto.AirTotem_WindfuryTotem
			enhOptions.Rotation.Totems.TwistWindfury = false
		}

		selfBuffs.TwistWindfury = enhOptions.Rotation.Totems.TwistWindfury
		selfBuffs.WindfuryTotemRank = enhOptions.Rotation.Totems.WindfuryTotemRank
		if selfBuffs.TwistWindfury {
			selfBuffs.NextTotemDropType[shaman.AirTotem] = int32(proto.AirTotem_WindfuryTotem)
			selfBuffs.NextTotemDrops[shaman.AirTotem] = 0 // drop windfury immediately
		}

		selfBuffs.TwistFireNova = enhOptions.Rotation.Totems.TwistFireNova
		if selfBuffs.TwistFireNova {
			selfBuffs.NextTotemDropType[shaman.FireTotem] = int32(proto.FireTotem_FireNovaTotem) // start by dropping nova, then alternating.
		}
	}
	enh := &EnhancementShaman{
		Shaman:   shaman.NewShaman(character, *enhOptions.Talents, selfBuffs),
		Rotation: *enhOptions.Rotation,
	}
	// Enable Auto Attacks for this spec
	enh.EnableAutoAttacks(enhOptions.Options.DelayOffhandSwings)

	// Modify auto attacks multiplier from weapon mastery.
	enh.AutoAttacks.Effect.DamageMultiplier *= 1 + 0.02*float64(enhOptions.Talents.WeaponMastery)
	enh.ApplyWindfuryImbue(
		enhOptions.Options.MainHandImbue == proto.ShamanWeaponImbue_ImbueWindfury,
		enhOptions.Options.OffHandImbue == proto.ShamanWeaponImbue_ImbueWindfury)
	enh.ApplyFlametongueImbue(
		enhOptions.Options.MainHandImbue == proto.ShamanWeaponImbue_ImbueFlametongue,
		enhOptions.Options.OffHandImbue == proto.ShamanWeaponImbue_ImbueFlametongue)
	enh.ApplyFrostbrandImbue(
		enhOptions.Options.MainHandImbue == proto.ShamanWeaponImbue_ImbueFrostbrand,
		enhOptions.Options.OffHandImbue == proto.ShamanWeaponImbue_ImbueFrostbrand)
	enh.ApplyRockbiterImbue(
		enhOptions.Options.MainHandImbue == proto.ShamanWeaponImbue_ImbueRockbiter,
		enhOptions.Options.OffHandImbue == proto.ShamanWeaponImbue_ImbueRockbiter)

	if enhOptions.Options.MainHandImbue != proto.ShamanWeaponImbue_ImbueNone {
		enh.HasMHWeaponImbue = true
	}

	return enh
}

type EnhancementShaman struct {
	*shaman.Shaman

	Rotation proto.EnhancementShaman_Rotation
}

func (enh *EnhancementShaman) GetShaman() *shaman.Shaman {
	return enh.Shaman
}

func (enh *EnhancementShaman) Reset(sim *core.Simulation) {
	enh.Shaman.Reset(sim)
}

func (enh *EnhancementShaman) Act(sim *core.Simulation) time.Duration {
	// Redrop totems when needed.
	dropTime, dropSuccess := enh.TryDropTotems(sim)
	if dropTime > 0 {
		nextEventTime := core.MinDuration(dropTime, enh.AutoAttacks.NextAttackAt())
		if !dropSuccess {
			enh.Metrics.MarkOOM(sim, &enh.Character, nextEventTime-sim.CurrentTime)
		}
		return nextEventTime
	}

	target := sim.GetPrimaryTarget()

	success := true
	cost := 0.0
	if !enh.IsOnCD(shaman.StormstrikeCD, sim.CurrentTime) {
		ss := enh.NewStormstrike(sim, target)
		cost = ss.Cost.Value
		if success = ss.Attack(sim); success {
			return enh.AutoAttacks.NextEventAt(sim)
		}
	} else if !enh.IsOnCD(shaman.ShockCooldownID, sim.CurrentTime) {
		var shock *core.SimpleSpell
		if enh.Rotation.WeaveFlameShock && !enh.FlameShockSpell.IsInUse() {
			shock = enh.NewFlameShock(sim, target)
		} else if enh.Rotation.PrimaryShock == proto.EnhancementShaman_Rotation_Earth {
			shock = enh.NewEarthShock(sim, target)
		} else if enh.Rotation.PrimaryShock == proto.EnhancementShaman_Rotation_Frost {
			shock = enh.NewFrostShock(sim, target)
		}

		if shock != nil {
			cost = shock.ManaCost
			if success = shock.Cast(sim); success {
				return enh.AutoAttacks.NextEventAt(sim)
			}
		}
	}
	if !success {
		regenTime := enh.TimeUntilManaRegen(cost)
		nextActionAt := core.MinDuration(sim.CurrentTime+regenTime, enh.AutoAttacks.NextAttackAt())
		enh.Metrics.MarkOOM(sim, &enh.Character, nextActionAt-sim.CurrentTime)
		return nextActionAt
	}

	// We didn't try to cast anything. Just wait for next auto.
	return enh.AutoAttacks.NextAttackAt()
}
