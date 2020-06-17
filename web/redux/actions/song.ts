import { ThunkAction, ThunkDispatch } from "redux-thunk"
import { Song, ScrobbledResponse } from "../../shared/types"
import { PartialState } from "../types/u"

export type SongAction =
	| ReturnType<typeof songStart>
	| ReturnType<typeof songSuccess>
	| ReturnType<typeof songFail>

type SongThunkDispatch = ThunkDispatch<PartialState, undefined, SongAction>
type SongThunkResult<R> = ThunkAction<R, PartialState, undefined, SongAction>

const songStart = (username: string, ident: string) => {
	return {
		type: "SONG_START" as const,
		username,
		ident,
	}
}

const songSuccess = (username: string, ident: string, song: Song | null, noSuchSong: boolean, priv: boolean) => {
	return {
		type: "SONG_SUCCESS" as const,
		username,
		ident,
		song,
		noSuchSong,
		private: priv,
	}
}

const songFail = (ident: string, err: any) => {
	return {
		type: "SONG_FAIL" as const,
		ident,
		err,
	}
}

type FetchSongResult = {
	song: Song | null
	noSuchSong: boolean
	private: boolean
	err: any | null
}

export const fetchSong = (username: string, ident: string): SongThunkResult<void> => {
	return async (dispatch) => {
		dispatch(songStart(username, ident))
		try {
			const result = await _fetchSong(username, ident)
			dispatch(songSuccess(username, ident, result.song, result.noSuchSong, result.private))
		} catch (e) {
			dispatch(songFail(ident, e))
		}
	}
}

const _fetchSong = async (username: string, ident: string): Promise<FetchSongResult> => {
	const url = "/api/v1/scrobbled?username=" + username + "&song=" + encodeURIComponent(ident)
	const r = await fetch(url)
	switch (r.status) {
		case 200: {
			const rsp: ScrobbledResponse = await r.json()
			if (rsp.songs.length === 0) {
				// no song for given ident
				return { song: null, noSuchSong: true, private: false, err: null }
			}
			return { song: rsp.songs[0], noSuchSong: false, private: false, err: null }
		}
		case 403:
			return { song: null, noSuchSong: false, private: true, err: null }
		default:
			throw "bad status: " + r.status
	}
}
