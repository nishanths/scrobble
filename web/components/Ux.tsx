import React, { useState, useEffect, useRef } from "react"
import { useSelector, useDispatch } from "react-redux";
import { UArgs, Song } from "../shared/types"
import { trimPrefix, assertExhaustive, pathComponents } from "../shared/util"
import { Header } from "./Header"
import { Songs } from "./Songs"
import { SegmentedControl } from "./SegmentedControl"
import { State } from "../redux/types/u"
import { fetchScrobbles } from "../redux/actions/scrobbles"
import { useStateRef } from "../shared/hooks"
import "../scss/u.scss"

declare var NProgress: {
  start(): void
  done(): void
  configure(opts: { [k: string]: any }): void
}

enum Mode {
  All, Loved
}

const controlValues = ["All", "Loved"] as const

type ControlValue = typeof controlValues[number]

const controlValueForMode = (m: Mode): ControlValue => {
  switch (m) {
    case Mode.All: return "All"
    case Mode.Loved: return "Loved"
    default: throw assertExhaustive(m)
  }
}

const modeFromControlValue = (v: ControlValue): Mode => {
  switch (v) {
    case "All": return Mode.All
    case "Loved": return Mode.Loved
    default: throw assertExhaustive(v)
  }
}

const songsForMode = (m: Mode, s: Song[]): Song[] => {
  switch (m) {
    case Mode.All:
      return s
    case Mode.Loved:
      return s.filter(song => song.loved)
  }
  assertExhaustive(m)
}

type UProps = UArgs & {
  wnd: Window
}

// U is the root component for the username page, e.g.,
// https://scrobble.allele.cc/u/whatever.
export const U: React.FC<UProps> = ({
  wnd,
  host,
  artworkBaseURL,
  profileUsername,
  logoutURL,
  account,
  self,
}) => {
  // Divisble by 2, 3, and 4. This is appropriate because these are the number
  // of cards typically displayed per row. Using such a number ensures that
  // the last row isn't an incomplete row.
  const MoreIncrement = 36;
  const dispatch = useDispatch()
  const header = <Header username={profileUsername} signedIn={!!logoutURL} />;

  const [endIdx, endIdxRef, setEndIdx] = useStateRef(0)
  const [mode, modeRef, setMode] = useStateRef(Mode.All) // TODO: from URL

  const scrobbles = useSelector((s: State) => s.scrobbles)
  const scrobblesRef = useRef(scrobbles)
  useEffect(() => { scrobblesRef.current = scrobbles }, [scrobbles])

  const nextEndIdx = (currentEndIdx: number, totalSongs: number): number => {
    // increment, but don't go over the number of songs itself
    const b = Math.min(currentEndIdx + MoreIncrement, totalSongs)
    // if there aren't sufficient songs left for the next time, just include them now
    return totalSongs - b < MoreIncrement ? totalSongs : b;
  }

  const initEnd = useRef(false)

  useEffect(() => {
    if (initEnd.current === true) { return }
    initEnd.current = false
    setEndIdx(scrobbles.error === false ? nextEndIdx(0, scrobbles.songs.length) : 0)
  }, [scrobbles])

  useEffect(() => {
    NProgress.configure({ showSpinner: false, minimum: 0.1, trickleSpeed: 25, speed: 500 })
  }, [])

  useEffect(() => {
    dispatch(fetchScrobbles(profileUsername))
  }, [])

  useEffect(() => {
    const f = () => {
      const leeway = 250
      if ((wnd.innerHeight + wnd.pageYOffset) >= (wnd.document.body.offsetHeight - leeway)) {
        const newEnd = nextEndIdx(endIdxRef.current, songsForMode(modeRef.current, scrobblesRef.current.songs).length)
        setEndIdx(Math.max(newEnd, endIdxRef.current))
      }
    }
    wnd.addEventListener("scroll", f)
    return () => { wnd.removeEventListener("scroll", f) }
  }, [])

  // ... render ...

  if (scrobbles.error === true) {
    return <>
      {header}
      <div className="info">(Failed to fetch scrobbles.)</div>
    </>
  }

  if (scrobbles.fetching) {
    NProgress.start()
    return <>{header}</>
  }

  NProgress.done()

  if (scrobbles.private) {
    return <>
      {header}
      <div className="info">(This user's scrobbles are private.)</div>
    </>
  }

  if (scrobbles.songs.length === 0) {
    return <>
      {header}
      <div className="info">({self ? "You haven't" : "This user hasn't"} scrobbled yet.)</div>
    </>
  }

  return <>
    {header}
    <div className="control">
      <SegmentedControl
        afterChange={(v) => { setMode(modeFromControlValue(v)) }} // TODO: update URL
        values={controlValues}
        initialValue={controlValueForMode(mode)}
      />
    </div>
    <div className="songs">
      <Songs
        songs={songsForMode(mode, scrobbles.songs).slice(0, endIdx)}
        artworkBaseURL={artworkBaseURL}
        now={() => new Date()}
      />
    </div>
  </>
}


// private static modeFromURL(wnd: Window): Mode {
//     let components = pathComponents(wnd.location.pathname)
//     if (components.length < 3 || // /u/username,
//       components[2] != "loved" // /u/username/song, /u/username/<gibberish>
//     ) {
//       return Mode.All
//     }
//     return Mode.Loved // /u/username/loved, /u/username/loved/song, /u/username/loved/<gibberish>
//   }

//   private static urlFromMode(m: Mode, username: string): string {
//     switch (m) {
//       case Mode.All: return "/u/" + username
//       case Mode.Loved: return "/u/" + username + "/loved"
//     }
//     assertExhaustive(m)
//   }
