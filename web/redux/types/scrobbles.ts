import { Song, ArtworkHash } from "../../shared/types"

export type ScrobblesState<T extends Song | ArtworkHash> = {
  fetching: boolean
  items: T[]
  total?: number
  private: boolean
  error: boolean
  done: boolean
}
