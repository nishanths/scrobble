import { ThunkAction, ThunkDispatch } from "redux-thunk"
import { Song, ScrobbledResponse } from "../../shared/types"
import { PartialState } from "../types/u"

export type AllScrobblesAction =
	ReturnType<typeof allScrobblesStart> |
	ReturnType<typeof allScrobblesSuccess> |
	ReturnType<typeof allScrobblesFail>

export type LovedScrobblesAction =
	ReturnType<typeof lovedScrobblesStart> |
	ReturnType<typeof lovedScrobblesSuccess> |
	ReturnType<typeof lovedScrobblesFail>

export type ColorScrobblesAction =
	ReturnType<typeof colorScrobblesStart> |
	ReturnType<typeof colorScrobblesSuccess> |
	ReturnType<typeof colorScrobblesFail>

type AllScrobblesThunkDispatch = ThunkDispatch<PartialState, undefined, AllScrobblesAction>
type AllScrobblesThunkResult<R> = ThunkAction<R, PartialState, undefined, AllScrobblesAction>

type LovedScrobblesThunkDispatch = ThunkDispatch<PartialState, undefined, LovedScrobblesAction>
type LovedScrobblesThunkResult<R> = ThunkAction<R, PartialState, undefined, LovedScrobblesAction>

type ColorScrobblesThunkDispatch = ThunkDispatch<PartialState, undefined, ColorScrobblesAction>
type ColorScrobblesThunkResult<R> = ThunkAction<R, PartialState, undefined, ColorScrobblesAction>

const allScrobblesStart = (username: string) => {
	return {
		type: "ALL_SCROBBLES_START" as const,
		username,
	}
}

const allScrobblesSuccess = (username: string, songs: Song[], total: number, priv: boolean) => {
	return {
		type: "ALL_SCROBBLES_SUCCESS" as const,
		username,
		songs,
		total,
		private: priv,
	}
}

const allScrobblesFail = (err: any) => {
	return {
		type: "ALL_SCROBBLES_FAIL" as const,
		err,
	}
}

const lovedScrobblesStart = (username: string) => {
	return {
		type: "LOVED_SCROBBLES_START" as const,
		username,
	}
}

const lovedScrobblesSuccess = (username: string, songs: Song[], total: number, priv: boolean) => {
	return {
		type: "LOVED_SCROBBLES_SUCCESS" as const,
		username,
		songs,
		total,
		private: priv,
	}
}

const lovedScrobblesFail = (err: any) => {
	return {
		type: "LOVED_SCROBBLES_FAIL" as const,
		err,
	}
}

export const fetchAllScrobbles = (username: string, limit: number): AllScrobblesThunkResult<void> => {
	return async (dispatch) => {
		dispatch(allScrobblesStart(username))

		try {
			const result = await _fetchScrobbledSongs(username, limit, false)
			dispatch(allScrobblesSuccess(username, result.songs, result.total, result.private))
		} catch (e) {
			dispatch(allScrobblesFail(e))
		}
	}
}

export const fetchLovedScrobbles = (username: string): LovedScrobblesThunkResult<void> => {
	return async (dispatch) => {
		dispatch(lovedScrobblesStart(username))

		try {
			const result = await _fetchScrobbledSongs(username, null, true)
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

const _fetchScrobbledSongs = async (username: string, limit: number | null, loved: boolean): Promise<FetchSongsResult> => {
	let url = "/api/v1/scrobbled?username=" + username
	if (limit !== null) {
		url += "&limit=" + limit;
	}
	if (loved === true) {
		url += "&loved=true"
	}
	const r = await fetch(url)
	switch (r.status) {
		case 200: {
			const rsp: ScrobbledResponse = await r.json()
			return { songs: rsp.songs, total: rsp.total, private: false, err: null }
		}
		case 404:
			return { songs: [], total: 0, private: true, err: null }
		// TODO: if we had the ability to display toast notifications, we could show
		// "please sign in again" for 401 status
		default:
			throw "bad status: " + r.status
	}
}

const colorScrobblesStart = (color: string, username: string) => {
	return {
		type: "COLOR_SCROBBLES_START" as const,
		color,
		username,
	}
}

const colorScrobblesSuccess = (color: string, username: string, songs: Song[], priv: boolean) => {
	return {
		type: "COLOR_SCROBBLES_SUCCESS" as const,
		color,
		username,
		songs,
		private: priv,
	}
}

const colorScrobblesFail = (err: any, color: string) => {
	return {
		type: "COLOR_SCROBBLES_FAIL" as const,
		color,
		err,
	}
}

export const fetchColorScrobbles = (color: string, username: string): ColorScrobblesThunkResult<void> => {
	return async (dispatch) => {
		dispatch(colorScrobblesStart(color, username))

		try {
			const result = await _fetchColorScrobbles(color, username)
			dispatch(colorScrobblesSuccess(color, username, result.songs, result.private))
		} catch (e) {
			dispatch(colorScrobblesFail(e, color))
		}
	}
}

type FetchColorResult = {
	songs: Song[]
	private: boolean
	err: any | null
}

const _fetchColorScrobbles = async (color: string, username: string): Promise<FetchColorResult> => {
	const url = "/api/v1/scrobbled/color?username=" + username + "&color=" + color;

	const r = await fetch(url)
	switch (r.status) {
		case 200: {
			const rsp: ScrobbledResponse = await r.json()
			return { songs: rsp.songs, private: false, err: null }
		}
		case 403:
			return { songs: [], private: true, err: null }
		// TODO: if we had the ability to display toast notifications, we could show
		// "please sign in again" for 401 status
		default:
			throw "bad status: " + r.status
	}
}
