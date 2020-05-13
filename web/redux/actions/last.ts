import { Color } from "../../components/colorpicker"

export type LastAction =
	| ReturnType<typeof setLastColor>
	| ReturnType<typeof setLastScrobblesEndIdx>

export const setLastColor = (c: Color) => {
  return {
    type: "SET_LAST_COLOR" as const,
    color: c,
  }
}

export const setLastScrobblesEndIdx = (i: number) => {
  return {
    type: "SET_LAST_SCROBBLES_END_IDX" as const,
    scrobblesEndIdx: i,
  }
}
