import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
import { Raid } from "./api";
import { RaidTarget } from "./common";
import { Cooldowns } from "./common";
import { Race } from "./common";
import { Consumes } from "./common";
import { IndividualBuffs } from "./common";
import { EquipmentSpec } from "./common";
import { Encounter } from "./common";
import { Player } from "./api";
import { PartyBuffs } from "./common";
import { RaidBuffs } from "./common";
/**
 * @generated from protobuf message proto.SimSettings
 */
export interface SimSettings {
    /**
     * @generated from protobuf field: int32 iterations = 1;
     */
    iterations: number;
    /**
     * @generated from protobuf field: int32 phase = 2;
     */
    phase: number;
    /**
     * @generated from protobuf field: int64 fixed_rng_seed = 3;
     */
    fixedRngSeed: bigint;
}
/**
 * Contains all information that is imported/exported from an individual sim.
 *
 * @generated from protobuf message proto.IndividualSimSettings
 */
export interface IndividualSimSettings {
    /**
     * @generated from protobuf field: proto.SimSettings settings = 5;
     */
    settings?: SimSettings;
    /**
     * @generated from protobuf field: proto.RaidBuffs raid_buffs = 1;
     */
    raidBuffs?: RaidBuffs;
    /**
     * @generated from protobuf field: proto.PartyBuffs party_buffs = 2;
     */
    partyBuffs?: PartyBuffs;
    /**
     * @generated from protobuf field: proto.Player player = 3;
     */
    player?: Player;
    /**
     * @generated from protobuf field: proto.Encounter encounter = 4;
     */
    encounter?: Encounter;
    /**
     * @generated from protobuf field: repeated double ep_weights = 6;
     */
    epWeights: number[];
}
/**
 * Local storage data for gear settings.
 *
 * @generated from protobuf message proto.SavedGearSet
 */
export interface SavedGearSet {
    /**
     * @generated from protobuf field: proto.EquipmentSpec gear = 1;
     */
    gear?: EquipmentSpec;
    /**
     * @generated from protobuf field: repeated double bonus_stats = 2;
     */
    bonusStats: number[];
}
/**
 * Local storage data for other settings.
 *
 * @generated from protobuf message proto.SavedSettings
 */
export interface SavedSettings {
    /**
     * @generated from protobuf field: proto.RaidBuffs raid_buffs = 1;
     */
    raidBuffs?: RaidBuffs;
    /**
     * @generated from protobuf field: proto.PartyBuffs party_buffs = 2;
     */
    partyBuffs?: PartyBuffs;
    /**
     * @generated from protobuf field: proto.IndividualBuffs player_buffs = 3;
     */
    playerBuffs?: IndividualBuffs;
    /**
     * @generated from protobuf field: proto.Consumes consumes = 4;
     */
    consumes?: Consumes;
    /**
     * @generated from protobuf field: proto.Race race = 5;
     */
    race: Race;
    /**
     * @generated from protobuf field: proto.Cooldowns cooldowns = 6;
     */
    cooldowns?: Cooldowns;
}
/**
 * @generated from protobuf message proto.SavedTalents
 */
export interface SavedTalents {
    /**
     * @generated from protobuf field: string talents_string = 1;
     */
    talentsString: string;
}
/**
 * A buff bot placed in a raid.
 *
 * @generated from protobuf message proto.BuffBot
 */
export interface BuffBot {
    /**
     * Uniquely identifies which buffbot this is.
     *
     * @generated from protobuf field: string id = 1;
     */
    id: string;
    /**
     * @generated from protobuf field: int32 raid_index = 2;
     */
    raidIndex: number;
    /**
     * The assigned player to innervate. Only used for druid buffbots.
     *
     * @generated from protobuf field: proto.RaidTarget innervate_assignment = 3;
     */
    innervateAssignment?: RaidTarget;
    /**
     * The assigned player to PI. Only used for disc priest buffbots.
     *
     * @generated from protobuf field: proto.RaidTarget power_infusion_assignment = 4;
     */
    powerInfusionAssignment?: RaidTarget;
}
/**
 * @generated from protobuf message proto.BlessingsAssignment
 */
export interface BlessingsAssignment {
    /**
     * Index corresponds to Spec that the blessing should be applied to.
     *
     * @generated from protobuf field: repeated proto.Blessings blessings = 1;
     */
    blessings: Blessings[];
}
/**
 * @generated from protobuf message proto.BlessingsAssignments
 */
export interface BlessingsAssignments {
    /**
     * Assignments for each paladin.
     *
     * @generated from protobuf field: repeated proto.BlessingsAssignment paladins = 1;
     */
    paladins: BlessingsAssignment[];
}
/**
 * Local storage data for a saved encounter.
 *
 * @generated from protobuf message proto.SavedEncounter
 */
export interface SavedEncounter {
    /**
     * @generated from protobuf field: proto.Encounter encounter = 1;
     */
    encounter?: Encounter;
}
/**
 * Local storage data for raid sim settings.
 *
 * @generated from protobuf message proto.SavedRaid
 */
export interface SavedRaid {
    /**
     * @generated from protobuf field: proto.Raid raid = 1;
     */
    raid?: Raid;
    /**
     * @generated from protobuf field: repeated proto.BuffBot buff_bots = 2;
     */
    buffBots: BuffBot[];
    /**
     * @generated from protobuf field: proto.BlessingsAssignments blessings = 3;
     */
    blessings?: BlessingsAssignments;
}
/**
 * Contains all information that is imported/exported from a raid sim.
 *
 * @generated from protobuf message proto.RaidSimSettings
 */
export interface RaidSimSettings {
    /**
     * @generated from protobuf field: proto.SimSettings settings = 5;
     */
    settings?: SimSettings;
    /**
     * @generated from protobuf field: proto.Raid raid = 1;
     */
    raid?: Raid;
    /**
     * @generated from protobuf field: repeated proto.BuffBot buff_bots = 2;
     */
    buffBots: BuffBot[];
    /**
     * @generated from protobuf field: proto.BlessingsAssignments blessings = 3;
     */
    blessings?: BlessingsAssignments;
    /**
     * @generated from protobuf field: proto.Encounter encounter = 4;
     */
    encounter?: Encounter;
}
/**
 * @generated from protobuf enum proto.Blessings
 */
export declare enum Blessings {
    /**
     * @generated from protobuf enum value: BlessingUnknown = 0;
     */
    BlessingUnknown = 0,
    /**
     * @generated from protobuf enum value: BlessingOfKings = 1;
     */
    BlessingOfKings = 1,
    /**
     * @generated from protobuf enum value: BlessingOfMight = 2;
     */
    BlessingOfMight = 2,
    /**
     * @generated from protobuf enum value: BlessingOfSalvation = 3;
     */
    BlessingOfSalvation = 3,
    /**
     * @generated from protobuf enum value: BlessingOfWisdom = 4;
     */
    BlessingOfWisdom = 4
}
declare class SimSettings$Type extends MessageType<SimSettings> {
    constructor();
    create(value?: PartialMessage<SimSettings>): SimSettings;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SimSettings): SimSettings;
    internalBinaryWrite(message: SimSettings, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.SimSettings
 */
export declare const SimSettings: SimSettings$Type;
declare class IndividualSimSettings$Type extends MessageType<IndividualSimSettings> {
    constructor();
    create(value?: PartialMessage<IndividualSimSettings>): IndividualSimSettings;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: IndividualSimSettings): IndividualSimSettings;
    internalBinaryWrite(message: IndividualSimSettings, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.IndividualSimSettings
 */
export declare const IndividualSimSettings: IndividualSimSettings$Type;
declare class SavedGearSet$Type extends MessageType<SavedGearSet> {
    constructor();
    create(value?: PartialMessage<SavedGearSet>): SavedGearSet;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SavedGearSet): SavedGearSet;
    internalBinaryWrite(message: SavedGearSet, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.SavedGearSet
 */
export declare const SavedGearSet: SavedGearSet$Type;
declare class SavedSettings$Type extends MessageType<SavedSettings> {
    constructor();
    create(value?: PartialMessage<SavedSettings>): SavedSettings;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SavedSettings): SavedSettings;
    internalBinaryWrite(message: SavedSettings, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.SavedSettings
 */
export declare const SavedSettings: SavedSettings$Type;
declare class SavedTalents$Type extends MessageType<SavedTalents> {
    constructor();
    create(value?: PartialMessage<SavedTalents>): SavedTalents;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SavedTalents): SavedTalents;
    internalBinaryWrite(message: SavedTalents, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.SavedTalents
 */
export declare const SavedTalents: SavedTalents$Type;
declare class BuffBot$Type extends MessageType<BuffBot> {
    constructor();
    create(value?: PartialMessage<BuffBot>): BuffBot;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: BuffBot): BuffBot;
    internalBinaryWrite(message: BuffBot, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.BuffBot
 */
export declare const BuffBot: BuffBot$Type;
declare class BlessingsAssignment$Type extends MessageType<BlessingsAssignment> {
    constructor();
    create(value?: PartialMessage<BlessingsAssignment>): BlessingsAssignment;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: BlessingsAssignment): BlessingsAssignment;
    internalBinaryWrite(message: BlessingsAssignment, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.BlessingsAssignment
 */
export declare const BlessingsAssignment: BlessingsAssignment$Type;
declare class BlessingsAssignments$Type extends MessageType<BlessingsAssignments> {
    constructor();
    create(value?: PartialMessage<BlessingsAssignments>): BlessingsAssignments;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: BlessingsAssignments): BlessingsAssignments;
    internalBinaryWrite(message: BlessingsAssignments, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.BlessingsAssignments
 */
export declare const BlessingsAssignments: BlessingsAssignments$Type;
declare class SavedEncounter$Type extends MessageType<SavedEncounter> {
    constructor();
    create(value?: PartialMessage<SavedEncounter>): SavedEncounter;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SavedEncounter): SavedEncounter;
    internalBinaryWrite(message: SavedEncounter, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.SavedEncounter
 */
export declare const SavedEncounter: SavedEncounter$Type;
declare class SavedRaid$Type extends MessageType<SavedRaid> {
    constructor();
    create(value?: PartialMessage<SavedRaid>): SavedRaid;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SavedRaid): SavedRaid;
    internalBinaryWrite(message: SavedRaid, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.SavedRaid
 */
export declare const SavedRaid: SavedRaid$Type;
declare class RaidSimSettings$Type extends MessageType<RaidSimSettings> {
    constructor();
    create(value?: PartialMessage<RaidSimSettings>): RaidSimSettings;
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: RaidSimSettings): RaidSimSettings;
    internalBinaryWrite(message: RaidSimSettings, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter;
}
/**
 * @generated MessageType for protobuf message proto.RaidSimSettings
 */
export declare const RaidSimSettings: RaidSimSettings$Type;
export {};
