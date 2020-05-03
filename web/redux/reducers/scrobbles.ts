import { AllScrobblesAction, LovedScrobblesAction } from "../actions/scrobbles"
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

export const allScrobblesReducer = (state: ScrobblesState<Song> = defaultAllScrobblesState, action: AllScrobblesAction): ScrobblesState<Song> => {
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

export const lovedScrobblesReducer = (state: ScrobblesState<Song> = defaultLovedScrobblesState, action: LovedScrobblesAction): ScrobblesState<Song> => {
  switch (action.type) {
    case "LOVED_SCROBBLES_START":
      return { ...state, fetching: true, done: false }
    case "LOVED_SCROBBLES_SUCCESS":
      return { ...state, items: action.songs, total: action.total, private: action.private, fetching: false, error: false, done: true }
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
  const m: { [k: string]: ScrobblesState<ArtworkHash> } = {}
  colors.forEach(c => { m[c] = defaultState<ArtworkHash>() })
  return m
})()
