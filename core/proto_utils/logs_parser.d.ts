import { RaidSimResult } from '/tbc/core/proto/api.js';
export declare class Entity {
    readonly name: string;
    readonly ownerName: string;
    readonly index: number;
    readonly isTarget: boolean;
    readonly isPet: boolean;
    constructor(name: string, ownerName: string, index: number, isTarget: boolean, isPet: boolean);
    static parseRegex: RegExp;
    static parseAll(str: string): Array<Entity>;
}
interface SimLogParams {
    raw: string;
    timestamp: number;
    source: Entity | null;
    target: Entity | null;
}
export declare class SimLog {
    readonly raw: string;
    readonly timestamp: number;
    readonly source: Entity | null;
    readonly target: Entity | null;
    constructor(params: SimLogParams);
    static parseAll(result: RaidSimResult): Array<SimLog>;
}
export declare class AuraGainedLog extends SimLog {
    readonly auraName: string;
    constructor(params: SimLogParams, auraName: string);
    static parse(params: SimLogParams): AuraGainedLog | null;
}
export declare class AuraFadedLog extends SimLog {
    readonly auraName: string;
    constructor(params: SimLogParams, auraName: string);
    static parse(params: SimLogParams): AuraFadedLog | null;
}
export declare class ManaChangedLog extends SimLog {
    readonly manaBefore: number;
    readonly manaAfter: number;
    readonly cause: string;
    constructor(params: SimLogParams, manaBefore: number, manaAfter: number, cause: string);
    static parse(params: SimLogParams): ManaChangedLog | null;
}
export {};
