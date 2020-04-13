import { AllScrobblesAction, LovedScrobblesAction } from "../actions/scrobbles"
import { ScrobblesState } from "../types/scrobbles"

const defaultState = (): ScrobblesState => {
  return {
    fetching: false,
    songs: [],
    total: 0,
    private: false,
    error: false,
    done: false,
  }
}

const defaultAllScrobblesState = defaultState()
const defaultLovedScrobblesState = defaultState()

export const allScrobblesReducer = (state: ScrobblesState = defaultAllScrobblesState, action: AllScrobblesAction): ScrobblesState => {
  switch (action.type) {
    case "ALL_SCROBBLES_START":
      return { ...state, fetching: true, done: false }
    case "ALL_SCROBBLES_SUCCESS":
      return { songs: action.songs, total: action.total, private: action.private, fetching: false, error: false, done: true }
    case "ALL_SCROBBLES_FAIL":
      return { ...state, fetching: false, error: true, done: true }
    default:
      return state
  }
}

export const lovedScrobblesReducer = (state: ScrobblesState = defaultLovedScrobblesState, action: LovedScrobblesAction): ScrobblesState => {
  switch (action.type) {
    case "LOVED_SCROBBLES_START":
      return { ...state, fetching: true, done: false }
    case "LOVED_SCROBBLES_SUCCESS":
      return { ...state, songs: action.songs, total: action.total, private: action.private, fetching: false, error: false, done: true }
    case "LOVED_SCROBBLES_FAIL":
      return { ...state, fetching: false, error: true, done: true }
    default:
      return state
  }
}
