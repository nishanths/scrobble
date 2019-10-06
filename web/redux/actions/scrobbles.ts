import { ThunkAction, ThunkDispatch } from "redux-thunk"
import { Dispatch } from "redux"
import { Song } from "../../shared/types"
import { PartialState } from "../types/u"

export type ScrobblesAction =
  ReturnType<typeof scrobblesStart> |
  ReturnType<typeof scrobblesSuccess> |
  ReturnType<typeof scrobblesFail>


type ScrobblesThunkDispatch = ThunkDispatch<PartialState, void, ScrobblesAction>
type ScrobblesThunkAction<R> = ThunkAction<R, PartialState, void, ScrobblesAction>

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

export const scrobblesFail = (err: any) => {
  return {
    type: "SCROBBLES_FAIL" as const,
    err,
  }
}

export const fetchScrobbles = (username: string): ScrobblesThunkAction<void> => {
  return async (dispatch: ScrobblesThunkDispatch, _: () => PartialState) => {
    dispatch(scrobblesStart(username))
    try {
      const result = await _fetchScrobbles(username)
      dispatch(scrobblesSuccess(username, result.songs, result.private))
    } catch (e) {
      dispatch(scrobblesFail(e))
    }
  }
}

type FetchScrobblesResult = {
  songs: Song[]
  private: boolean
  err: any | null
}

const _fetchScrobbles = async (username: string): Promise<FetchScrobblesResult> => {
  const r = await fetch("/api/v1/scrobbled?username=" + username)
  switch (r.status) {
    case 200:
      const songs: Song[] = await r.json()
      return { songs, private: false, err: null }
    case 404:
      throw { songs: [], private: true, err: null }
    default:
      throw "bad status: " + r.status
  }
}
