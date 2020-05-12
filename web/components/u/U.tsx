import React, { useState, useEffect, useRef } from "react"
import { useSelector, useDispatch } from "react-redux";
import { RouteComponentProps, Redirect } from "react-router-dom";
import { UArgs, Song, NProgress } from "../../shared/types"
import { Mode, DetailKind, pathForMode, pathForColor, modeFromControlValue } from "./shared"
import { Header, ColorPicker, Top } from "./top"
import { trimPrefix, assertExhaustive, pathComponents, assert, hexDecode, hexEncode } from "../../shared/util"
import { Songs } from "../Songs"
import { SegmentedControl } from "../SegmentedControl"
import { CloseIcon } from "../CloseIcon"
import { Color } from "../colorpicker"
import { State } from "../../redux/types/u"
import { fetchAllScrobbles, fetchLovedScrobbles, fetchColorScrobbles } from "../../redux/actions/scrobbles"
import { fetchSong } from "../../redux/actions/song"
import { useStateRef } from "../../shared/hooks"
import "../../scss/u.scss"
import "../../scss/detail-modal.scss"
import 'react-responsive-modal/styles.css';
import { Modal } from 'react-responsive-modal';

type UProps = UArgs & {
  wnd: Window
  mode: Mode
  color?: Color
  detail?: {
    kind: DetailKind
    hexIdent: string // hex-encoded song ident
  }
} & RouteComponentProps;

declare var NProgress: NProgress

// U is the root component for the username page, e.g.,
// https://scrobble.littleroot.org/u/whatever.
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
  color,
  detail,
  history,
}) => {
  // Divisble by 2, 3, and 4. This is appropriate because these are the number
  // of cards typically displayed per row. Using such a number ensures that
  // the last row isn't an incomplete row.
  const moreIncrement = 36;
  const limit = 504; // moreIncrement * 14

  const dispatch = useDispatch()
  const [endIdx, endIdxRef, setEndIdx] = useStateRef(0)
  const [lastColor, setLastColor] = useState(color) // save latest color when switching between other modes
  useEffect(() => {
    if (mode === Mode.Color) {
      setLastColor(color)
    }
  }, [color])

  const onControlChange = (newMode: Mode) => {
    NProgress.done()
    let u: string;
    switch (newMode) {
      case Mode.All:
      case Mode.Loved:
        u = "/u/" + profileUsername + pathForMode(newMode)
        break
      case Mode.Color:
        u = "/u/" + profileUsername + pathForMode(newMode) + pathForColor(lastColor)
        break
      default:
        assertExhaustive(newMode)
    }
    history.push(u)
  }

  const onColorChange = (newColor: Color) => {
    NProgress.done()
    console.log(newColor)
    assert(mode === Mode.Color, "mode should be Color")
    history.push("/u/" + profileUsername + pathForMode(mode) + pathForColor(newColor))
  }

  const nextEndIdx = (currentEndIdx: number, total: number): number => {
    // increment, but don't go over the number of items itself
    const b = Math.min(currentEndIdx + moreIncrement, total)
    // if there aren't sufficient items left for the next time, just include them now
    return total - b < moreIncrement ? total : b;
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

  // song detail redux state
  const detailSong = useSelector((s: State) => {
    if (detail === undefined) {
      return null
    }
    const key = hexDecode(detail.hexIdent)
    return s.songs.getOrDefault(key)
  })
  const detailSongRef = useRef(detailSong)
  useEffect(() => { detailSongRef.current = detailSong }, [detailSong])

  // initial end idx
  useEffect(() => {
    const s = scrobblesRef.current
    if (s === null) {
      return
    }
    const e = s.error === false ? nextEndIdx(0, s.items.length) : 0
    setEndIdx(e)
  }, [scrobbles, mode])

  // fetch scrobbles
  useEffect(() => {
    if (detail !== undefined) {
      return
    }

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
  }, [profileUsername, mode, color, detail])

  // fetch song detail
  useEffect(() => {
    if (detail === undefined) {
      return
    }
    const s = detailSongRef.current
    if (s === null || (s.done === false && s.fetching === false) || s.error === true) {
      dispatch(fetchSong(profileUsername, hexDecode(detail.hexIdent)))
    }
  }, [profileUsername, detail])

  // infinite scroll
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

  const header = Header(profileUsername, logoutURL)
  const colorPicker = ColorPicker(color, onColorChange)
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

  NProgress.configure({ showSpinner: false, minimum: 0.1, trickleSpeed: 150, speed: 500 })

  if (scrobbles.fetching === true) {
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

  // handle initial redux state
  if (scrobbles.done === false) {
    return null
  }

  NProgress.done()

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
        <div className="songs">
          <Songs
            songs={itemsToShow}
            artworkBaseURL={artworkBaseURL}
            albumCentric={false}
            more={scrobbles.total! - itemsToShow.length}
            // "showing all songs that are available on the client" && "more number of songs present for the user "
            showMoreCard={(itemsToShow.length === scrobbles.items.length) && (scrobbles.total! > scrobbles.items.length)}
            showDates={true}
            now={() => new Date()}
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
          />
        </div>
      </>
    }

    default: {
      assertExhaustive(mode)
    }
  }
}

