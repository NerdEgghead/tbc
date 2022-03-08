import { Stat } from '/tbc/core/proto/common.js';
import { statNames } from '/tbc/core/proto_utils/names.js';
import { TypedEvent } from '/tbc/core/typed_event.js';
import { NumberPicker } from '/tbc/core/components/number_picker.js';
import { getEnumValues } from '/tbc/core/utils.js';
import { Popup } from './popup.js';
export class SettingsMenu extends Popup {
    constructor(parent, simUI) {
        super(parent);
        this.rootElem.classList.add('settings-menu');
        this.simUI = simUI;
        this.rootElem.innerHTML = `
			<div class="settings-menu-title">
				<span>SETTINGS</span>
			</div>
			<div class="settings-menu-content">
				<div class="settings-menu-content-left">
					<button class="restore-defaults-button sim-button">RESTORE DEFAULTS</button>
					<div class="settings-menu-section">
						<div class="fixed-rng-seed">
						</div>
						<div class="last-used-rng-seed-container">
							<span>Last used RNG seed:</span><span class="last-used-rng-seed">0</span>
						</div>
					</div>
				</div>
				<div class="settings-menu-content-right">
					<div class="settings-menu-section settings-menu-ep-weights">
					</div>
				</div>
			</div>
		`;
        this.addCloseButton();
        const restoreDefaultsButton = this.rootElem.getElementsByClassName('restore-defaults-button')[0];
        restoreDefaultsButton.addEventListener('click', event => {
            this.simUI.applyDefaults(TypedEvent.nextEventID());
        });
        tippy(restoreDefaultsButton, {
            'content': 'Restores all default settings (gear, consumes, buffs, talents, EP weights, etc).',
            'allowHTML': true,
        });
        const fixedRngSeed = this.rootElem.getElementsByClassName('fixed-rng-seed')[0];
        new NumberPicker(fixedRngSeed, this.simUI.sim, {
            label: 'Fixed RNG Seed',
            labelTooltip: 'Seed value for the random number generator used during sims, or 0 to use different randomness each run. Use this to share exact sim results or for debugging.',
            changedEvent: (sim) => sim.fixedRngSeedChangeEmitter,
            getValue: (sim) => sim.getFixedRngSeed(),
            setValue: (eventID, sim, newValue) => {
                sim.setFixedRngSeed(eventID, newValue);
            },
        });
        const lastUsedRngSeed = this.rootElem.getElementsByClassName('last-used-rng-seed')[0];
        lastUsedRngSeed.textContent = String(this.simUI.sim.getLastUsedRngSeed());
        this.simUI.sim.lastUsedRngSeedChangeEmitter.on(() => lastUsedRngSeed.textContent = String(this.simUI.sim.getLastUsedRngSeed()));
        this.setupEpWeightsSettings();
    }
    setupEpWeightsSettings() {
        const sectionRoot = this.rootElem.getElementsByClassName('settings-menu-ep-weights')[0];
        const label = document.createElement('span');
        label.classList.add('ep-weights-label');
        label.textContent = 'EP Weights';
        tippy(label, {
            'content': 'EP Weights for sorting the item selector.',
            'allowHTML': true,
        });
        sectionRoot.appendChild(label);
        //const epStats = this.simUI.individualConfig.epStats;
        const epStats = getEnumValues(Stat).filter(stat => ![Stat.StatMana, Stat.StatEnergy, Stat.StatRage].includes(stat));
        const weightPickers = epStats.map(stat => new NumberPicker(sectionRoot, this.simUI.player, {
            label: statNames[stat],
            changedEvent: (player) => player.epWeightsChangeEmitter,
            getValue: (player) => player.getEpWeights().getStat(stat),
            setValue: (eventID, player, newValue) => {
                const epWeights = player.getEpWeights().withStat(stat, newValue);
                player.setEpWeights(eventID, epWeights);
            },
        }));
    }
}
