import { combineReducers } from "redux"
import { songsReducer } from "./songs"
import { UArgs } from "../../shared/types"

export type State = {
  songs: ReturnType<typeof songsReducer>
  uargs: UArgs
}

export default combineReducers<State>({
  songs: songsReducer,

  // TODO: is there a better way to specify a no-op reducer for the read-only
  // preloadedState property uargs?
  uargs: (s: UArgs | undefined): UArgs => (s !== undefined ? s : {
    artworkBaseURL: "",
    host: "",
    self: false,
    profileUsername: "",
    logoutURL: "",
    account: {
      apiKey: "",
      username: "",
      private: false,
    },
  }),
})
