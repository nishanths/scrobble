import { ScrobblesAction } from "../actions/scrobbles"
import { ScrobblesState } from "../types/scrobbles"

const defaultState: ScrobblesState = {
  fetching: false,
  songs: [],
  private: false,
  error: false,
}

export const scrobblesReducer = (state: ScrobblesState = defaultState, action: ScrobblesAction): ScrobblesState => {
  switch (action.type) {
    case "SCROBBLES_START":
      return { ...state, fetching: true }
    case "SCROBBLES_SUCCESS":
      return { songs: action.songs, private: action.private, fetching: false, error: false }
    case "SCROBBLES_FAIL":
      return { ...state, fetching: false, error: true }
    default:
      return state
  }
}
