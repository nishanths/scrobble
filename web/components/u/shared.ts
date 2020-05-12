import { assertExhaustive } from "../../shared/util"
import { Color } from "../colorpicker"

export enum DetailKind {
  Song, Album
}

export enum Mode {
  All, Loved, Color
}

export const controlValues = ["All", "Loved", "By color"] as const

export type ControlValue = typeof controlValues[number]

export const controlValueForMode = (m: Mode): ControlValue => {
  switch (m) {
    case Mode.All: return "All"
    case Mode.Loved: return "Loved"
    case Mode.Color: return "By color"
    default: assertExhaustive(m)
  }
}

export const modeFromControlValue = (v: ControlValue): Mode => {
  switch (v) {
    case "All": return Mode.All
    case "Loved": return Mode.Loved
    case "By color": return Mode.Color
    default: assertExhaustive(v)
  }
}

export const pathForMode = (m: Mode): string => {
  switch (m) {
    case Mode.All: return ""
    case Mode.Loved: return "/loved"
    case Mode.Color: return "/color"
  }
  assertExhaustive(m)
}

export const pathForColor = (c: Color | undefined): string => {
  if (c === undefined) {
    return ""
  }
  return "/" + c
}
