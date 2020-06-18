import { Color } from "../../components/colorpicker"
import { Mode } from "../../components/u/shared"

export type LastAction =
	| ReturnType<typeof setLastColor>
	| ReturnType<typeof setLastScrobblesEndIdx>
	| ReturnType<typeof setLastScrobblesScrollY>
	| ReturnType<typeof setLastSearch>

export const setLastColor = (c: Color) => {
	return {
		type: "SET_LAST_COLOR" as const,
		color: c,
	}
}

export const setLastScrobblesEndIdx = (i: number) => {
	return {
		type: "SET_LAST_SCROBBLES_END_IDX" as const,
		endIdx: i,
	}
}

export const setLastScrobblesScrollY = (y: number | undefined) => {
	return {
		type: "SET_LAST_SCROBBLES_SCROLL_Y" as const,
		scrollY: y,
	}
}

export const setLastSearch = (s: string) => {
	return {
		type: "SET_LAST_SEARCH" as const,
		search: s,
	}
}
