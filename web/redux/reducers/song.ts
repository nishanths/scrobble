import { SongAction } from "../actions/song"
import { SongState } from "../types/song"
import { MapDefault } from "../../shared/util"

const defaultStateValueFunc = (): SongState => {
  return {
    fetching: false,
    private: false,
    notFound: false,
    item: null,
    error: false,
    done: false,
  }
}

const defaultState = new MapDefault<string, SongState>(defaultStateValueFunc)

export const songReducer = (state = defaultState, action: SongAction): typeof defaultState => {
  switch (action.type) {
    case "SONG_START": {
      const s = state.copy()
      const v = s.getOrDefault(action.ident)
      s.set(action.ident, { ...v, fetching: true })
      return s
    }
    case "SONG_SUCCESS": {
      const s = state.copy()
      const v = s.getOrDefault(action.ident)
      s.set(action.ident, { fetching: false, done: true, error: false, notFound: action.notFound, item: action.song, private: action.private  })
      return s
    }
    case "SONG_FAIL": {
      const s = state.copy()
      const v = s.getOrDefault(action.ident)
      s.set(action.ident, { ...v, fetching: false, done: true, error: true  })
      return s
    }
    default: {
      return state
    }
  }
}
