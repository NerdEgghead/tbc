// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
export const RogueRotationConfig = {
    inputs: [
        {
            type: 'boolean', cssClass: 'maintain-expose-armor-picker',
            getModObject: (simUI) => simUI.player,
            config: {
                label: 'Maintain EA',
                labelTooltip: 'Keeps Expose Armor active on the primary target.',
                changedEvent: (player) => player.rotationChangeEmitter,
                getValue: (player) => player.getRotation().maintainExposeArmor,
                setValue: (eventID, player, newValue) => {
                    const newRotation = player.getRotation();
                    newRotation.maintainExposeArmor = newValue;
                    player.setRotation(eventID, newRotation);
                },
            },
        },
        {
            type: 'boolean', cssClass: 'use-rupture-picker',
            getModObject: (simUI) => simUI.player,
            config: {
                label: 'Use Rupture',
                labelTooltip: 'Uses Rupture over Eviscerate when appropriate.',
                changedEvent: (player) => player.rotationChangeEmitter,
                getValue: (player) => player.getRotation().useRupture,
                setValue: (eventID, player, newValue) => {
                    const newRotation = player.getRotation();
                    newRotation.useRupture = newValue;
                    player.setRotation(eventID, newRotation);
                },
            },
        },
    ],
};
