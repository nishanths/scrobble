import React, { useState, useEffect, useRef } from "react"
import { useSelector, useDispatch } from "react-redux";
import { RouteComponentProps, Redirect } from "react-router-dom";
import { UArgs, Song } from "../shared/types"
import { trimPrefix, assertExhaustive, pathComponents } from "../shared/util"
import { Header } from "./Header"
import { Songs } from "./Songs"
import { Artworks } from "./Artworks"
import { SegmentedControl } from "./SegmentedControl"
import { Color, ColorPicker } from "./colorpicker"
import { State } from "../redux/types/u"
import { fetchAllScrobbles, fetchLovedScrobbles, fetchColorScrobbles } from "../redux/actions/scrobbles"
import { useStateRef } from "../shared/hooks"
import "../scss/u.scss"

declare var NProgress: {
  start(): void
  done(): void
  configure(opts: { [k: string]: any }): void
}

export enum Mode {
  All, Loved, Color
}

const controlValues = ["All", "Loved", "By color"] as const

type ControlValue = typeof controlValues[number]

const controlValueForMode = (m: Mode): ControlValue => {
  switch (m) {
    case Mode.All: return "All"
    case Mode.Loved: return "Loved"
    case Mode.Color: return "By color"
    default: assertExhaustive(m)
  }
}

const modeFromControlValue = (v: ControlValue): Mode => {
  switch (v) {
    case "All": return Mode.All
    case "Loved": return Mode.Loved
    case "By color": return Mode.Color
    default: assertExhaustive(v)
  }
}

const pathForMode = (m: Mode): string => {
  switch (m) {
    case Mode.All: return ""
    case Mode.Loved: return "/loved"
    case Mode.Color: return "/color"
  }
  assertExhaustive(m)
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
  private: priv,
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
  const [color, colorRef, setColor] = useStateRef<Color | undefined>(undefined)

  const onControlChange = (newMode: Mode) => { history.push("/u/" + profileUsername + pathForMode(newMode)) }

  const scrobbles = useSelector((s: State) => {
    switch (mode) {
      case Mode.All: return s.allScrobbles
      case Mode.Loved: return s.lovedScrobbles
      case Mode.Color: return color !== undefined ? s.colorScrobbles.get(color)! : null
    }
    throw assertExhaustive(mode)
  })
  const scrobblesRef = useRef(scrobbles)
  useEffect(() => { scrobblesRef.current = scrobbles }, [scrobbles])

  const nextEndIdx = (currentEndIdx: number, total: number): number => {
    // increment, but don't go over the number of items itself
    const b = Math.min(currentEndIdx + moreIncrement, total)
    // if there aren't sufficient items left for the next time, just include them now
    return total - b < moreIncrement ? total : b;
  }

  useEffect(() => {
    const s = scrobblesRef.current
    if (s === null) {
      return
    }
    const e = s.error === false ? nextEndIdx(0, s.items.length) : 0
    setEndIdx(e)
  }, [scrobbles, mode])

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
        if (colorRef.current === undefined) {
          break
        }
        if (s === null || (s!.done === false && s!.fetching === false) || s!.error === true) {
          dispatch(fetchColorScrobbles(colorRef.current, profileUsername))
        }
        break
      }
      default: {
        throw assertExhaustive(mode)
      }
    }
  }, [profileUsername, mode, color])

  useEffect(() => {
    const f = () => {
      const s = scrobblesRef.current
      if (s === null) {
        return
      }

      const leeway = 250

      if ((wnd.innerHeight + wnd.pageYOffset) >= (wnd.document.body.offsetHeight - leeway)) {
        const newEnd = nextEndIdx(endIdxRef.current, s.items.length)
        const e = Math.max(newEnd, endIdxRef.current)
        setEndIdx(e)
      }
    }

    wnd.addEventListener("scroll", f)
    return () => { wnd.removeEventListener("scroll", f) }
  }, [scrobbles, mode])

  // ... render ...

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

  const colorPicker = <div className="colorPicker">
    <ColorPicker initialSelection={color} prompt="Pick a color to see artwork of that color." afterSelect={(c) => { setColor(c) }} />
  </div>


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
      {colorPicker}
    </>
  }

  const s = scrobbles!

  NProgress.configure({ showSpinner: false, minimum: 0.1, trickleSpeed: 150, speed: 500 })

  if (s.done === false) {
    NProgress.start()
    return <>{top}</>
  }

  if (s.error === true) {
    NProgress.done()
    return <>
      {header}
      <div className="info">(Failed to fetch scrobbles.)</div>
    </>
  }

  NProgress.done()

  // can happen if the privacy was changed after the initial server page load
  if (s.private) {
    return <>
      {header}
      <div className="info">(This user's scrobbles are private.)</div>
    </>
  }

  if (s.items.length === 0) {
    return <>
      {header}
      <div className="info">({self ? "You haven't" : "This user hasn't"} scrobbled yet.)</div>
    </>
  }

  const itemsToShow = s.items.slice(0, endIdx);

  switch (mode) {
    case Mode.All:
    case Mode.Loved: {
      return <>
        {top}
        <div className="songs">
          <Songs
            songs={itemsToShow as Song[]}
            more={s.total! - itemsToShow.length}
            // "showing all songs that are available on the client" && "more number of songs present for the user "
            showMore={(itemsToShow.length === s.items.length) && (s.total! > s.items.length)}
            artworkBaseURL={artworkBaseURL}
            now={() => new Date()}
          />
        </div>
      </>
    }

    case Mode.Color: {
      return <>
        {top}
        {colorPicker}
        <div className="songs">
          <Artworks
            hashes={itemsToShow as string[]}
            artworkBaseURL={artworkBaseURL}
          />
        </div>
      </>
    }

    default: {
      assertExhaustive(mode)
    }
  }
}

