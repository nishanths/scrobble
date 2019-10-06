import { Song } from "../../shared/types"

export const scrobblesStart = (username: string) => {
  return {
    type: "SCROBBLES_START" as const,
    username,
  }
}

export const scrobblesSuccess = (username: string, songs: Song[], priv: boolean) => {
  return {
    type: "SCROBBLES_SUCCESS" as const,
    username,
    songs,
    private: priv,
  }
}

export const scrobblesFail = () => {
  return {
    type: "SCROBBLES_FAIL" as const,
  }
}

export type ScrobblesAction =
  ReturnType<typeof scrobblesStart> |
  ReturnType<typeof scrobblesSuccess> |
  ReturnType<typeof scrobblesFail>


type FetchScrobblesResult = {
  songs: Song[]
  private: boolean
  err: any | null
}

const _fetchScrobbles = async (username: string): Promise<FetchScrobblesResult> => {
  try {
    const r = await fetch("/api/v1/scrobbled?username=" + username)
    switch (r.status) {
      case 200:
        const songs: Song[] = await r.json()
      case 404:
        return { songs: [], private: true, err: null }
      default:
        return { songs: [], private: false, err: "bad status: " + r.status }
    }
  } catch(err) {
    return { songs: [], private: false, err }
  }
}
