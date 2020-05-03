import { AllScrobblesAction, LovedScrobblesAction, ColorScrobblesAction } from "../actions/scrobbles"
import { ScrobblesState } from "../types/scrobbles"
import { Song, ArtworkHash } from "../../shared/types"

const defaultState = <T extends Song | ArtworkHash>(): ScrobblesState<T> => {
  return {
    fetching: false,
    items: [],
    total: 0,
    private: false,
    error: false,
    done: false,
  }
}

const defaultAllScrobblesState = defaultState<Song>()
const defaultLovedScrobblesState = defaultState<Song>()

export const allScrobblesReducer = (state = defaultAllScrobblesState, action: AllScrobblesAction): ScrobblesState<Song> => {
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

export const lovedScrobblesReducer = (state = defaultLovedScrobblesState, action: LovedScrobblesAction): ScrobblesState<Song> => {
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
  const m: { [color: string]: ScrobblesState<ArtworkHash> } = {}
  colors.forEach(c => { m[c] = defaultState<ArtworkHash>() })
  return m
})()

type ColorScrobblesState = typeof defaultColorScrobblesState

export const colorScrobblesReducer = (state = defaultColorScrobblesState, action: ColorScrobblesAction): ColorScrobblesState => {
  const color = action.color

  switch (action.type) {
    case "COLOR_SCROBBLES_START": {
      const s = copy(state)
      s[color] = { ...s[color], fetching: true, done: false }
      return s
    }
    case "COLOR_SCROBBLES_SUCCESS": {
      const s = copy(state)
      s[color] = { items: action.hashes, private: action.private, fetching: false, error: false, done: true }
      return s
    }
    case "COLOR_SCROBBLES_FAIL": {
      const s = copy(state)
      s[color] = { ...s[color], fetching: false, error: true, done: true }
      return s
    }
    default: {
      return state
    }
  }
}

const copy = (m: ColorScrobblesState): ColorScrobblesState => {
  const n: ColorScrobblesState = {}
  for (const key in m) {
    n[key] = { ...m[key] }
  }
  return n
}
