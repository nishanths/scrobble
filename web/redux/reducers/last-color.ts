import { LastColorState } from "../types/last-color"
import { LastColorAction } from "../actions/last-color"

const defaultState: LastColorState = {
  color: undefined
}

export const lastColorReducer = (state = defaultState, action: LastColorAction): LastColorState => {
  switch (action.type) {
    case "SET_LAST_COLOR":
      return { color: action.color }
    default:
      return state
  }
}
