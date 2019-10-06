import { combineReducers } from "redux"
import { scrobblesReducer } from "./scrobbles"
import { UArgs } from "../../shared/types"

export type State = {
  scrobbles: ReturnType<typeof scrobblesReducer>
  uargs: UArgs
}

const reducer = combineReducers({
  scrobbles: scrobblesReducer,

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

export default reducer
