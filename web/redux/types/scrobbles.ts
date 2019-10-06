import { Song } from "../../shared/types"

export type ScrobblesState = {
  fetching: boolean
  songs: Song[]
  private: boolean
  error: boolean
  done: boolean
}
