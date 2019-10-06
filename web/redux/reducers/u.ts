import { combineReducers } from "redux"
import { scrobblesReducer } from "./scrobbles"
import { uargsReducer } from "./uargs"

const reducer = combineReducers({
  scrobbles: scrobblesReducer,

  // TODO: is there a better way to specify a no-op reducer for the read-only
  // preloadedState property uargs?
  uargs: uargsReducer,
})

export default reducer
