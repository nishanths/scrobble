import React, { useState, useEffect, useRef } from "react"
import { useSelector, useDispatch } from "react-redux";
import { RouteComponentProps } from "react-router-dom";
import { UArgs, NProgress } from "../../shared/types"
import { Mode, DetailKind, pathForMode, pathForColor, modeFromControlValue } from "./shared"
import { Header, ColorPicker, Top } from "./top"
import { Scrobbles } from "./Scrobbles"
import { Detail } from "./Detail"
import { assertExhaustive, hexDecode, assert } from "../../shared/util"
import { Color } from "../colorpicker"
import { State } from "../../redux/types/u"
import { fetchAllScrobbles, fetchLovedScrobbles, fetchColorScrobbles } from "../../redux/actions/scrobbles"
import { fetchSong } from "../../redux/actions/song"
import { useStateRef } from "../../shared/hooks"
import "../../scss/u.scss"

type UProps = UArgs & {
  wnd: Window
  mode: Mode
  color?: Color
  detail?: {
    kind: DetailKind
    hexIdent: string // hex-encoded song ident
  }
} & RouteComponentProps;

declare const NProgress: NProgress

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

// U is the root component for the username page, e.g.,
// https://scrobble.littleroot.org/u/whatever.
export const U: React.FC<UProps> = ({
  artworkBaseURL,
  profileUsername,
  logoutURL,
  self,
  private: priv,
  wnd,
  mode,
  color,
  detail,
  history,
}) => {
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
    const song = detailSongRef.current
    if (song === null || (song.done === false && song.fetching === false) || song.error === true) {
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

  if (detail === undefined) {
    return <Scrobbles
      scrobbles={scrobbles}
      artworkBaseURL={artworkBaseURL}
      endIdx={endIdx}
      private={priv}
      self={self}
      mode={mode}
      color={color}
      header={header}
      top={top}
      nProgress={NProgress}
    />
  } else {
    return <Detail
      song={detailSong!}
      profileUsername={profileUsername}
      mode={mode}
      color={color}
      private={priv}
      self={self}
      detailKind={detail.kind}
      nProgress={NProgress}
      history={history}
    />
  }
}

