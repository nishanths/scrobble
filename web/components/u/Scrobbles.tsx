import React, { useEffect, useRef, useState } from "react"
import { useDispatch, useSelector } from "react-redux"
import { RouteComponentProps } from "react-router-dom";
import { State } from "../../redux/types/u"
import { assertExhaustive, assert, hexEncode } from "../../shared/util"
import { NProgress, Song } from "../../shared/types"
import { Mode, DetailKind, pathForMode, pathForColor, pathForDetailKind, modeFromControlValue } from "./shared"
import { Color } from "../colorpicker"
import { Songs } from "../Songs"
import { setLastColor, setLastScrobblesEndIdx, setLastScrobblesScrollY } from "../../redux/actions/last"
import { fetchAllScrobbles, fetchLovedScrobbles, fetchColorScrobbles } from "../../redux/actions/scrobbles"
import { Header, ColorPicker, Top } from "./top"
import { SearchBox } from "../searchbox"

// Divisble by 2, 3, and 4. This is appropriate because these are the number
// of cards typically displayed per row. Using such a number ensures that
// the last row isn't an incomplete row.
const moreIncrement = 36;

const limit = 504; // `moreIncrement` * 14

const nextEndIdx = (currentEndIdx: number, total: number): number => {
	// increment, but don't go over the number of items itself
	const b = Math.min(currentEndIdx + moreIncrement, total)
	// if there aren't sufficient items left for the next time, just include them now
	return total - b < moreIncrement ? total : b;
}

type History = RouteComponentProps["history"]

const searchPlaceholder = "Filter by album, artist, or song title…"

export const Scrobbles: React.StatelessComponent<{
	profileUsername: string
	signedIn: boolean
	artworkBaseURL: string
	private: boolean
	self: boolean
	mode: Mode
	color: Color | undefined
	nProgress: NProgress
	history: History
	wnd: Window
}> = ({
	profileUsername,
	signedIn,
	artworkBaseURL,
	private: priv,
	self,
	mode,
	color,
	nProgress,
	history,
	wnd,
}) => {
		const dispatch = useDispatch()
		const last = useSelector((s: State) => s.last)

		const [endIdx, setEndIdx] = useState(last.scrobblesEndIdx || 0)
		const endIdxRef = useRef(endIdx)
		useEffect(() => { endIdxRef.current = endIdx }, [endIdx])

		const [scrollY, setScrollY] = useState(wnd.pageYOffset)

		const shouldUpdateScrollTo = () => last.scrobblesScrollY !== undefined
		useEffect(() => {
			if (shouldUpdateScrollTo()) {
				wnd.scrollTo({ top: last.scrobblesScrollY })
				dispatch(setLastScrobblesScrollY(undefined))
			}
		})

		const [searchValue, setSearchValue] = useState("")

		const onControlChange = (newMode: Mode) => {
			nProgress.done()
			dispatch(setLastScrobblesEndIdx(0))

			let u: string;
			switch (newMode) {
				case Mode.All:
				case Mode.Loved:
					u = "/u/" + profileUsername + pathForMode(newMode)
					break
				case Mode.Color:
					u = "/u/" + profileUsername + pathForMode(newMode) + pathForColor(last.color)
					break
				default:
					assertExhaustive(newMode)
			}
			history.push(u)
		}

		const onColorChange = (newColor: Color) => {
			nProgress.done()
			dispatch(setLastScrobblesEndIdx(0))
			dispatch(setLastColor(newColor))
			assert(mode === Mode.Color, "mode should be Color")
			history.push("/u/" + profileUsername + pathForMode(mode) + pathForColor(newColor))
		}

		const onSongClick = (s: Song) => {
			nProgress.done()
			dispatch(setLastScrobblesEndIdx(endIdx))
			dispatch(setLastScrobblesScrollY(Math.floor(scrollY)))

			let kind: DetailKind
			switch (mode) {
				case Mode.All:
				case Mode.Loved:
					kind = DetailKind.Song
					break
				case Mode.Color:
					kind = DetailKind.Album
					break
				default:
					assertExhaustive(mode)
			}

			history.push("/u/" + profileUsername + pathForMode(mode) + pathForColor(color) + pathForDetailKind(kind) + "/" + hexEncode(s.ident))
		}

		// scrobbles redux state
		const scrobbles = useSelector((s: State) => {
			switch (mode) {
				case Mode.All: return s.allScrobbles
				case Mode.Loved: return s.lovedScrobbles
				case Mode.Color: return color !== undefined ? s.colorScrobbles.get(color)! : null
			}
			assertExhaustive(mode)
		})
		const scrobblesRef = useRef(scrobbles)
		useEffect(() => { scrobblesRef.current = scrobbles }, [scrobbles])

		// initial end idx
		useEffect(() => {
			const s = scrobblesRef.current
			if (s === null) {
				return
			}
			const e = s.error === false ? nextEndIdx(0, s.items.length) : 0
			if (e > endIdx) {
				setEndIdx(e)
			}
		}, [scrobbles, mode])

		// fetch scrobbles
		useEffect(() => {
			const s = scrobblesRef.current

			switch (mode) {
				case Mode.All: {
					if ((s!.done === false && s!.fetching === false) || s!.error === true) {
						dispatch(fetchAllScrobbles(profileUsername, limit))
					}
					break
				}
				case Mode.Loved: {
					if ((s!.done === false && s!.fetching === false) || s!.error === true) {
						dispatch(fetchLovedScrobbles(profileUsername, limit))
					}
					break
				}
				case Mode.Color: {
					if (color === undefined) {
						break
					}
					if (s === null || (s.done === false && s.fetching === false) || s.error === true) {
						dispatch(fetchColorScrobbles(color, profileUsername))
					}
					break
				}
				default: {
					assertExhaustive(mode)
				}
			}
		}, [profileUsername, mode, color])

		// infinite scroll
		useEffect(() => {
			const f = () => {
				const y = wnd.pageYOffset
				setScrollY(y)

				const s = scrobblesRef.current
				if (s === null) {
					return
				}

				const leeway = 250
				if ((wnd.innerHeight + y) >= (wnd.document.body.offsetHeight - leeway)) {
					const newEnd = nextEndIdx(endIdxRef.current, s.items.length)
					const e = Math.max(newEnd, endIdxRef.current)
					setEndIdx(e)
				}
			}

			wnd.addEventListener("scroll", f)
			return () => { wnd.removeEventListener("scroll", f) }
		}, [scrobbles, mode])

		// ... render ...

		const header = Header(profileUsername, signedIn, true)
		const colorPicker = ColorPicker(color, onColorChange)
		const searchBox = <div className="searchBox">
			<SearchBox
				value={searchValue}
				onChange={(v) => setSearchValue(v)}
				placeholder={searchPlaceholder}
			/>
		</div>
		const top = Top(header, colorPicker, mode, (v) => { onControlChange(modeFromControlValue(v)) })

		// Easy case. For private accounts that aren't the current user, render the
		// private info-message.
		if (priv === true && self === false) {
			return <>
				{header}
				<div className="info">(This user's scrobbles are private.)</div>
			</>
		}

		// If in the Color mode and no color is selected, render the top area and
		// the color picker, and we're done.
		if (mode === Mode.Color && color === undefined) {
			return <>
				{top}
			</>
		}

		assert(scrobbles !== null, "scrobbles unexpectedly null")

		if (scrobbles.fetching === true) {
			nProgress.start()
			return <>{top}</>
		}

		if (scrobbles.error === true) {
			nProgress.done()
			return <>
				{header}
				<div className="info">(Failed to fetch scrobbles.)</div>
			</>
		}

		// handle initial redux state
		if (scrobbles.done === false) {
			return null
		}

		nProgress.done()

		// can happen if the privacy was changed after the initial server page load
		if (scrobbles.private) {
			return <>
				{header}
				<div className="info">(This user's scrobbles are private.)</div>
			</>
		}

		if (scrobbles.items.length === 0) {
			return <>
				{top}
				<div className="info">({self ? "You haven't" : "This user hasn't"} scrobbled {mode != Mode.All ? "matching " : ""}songs yet.)</div>
			</>
		}

		const itemsToShow = scrobbles.items.slice(0, endIdx);

		switch (mode) {
			case Mode.All:
			case Mode.Loved: {
				return <>
					{top}
					{searchBox}
					<div className="songs">
						<Songs
							songs={itemsToShow}
							artworkBaseURL={artworkBaseURL}
							albumCentric={false}
							more={scrobbles.total! - itemsToShow.length}
							// "showing all songs that are available on the client" && "more number of songs present for the user"
							showMoreCard={(itemsToShow.length === scrobbles.items.length) && (scrobbles.total! > scrobbles.items.length)}
							showDates={true}
							now={() => new Date()}
							onSongClick={onSongClick}
						/>
					</div>
				</>
			}

			case Mode.Color: {
				return <>
					{top}
					<div className="songs">
						<Songs
							songs={itemsToShow}
							artworkBaseURL={artworkBaseURL}
							albumCentric={true}
							showMoreCard={false}
							showDates={false}
							onSongClick={onSongClick}
						/>
					</div>
				</>
			}

			default: {
				assertExhaustive(mode)
			}
		}
	}
