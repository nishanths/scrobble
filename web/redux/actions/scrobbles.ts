import { ThunkAction, ThunkDispatch } from "redux-thunk"
import { Dispatch } from "redux"
import { Song, ScrobbledResponse, ArtworkHash } from "../../shared/types"
import { PartialState } from "../types/u"

export type AllScrobblesAction =
  ReturnType<typeof allScrobblesStart> |
  ReturnType<typeof allScrobblesSuccess> |
  ReturnType<typeof allScrobblesFail>

export type LovedScrobblesAction =
  ReturnType<typeof lovedScrobblesStart> |
  ReturnType<typeof lovedScrobblesSuccess> |
  ReturnType<typeof lovedScrobblesFail>

type AllScrobblesThunkDispatch = ThunkDispatch<PartialState, undefined, AllScrobblesAction>
type AllScrobblesThunkResult<R> = ThunkAction<R, PartialState, undefined, AllScrobblesAction>

type LovedScrobblesThunkDispatch = ThunkDispatch<PartialState, undefined, LovedScrobblesAction>
type LovedScrobblesThunkResult<R> = ThunkAction<R, PartialState, undefined, LovedScrobblesAction>

export const allScrobblesStart = (username: string) => {
  return {
    type: "ALL_SCROBBLES_START" as const,
    username,
  }
}

export const allScrobblesSuccess = (username: string, songs: Song[], total: number, priv: boolean) => {
  return {
    type: "ALL_SCROBBLES_SUCCESS" as const,
    username,
    songs,
    total,
    private: priv,
  }
}

export const allScrobblesFail = (err: any) => {
  return {
    type: "ALL_SCROBBLES_FAIL" as const,
    err,
  }
}

export const lovedScrobblesStart = (username: string) => {
  return {
    type: "LOVED_SCROBBLES_START" as const,
    username,
  }
}

export const lovedScrobblesSuccess = (username: string, songs: Song[], total: number, priv: boolean) => {
  return {
    type: "LOVED_SCROBBLES_SUCCESS" as const,
    username,
    songs,
    total,
    private: priv,
  }
}

export const lovedScrobblesFail = (err: any) => {
  return {
    type: "LOVED_SCROBBLES_FAIL" as const,
    err,
  }
}

export const fetchAllScrobbles = (username: string, limit: number): AllScrobblesThunkResult<void> => {
  return async (dispatch, store) => {
    dispatch(allScrobblesStart(username))

    try {
      const result = await _fetchScrobbledSongs(username, limit, false)
      dispatch(allScrobblesSuccess(username, result.songs, result.total, result.private))
    } catch (e) {
      dispatch(allScrobblesFail(e))
    }
  }
}

export const fetchLovedScrobbles = (username: string, limit: number): LovedScrobblesThunkResult<void> => {
  return async (dispatch, store) => {
    dispatch(lovedScrobblesStart(username))

    try {
      const result = await _fetchScrobbledSongs(username, limit, true)
      dispatch(lovedScrobblesSuccess(username, result.songs, result.total, result.private))
    } catch (e) {
      dispatch(lovedScrobblesFail(e))
    }
  }
}

type FetchSongsResult = {
  songs: Song[]
  total: number
  private: boolean
  err: any | null
}

const _fetchScrobbledSongs = async (username: string, limit: number, loved: boolean): Promise<FetchSongsResult> => {
  let url = "/api/v1/scrobbled?username=" + username + "&limit=" + limit;
  if (loved === true) {
    url += "&loved=true"
  }
  const r = await fetch(url)
  switch (r.status) {
    case 200:
      const rsp: ScrobbledResponse = await r.json()
      return { songs: rsp.songs, total: rsp.total, private: false, err: null }
    case 404:
      return { songs: [], total: 0, private: true, err: null }
    // TODO: if we had the ability to display toast notifications, we could show
    // "please sign in again" for 401 status
    default:
      throw "bad status: " + r.status
  }
}

// TODO

export const colorScrobblesStart = (color: string, username: string) => {
  return {
    type: "COLOR_SCROBBLES_START" as const,
    color,
    username,
  }
}

type FetchColorResult = {
  hashes: ArtworkHash[]
  private: boolean
  err: any | null
}
