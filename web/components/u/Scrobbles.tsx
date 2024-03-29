import { Index } from "flexsearch";
import React, { useEffect, useRef, useState } from "react";
import { useHotkeys } from "react-hotkeys-hook";
import { useDispatch, useSelector } from "react-redux";
import { RouteComponentProps } from "react-router-dom";
import { setLastColor, setLastScrobblesEndIdx, setLastScrobblesScrollY, setLastSearch } from "../../redux/actions/last";
import { fetchAllScrobbles, fetchColorScrobbles, fetchLovedScrobbles } from "../../redux/actions/scrobbles";
import { State } from "../../redux/types/u";
import { useStateRef } from "../../shared/hooks";
import { NProgress, Song } from "../../shared/types";
import { assert, assertExhaustive, copyMap, debounce } from "../../shared/util";
import { Songs } from "../Songs";
import { Color } from "../colorpicker";
import { SearchBox } from "../searchbox";
import { Doc, createIndex, hasActiveSearch, indexID } from "./search";
import { DetailKind, Mode, fullPath, modeFromControlValue } from "./shared";
import { ColorPicker, Header, SegmentedControl, Top } from "./top";

// Divisble by 2, 3, and 4. This is appropriate because these are the number
// of cards typically displayed per row. Using such a number ensures that
// the last row isn't an incomplete row.
const moreIncrement = 36;

const nextEndIdx = (currentEndIdx: number, total: number): number => {
	// increment, but don't go over the number of items itself
	const b = Math.min(currentEndIdx + moreIncrement, total)
	// if there aren't sufficient items left for the next time, just include them now
	return total - b < moreIncrement ? total : b;
}

type History = RouteComponentProps["history"]

const searchPlaceholder = "Filter by album, artist, or song title"

// required data to enable searching for a given mode.
type IndexForMode = {
	songIdents: Map<string, Song> // song ident -> Song
	searchIndex: Index<Doc>       // actual search index
}

type ScrobblesMode = Mode.All | Mode.Loved | Mode.Color

export const Scrobbles: React.FC<{
	profileUsername: string
	signedIn: boolean
	artworkBaseURL: string
	private: boolean
	self: boolean
	mode: ScrobblesMode
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
		const searchInputElem = React.createRef<HTMLInputElement>()
		useHotkeys("cmd+k", () => {
			if (searchInputElem.current !== null) {
				searchInputElem.current.focus()
			}
		}, [searchInputElem])

		const dispatch = useDispatch()
		const last = useSelector((s: State) => s.last)

		const modeRef = useRef(mode)
		useEffect(() => { modeRef.current = mode }, [mode])

		const [endIdx, endIdxRef, setEndIdx] = useStateRef(last.scrobblesEndIdx || 0)
		const [scrollY, setScrollY] = useState(wnd.pageYOffset)

		const shouldUpdateScrollTo = () => last.scrobblesScrollY !== undefined
		useEffect(() => {
			if (shouldUpdateScrollTo()) {
				wnd.scrollTo({ top: last.scrobblesScrollY })
				dispatch(setLastScrobblesScrollY(undefined))
			}
		})

		const [searchValue, setSearchValue] = useState(last.search)
		const [indexesByMode, indexesByModeRef, setIndexesByMode] = useStateRef<Map<Mode.All | Mode.Loved, IndexForMode>>(new Map())
		const [filteredSongs, filteredSongsRef, setFilteredSongs] = useStateRef<Song[] | null>(null)

		const onControlChange = (newMode: Mode) => {
			nProgress.done()
			dispatch(setLastScrobblesEndIdx(0))
			setSearchValue("")
			dispatch(setLastSearch(""))

			history.push(fullPath(profileUsername, newMode, last.color, undefined, undefined))
		}

		const onColorChange = (newColor: Color) => {
			nProgress.done()
			dispatch(setLastScrobblesEndIdx(0))
			dispatch(setLastColor(newColor))
			assert(mode === Mode.Color, "mode should be Color")
			history.push(fullPath(profileUsername, mode, newColor, undefined, undefined))
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

			history.push(fullPath(profileUsername, mode, color, undefined, { kind, songIdent: s.ident }))
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
						dispatch(fetchAllScrobbles(profileUsername, null))
					}
					break
				}
				case Mode.Loved: {
					if ((s!.done === false && s!.fetching === false) || s!.error === true) {
						dispatch(fetchLovedScrobbles(profileUsername))
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
				const f = filteredSongsRef.current

				const items = hasActiveSearch(searchValue) ?
					(f) :
					(s && s.items)
				if (items === null) {
					return
				}

				const leeway = 250
				if ((wnd.innerHeight + y) >= (wnd.document.body.offsetHeight - leeway)) {
					const newEnd = nextEndIdx(endIdxRef.current, items.length)
					const e = Math.max(newEnd, endIdxRef.current)
					if (e > endIdxRef.current) {
						setEndIdx(e)
					}
				}
			}

			wnd.addEventListener("scroll", f)
			return () => { wnd.removeEventListener("scroll", f) }
		}, [scrobbles, mode, searchValue])

		// search
		useEffect(() => {
			if (mode !== Mode.All && mode !== Mode.Loved) {
				return
			}

			if (scrobbles === null) {
				return
			}
			if (scrobbles.fetching === true) {
				return
			}
			if (scrobbles.error === true) {
				return
			}
			// handle initial redux state
			if (scrobbles.done === false) {
				return
			}

			const indexForMode = indexesByMode.get(mode)

			if (indexForMode === undefined) {
				const searchIndex = createIndex()
				const songIdents = new Map<string, Song>()

				for (const s of scrobbles.items) {
					const d: Doc = {
						id: indexID(s.ident),
						songIdent: s.ident,
						title: s.title,
						artist: s.artistName,
						album: s.albumTitle,
					}
					searchIndex.add(d)
					songIdents.set(s.ident, s)
				}

				const updated = copyMap(indexesByMode).set(mode, { searchIndex, songIdents })
				setIndexesByMode(updated)
			}
		}, [mode, scrobbles])

		const [debouncedSearch,] = useState(() => debounce((v) => {
			const m = modeRef.current
			if (m !== Mode.All && m !== Mode.Loved) {
				return
			}
			const i = indexesByModeRef.current.get(m)
			if (i === undefined) {
				return
			}
			const query = v.trim().length === 0 ? "" : v
			i.searchIndex.search({ query }, (docs: Doc[]) => {
				setFilteredSongs(docs.map(doc => i.songIdents.get(doc.songIdent)!))
			})
		}, 300))

		useEffect(() => {
			debouncedSearch(searchValue)
		}, [searchValue])

		const onSearchValueChange = (v: string): void => {
			setSearchValue(v)
			dispatch(setLastSearch(v))
			setFilteredSongs(null)
		}

		// ... render ...

		const header = Header(profileUsername, signedIn, true)
		const colorPicker = ColorPicker(color, onColorChange)
		const segmentedControl = SegmentedControl(mode, (v) => { onControlChange(modeFromControlValue(v)) })
		const top = Top(header, segmentedControl, colorPicker, mode)

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

		const noMatchingScrobbles = <>
			{top}
			<div className="info">({self ? "You haven't" : "This user hasn't"} scrobbled {mode != Mode.All ? "matching " : ""}songs yet.)</div>
		</>

		const searchBox = <div className="searchBox">
			<SearchBox
				ref={searchInputElem}
				value={searchValue}
				onChange={onSearchValueChange}
				placeholder={searchPlaceholder}
			/>
		</div>

		switch (mode) {
			case Mode.All:
			case Mode.Loved: {
				if (hasActiveSearch(searchValue)) {
					if (filteredSongs === null) {
						// waiting for search results state
						return <>
							{top}
							{searchBox}
						</>
					}
					if (filteredSongs.length === 0) {
						return <>
							{top}
							{searchBox}
							<div className="info">(No matching songs.)</div>
						</>
					}
				} else {
					// should render all scrobbles; no filtering
					if (scrobbles.items.length === 0) {
						return noMatchingScrobbles
					}
				}

				const itemsToShow = hasActiveSearch(searchValue) ?
					filteredSongs!.slice(0, endIdx) :
					scrobbles.items.slice(0, endIdx)

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
				if (scrobbles.items.length === 0) {
					return noMatchingScrobbles
				}
				const itemsToShow = scrobbles.items.slice(0, endIdx)
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
