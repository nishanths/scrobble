import React, { useState, useEffect, useRef } from "react"
import { useSelector, useDispatch } from "react-redux";
import { RouteComponentProps, Redirect } from "react-router-dom";
import { UArgs, Song } from "../shared/types"
import { trimPrefix, assertExhaustive, pathComponents } from "../shared/util"
import { Header } from "./Header"
import { Songs } from "./Songs"
import { SegmentedControl } from "./SegmentedControl"
import { State } from "../redux/types/u"
import { fetchAllScrobbles, fetchLovedScrobbles } from "../redux/actions/scrobbles"
import { useStateRef } from "../shared/hooks"
import "../scss/u.scss"

declare var NProgress: {
  start(): void
  done(): void
  configure(opts: { [k: string]: any }): void
}

export enum Mode {
  All, Loved
}

const controlValues = ["All", "Loved"] as const

type ControlValue = typeof controlValues[number]

const controlValueForMode = (m: Mode): ControlValue => {
  switch (m) {
    case Mode.All: return "All"
    case Mode.Loved: return "Loved"
    default: assertExhaustive(m)
  }
}

const modeFromControlValue = (v: ControlValue): Mode => {
  switch (v) {
    case "All": return Mode.All
    case "Loved": return Mode.Loved
    default: assertExhaustive(v)
  }
}

const pathForMode = (m: Mode): string => {
  switch (m) {
    case Mode.All: return ""
    case Mode.Loved: return "/loved"
  }
  assertExhaustive(m)
}

const modeFromPath = (p: string): Mode => {
  switch (p) {
    case "":
    case "/":
    case "/all":
      return Mode.All
    case "/loved":
    case "/loved/":
      return Mode.Loved
    default:
      return Mode.All
  }
}

type UProps = UArgs & {
  wnd: Window
  mode: Mode
} & RouteComponentProps;

// U is the root component for the username page, e.g.,
// https://scrobble.allele.cc/u/whatever.
export const U: React.FC<UProps> = ({
  host,
  artworkBaseURL,
  profileUsername,
  logoutURL,
  account,
  self,
  wnd,
  mode,
  history,
}) => {
  // Divisble by 2, 3, and 4. This is appropriate because these are the number
  // of cards typically displayed per row. Using such a number ensures that
  // the last row isn't an incomplete row.
  const moreIncrement = 36;
  const limit = 504; // moreIncrement * 14

  const dispatch = useDispatch()
  const [endIdx, endIdxRef, setEndIdx] = useStateRef(0)

  const scrobbles = useSelector((s: State) => {
    switch (mode) {
      case Mode.All: return s.allScrobbles
      case Mode.Loved: return s.lovedScrobbles
    }
    throw assertExhaustive(mode)
  })
  const scrobblesRef = useRef(scrobbles)
  useEffect(() => { scrobblesRef.current = scrobbles }, [scrobbles])

  const onControlChange = (newMode: Mode): void => {
    history.push("/u/" + profileUsername + pathForMode(newMode))
  }

  const nextEndIdx = (currentEndIdx: number, totalSongs: number): number => {
    // increment, but don't go over the number of songs itself
    const b = Math.min(currentEndIdx + moreIncrement, totalSongs)
    // if there aren't sufficient songs left for the next time, just include them now
    return totalSongs - b < moreIncrement ? totalSongs : b;
  }

  useEffect(() => {
    const e = scrobbles.error === false ? nextEndIdx(0, scrobblesRef.current.songs.length) : 0
    setEndIdx(e)
  }, [scrobbles, mode])

  useEffect(() => {
    switch (mode) {
      case Mode.All:
        if (scrobblesRef.current.done === false || scrobblesRef.current.error === true) {
          dispatch(fetchAllScrobbles(profileUsername, limit))
        }
        break
      case Mode.Loved:
        if (scrobblesRef.current.done === false || scrobblesRef.current.error === true) {
          dispatch(fetchLovedScrobbles(profileUsername, limit))
        }
        break
      default:
        throw assertExhaustive(mode)
    }
  }, [profileUsername, mode])

  useEffect(() => {
    const f = () => {
      const leeway = 250
      if ((wnd.innerHeight + wnd.pageYOffset) >= (wnd.document.body.offsetHeight - leeway)) {
        const newEnd = nextEndIdx(endIdxRef.current, scrobblesRef.current.songs.length)
        const e = Math.max(newEnd, endIdxRef.current)
        setEndIdx(e)
      }
    }
    wnd.addEventListener("scroll", f)
    return () => { wnd.removeEventListener("scroll", f) }
  }, [scrobbles, mode])

  const header = <Header username={profileUsername} signedIn={!!logoutURL} />

  const top = <>
    {header}
    <div className="control">
      <SegmentedControl
        afterChange={(v) => { onControlChange(modeFromControlValue(v)) }}
        values={controlValues}
        initialValue={controlValueForMode(mode)}
      />
    </div>
  </>

  // ... render ...

  NProgress.configure({ showSpinner: false, minimum: 0.1, trickleSpeed: 150, speed: 500 })

  if (scrobbles.done === false) {
    NProgress.start()
    return <>{top}</>
  }

  if (scrobbles.error === true) {
    NProgress.done()
    return <>
      {header}
      <div className="info">(Failed to fetch scrobbles.)</div>
    </>
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

  const songsToShow = scrobbles.songs.slice(0, endIdx);

  return <>
    {top}
    <div className="songs">
      <Songs
        songs={songsToShow}
        more={scrobbles.total - songsToShow.length}
        // "showing all songs that are available on the client" && "more number of songs present for the user "
        showMore={(songsToShow.length === scrobbles.songs.length) && (scrobbles.total > scrobbles.songs.length)}
        artworkBaseURL={artworkBaseURL}
        now={() => new Date()}
      />
    </div>
  </>
}
