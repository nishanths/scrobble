import { Song } from "../../shared/types"

export type ScrobblesState = {
  fetching: boolean
  songs: Song[]
  total: number
  private: boolean
  error: boolean
  done: boolean
}
