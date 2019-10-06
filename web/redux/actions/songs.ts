import { Song } from "../../shared/types"

export const fetchScrobblesRequest = (username: string) => {
  return {
    type: "FETCH_SCROBBLES_REQUEST" as const,
    username,
  }
}

export const fetchScrobblesSuccess = (username: string, songs: Song[]) => {
  return {
    type: "FETCH_SCROBBLES_SUCCESS" as const,
    username,
    songs,
  }
}

export const fetchScrobblesFail = (username: string) => {
  return {
    type: "FETCH_SCROBBLES_FAIL" as const,
    username,
  }
}

export type SongsAction =
  ReturnType<typeof fetchScrobblesRequest> |
  ReturnType<typeof fetchScrobblesSuccess> |
  ReturnType<typeof fetchScrobblesFail>
