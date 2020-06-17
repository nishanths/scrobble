import { AllScrobblesAction, LovedScrobblesAction, ColorScrobblesAction } from "../actions/scrobbles"
import { ScrobblesState } from "../types/scrobbles"
import { copyMap } from "../../shared/util"

const defaultStateFunc = (): ScrobblesState => {
    return {
        fetching: false,
        items: [],
        total: 0,
        private: false,
        error: false,
        done: false,
    }
}

const defaultAllScrobblesState = defaultStateFunc()
const defaultLovedScrobblesState = defaultStateFunc()

export const allScrobblesReducer = (state = defaultAllScrobblesState, action: AllScrobblesAction): ScrobblesState => {
    switch (action.type) {
        case "ALL_SCROBBLES_START":
            return { ...state, fetching: true, done: false }
        case "ALL_SCROBBLES_SUCCESS":
            return { items: action.songs, total: action.total, private: action.private, fetching: false, error: false, done: true }
        case "ALL_SCROBBLES_FAIL":
            return { ...state, fetching: false, error: true, done: true }
        default:
            return state
    }
}

export const lovedScrobblesReducer = (state = defaultLovedScrobblesState, action: LovedScrobblesAction): ScrobblesState => {
    switch (action.type) {
        case "LOVED_SCROBBLES_START":
            return { ...state, fetching: true, done: false }
        case "LOVED_SCROBBLES_SUCCESS":
            return { items: action.songs, total: action.total, private: action.private, fetching: false, error: false, done: true }
        case "LOVED_SCROBBLES_FAIL":
            return { ...state, fetching: false, error: true, done: true }
        default:
            return state
    }
}

const colors = [
    "gray",
    "red",
    "orange",
    "blue",
    "yellow",
    "white",
    "green",
    "violet",
    "pink",
    "brown",
    "black",
];

const defaultColorScrobblesState = (() => {
    // create a map that maps each color to a default scrobbles state for that color
    const m: Map<string, ScrobblesState> = new Map()
    colors.forEach(c => { m.set(c, defaultStateFunc()) })
    return m
})()

export const colorScrobblesReducer = (state = defaultColorScrobblesState, action: ColorScrobblesAction): Map<string, ScrobblesState> => {
    const color = action.color

    switch (action.type) {
        case "COLOR_SCROBBLES_START": {
            const s = copyMap(state)
            s.set(color, { ...s.get(color)!, fetching: true, done: false })
            return s
        }
        case "COLOR_SCROBBLES_SUCCESS": {
            const s = copyMap(state)
            s.set(color, { items: action.songs, private: action.private, fetching: false, error: false, done: true })
            return s
        }
        case "COLOR_SCROBBLES_FAIL": {
            const s = copyMap(state)
            s.set(color, { ...s.get(color)!, fetching: false, error: true, done: true })
            return s
        }
        default: {
            return state
        }
    }
}
