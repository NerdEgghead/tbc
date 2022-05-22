package druid

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

var RipActionID = core.ActionID{SpellID: 27008}
var RipEnergyCost = 30.0

func (druid *Druid) registerRipSpell(sim *core.Simulation) {
	druid.Rip = druid.RegisterSpell(core.SpellConfig{
		ActionID:    RipActionID,
		SpellSchool: core.SpellSchoolPhysical,
		SpellExtras: core.SpellExtrasMeleeMetrics | core.SpellExtrasIgnoreResists,

		ResourceType: stats.Energy,
		BaseCost:     RipEnergyCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: RipEnergyCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			OutcomeApplier:   core.OutcomeFuncMeleeSpecialHit(),
			OnSpellHit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.RipDot.Apply(sim)
					druid.SpendComboPoints(sim, spell.ActionID)
				}
			},
		}),
	})

	target := sim.GetPrimaryTarget()
	druid.RipDot = core.NewDot(core.Dot{
		Spell: druid.Rip,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Rip-" + strconv.Itoa(int(druid.Index)),
			ActionID: RipActionID,
		}),
		NumberOfTicks: 6,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			DamageMultiplier: 1 + core.TernaryFloat64(ItemSetThunderheartFeral.CharacterHasSetBonus(&druid.Character, 4), 0.15, 0),
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage: core.BuildBaseDamageConfig(func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				comboPoints := druid.ComboPoints()
				attackPower := hitEffect.MeleeAttackPower(spell.Character)

				if comboPoints < 3 {
					panic("Only 3-5 CP Rips are supported at present.")
				}
				if comboPoints == 3 {
					return (990 + 0.18*attackPower) / 6
				}
				if comboPoints == 4 {
					return (1272 + 0.24*attackPower) / 6
				}
				if comboPoints == 5 {
					return (1554 + 0.24*attackPower) / 6
				}
			}, 0),
			OutcomeApplier: core.OutcomeFuncTick(),
		}),
	})
}

func (druid *Druid) ShouldCastRip(sim *core.Simulation, rotation proto.FeralDruid_Rotation) bool {
	energy:= druid.CurrentEnergy()
	comboPoints := druid.ComboPoints()
	canPrimaryRip := (rotation.FinishingMove == proto.FeralDruid_Rotation_Rip)
	canWeaveRip := (rotation.FinishingMove == proto.FeralDruid_Rotation_Bite) && rotation.Ripweave && (energy >= 52) && !druid.PseudoStats.NoCost
	nearFightEnd := (sim.GetRemainingDuration() < time.Duration(10.0 * float64(time.Second)))
	return (canPrimaryRip || canWeaveRip) && (comboPoints >= rotation.RipCp) && !druid.RipDot.IsActive() && !nearFightEnd
}
