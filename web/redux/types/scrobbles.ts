import { Song } from "../../shared/types"

export type ScrobblesState = {
  fetching: boolean
  songs: Song[] | null
  private: boolean
  error: boolean
}
