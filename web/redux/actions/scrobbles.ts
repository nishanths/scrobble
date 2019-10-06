import { Song } from "../../shared/types"

export const scrobblesRequest = (username: string) => {
  return {
    type: "SCROBBLES_REQUEST" as const,
    username,
  }
}

export const scrobblesSuccess = (username: string, songs: Song[]) => {
  return {
    type: "SCROBBLES_SUCCESS" as const,
    username,
    songs,
  }
}

export const scrobblesFail = (username: string) => {
  return {
    type: "SCROBBLES_FAIL" as const,
    username,
  }
}

export type ScrobblesAction =
  ReturnType<typeof scrobblesRequest> |
  ReturnType<typeof scrobblesSuccess> |
  ReturnType<typeof scrobblesFail>
