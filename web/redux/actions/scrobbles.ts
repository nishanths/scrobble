import { Song } from "../../shared/types"

export const scrobblesStart = (username: string) => {
  return {
    type: "SCROBBLES_START" as const,
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
  ReturnType<typeof scrobblesStart> |
  ReturnType<typeof scrobblesSuccess> |
  ReturnType<typeof scrobblesFail>
