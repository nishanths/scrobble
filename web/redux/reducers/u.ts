import { combineReducers } from "redux"
import { allScrobblesReducer, lovedScrobblesReducer } from "./scrobbles"
import { uargsReducer } from "./uargs"

const reducer = combineReducers({
  allScrobbles: allScrobblesReducer,
  lovedScrobbles: lovedScrobblesReducer,

  // TODO: is there a better way to specify a no-op reducer for the read-only
  // preloadedState property uargs?
  uargs: uargsReducer,
})

export default reducer
