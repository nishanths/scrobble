import { ThunkAction, ThunkDispatch } from "redux-thunk"
import { Dispatch } from "redux"
import { Song } from "../../shared/types"
import { PartialState } from "../types/u"

export type ScrobblesAction =
  ReturnType<typeof scrobblesStart> |
  ReturnType<typeof scrobblesSuccess> |
  ReturnType<typeof scrobblesFail>

type ScrobblesThunkDispatch = ThunkDispatch<PartialState, undefined, ScrobblesAction>
type ScrobblesThunkResult<R> = ThunkAction<R, PartialState, undefined, ScrobblesAction>

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

export const fetchScrobbles = (username: string, limit: number): ScrobblesThunkResult<void> => {
  return async (dispatch, store) => {
    dispatch(scrobblesStart(username))
    try {
      const result = await _fetchScrobbles(username, limit)
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

const _fetchScrobbles = async (username: string, limit: number): Promise<FetchScrobblesResult> => {
  const r = await fetch("/api/v1/scrobbled?username=" + username + "&limit=" + limit)
  switch (r.status) {
    case 200:
      const songs: Song[] = await r.json()
      return { songs, private: false, err: null }
    case 404:
      return { songs: [], private: true, err: null }
    // TODO: if we had the ability to display toast notifications, we could show
    // "please sign in again" for 401 status
    default:
      throw "bad status: " + r.status
  }
}
