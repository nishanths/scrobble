import { ScrobblesAction } from "../actions/scrobbles"
import { ScrobblesState } from "../types/scrobbles"

const scrobblesReducer = (state: ScrobblesState, action: ScrobblesAction): ScrobblesState => {
  switch (action.type) {
    case "SCROBBLES_START":
      return { ...state, fetching: true, done: false }
    case "SCROBBLES_SUCCESS":
      return { songs: action.songs, total: action.total, private: action.private, fetching: false, error: false, done: true }
    case "SCROBBLES_FAIL":
      return { ...state, fetching: false, error: true, done: true }
    default:
      return state
  }
}

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

export const allScrobblesReducer = (state: ScrobblesState = defaultAllScrobblesState, action: ScrobblesAction): ScrobblesState => {
  return scrobblesReducer(state, action)
}

export const lovedScrobblesReducer = (state: ScrobblesState = defaultLovedScrobblesState, action: ScrobblesAction): ScrobblesState => {
  return scrobblesReducer(state, action)
}
