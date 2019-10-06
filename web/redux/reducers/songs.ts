import { Song } from "../../shared/types"
import { SongsAction } from "../actions/songs"

type SongsState = {
  fetching: boolean
  songs: Song[]
  error: boolean
}

const defaultSongsState: SongsState = {
  fetching: false,
  songs: [],
  error: false,
}

export const songsReducer = (state: SongsState = defaultSongsState, action: SongsAction): SongsState => {
  if (state === undefined) {
    return { fetching: false, songs: [], error: false }
  }
  switch (action.type) {
    case "FETCH_SCROBBLES_REQUEST":
      return { ...state, fetching: true }
    case "FETCH_SCROBBLES_SUCCESS":
      return { songs: action.songs, fetching: false, error: false }
    case "FETCH_SCROBBLES_FAIL":
      return { ...state, fetching: false, error: true }
    default:
      return state
  }
}
