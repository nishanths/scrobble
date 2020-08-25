import { assertExhaustive, hexEncode } from "../../shared/util"
import { Color } from "../colorpicker"

export enum DetailKind {
	Song, Album
}

export enum Mode {
	All, Loved, Color, Insights
}

export const controlValues = ["All", "Loved", "By color", "Insights"] as const

export type ControlValue = typeof controlValues[number]

export const controlValueForMode = (m: Mode): ControlValue => {
	switch (m) {
		case Mode.All: return "All"
		case Mode.Loved: return "Loved"
		case Mode.Color: return "By color"
		case Mode.Insights: return "Insights"
		default: assertExhaustive(m)
	}
}

export const modeFromControlValue = (v: ControlValue): Mode => {
	switch (v) {
		case "All": return Mode.All
		case "Loved": return Mode.Loved
		case "By color": return Mode.Color
		case "Insights": return Mode.Insights
		default: assertExhaustive(v)
	}
}

export const pathForMode = (m: Mode): string => {
	switch (m) {
		case Mode.All: return ""
		case Mode.Loved: return "/loved"
		case Mode.Color: return "/color"
		case Mode.Insights: return "/insights"
	}
	assertExhaustive(m)
}

export const pathForColor = (c: Color | undefined): string => {
	if (c === undefined) {
		return ""
	}
	return "/" + c
}

export const pathForDetailKind = (k: DetailKind): string => {
	switch (k) {
		case DetailKind.Song:
			return "/song"
		case DetailKind.Album:
			return "/album"
		default:
			assertExhaustive(k)
	}
}

export const fullPath = (
	profileUsername: string,
	mode: Mode,
	color: Color | undefined,
	insightType: InsightType | undefined,
	detail: {
		kind: DetailKind
		songIdent: string
	} | undefined,
): string => {
	let u: string;
	switch (mode) {
		case Mode.All:
		case Mode.Loved:
			u = "/u/" + profileUsername + pathForMode(mode)
			break
		case Mode.Insights:
			u = "/u/" + profileUsername + pathForMode(mode) + pathForInsightType(insightType)
			break
		case Mode.Color:
			u = "/u/" + profileUsername + pathForMode(mode) + pathForColor(color)
			break
		default:
			assertExhaustive(mode)
	}
	if (detail !== undefined) {
		u += pathForDetailKind(detail.kind) + "/" + hexEncode(detail.songIdent)
	}
	return u
}

export const pathForInsightType = (it: InsightType | undefined): string => {
	return it === undefined ? "" : "/" + it
}

export const insightTypes = [
	"most-played-songs",
	"most-listened-artists",
	"longest-songs",
	"artist-discovery",
] as const

export type InsightType = typeof insightTypes[number]

export const defaultInsightType: InsightType = "most-played-songs"
