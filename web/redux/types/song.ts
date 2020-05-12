import { Song } from "../../shared/types"

export type SongState = {
  fetching: boolean
  item: Song | null // non-null when done is true and private,error,noSuchSong are false
  noSuchSong: boolean
  private: boolean
  error: boolean
  done: boolean
}
