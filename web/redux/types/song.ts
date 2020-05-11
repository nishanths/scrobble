import { Song } from "../../shared/types"

export type SongState = {
  fetching: boolean
  item: Song | null
  notFound: boolean
  private: boolean
  error: boolean
  done: boolean
}
