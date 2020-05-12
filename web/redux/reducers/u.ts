import { combineReducers } from "redux"
import { allScrobblesReducer, lovedScrobblesReducer, colorScrobblesReducer } from "./scrobbles"
import { songsReducer } from "./song"
import { lastColorReducer } from "./last-color"
import { uargsReducer } from "./uargs"

const reducer = combineReducers({
  allScrobbles: allScrobblesReducer,
  lovedScrobbles: lovedScrobblesReducer,
  colorScrobbles: colorScrobblesReducer,

  songs: songsReducer,

  lastColor: lastColorReducer,

  // TODO: is there a better way to specify a no-op reducer for the read-only
  // preloadedState property uargs?
  uargs: uargsReducer,
})

export default reducer
