import { combineReducers } from "redux"
import { songs } from "./songs"
import { SongsState, UArgs } from "../../shared/types"

export type State = {
	songs: SongsState
	uargs: UArgs
}

export default combineReducers<State>({
	songs,
	uargs: () => ({}), // any non-null value will do
})
