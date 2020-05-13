import { LastState } from "../types/last"
import { LastAction } from "../actions/last"

const defaultState: LastState = {
  color: undefined,
  scrobblesEndIdx: undefined,
}

export const lastReducer = (state = defaultState, action: LastAction): LastState => {
  switch (action.type) {
    case "SET_LAST_COLOR":
      return { ...state, color: action.color }
    case "SET_LAST_SCROBBLES_END_IDX":
      return { ...state, scrobblesEndIdx: action.scrobblesEndIdx }
    default:
      return state
  }
}
