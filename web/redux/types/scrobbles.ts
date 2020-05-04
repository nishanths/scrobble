import { Song, ArtworkHash } from "../../shared/types"

export type ScrobblesState = {
  fetching: boolean
  items: Song[]
  total?: number
  private: boolean
  error: boolean
  done: boolean
}
