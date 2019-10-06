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

type UProps = UArgs & {
  wnd: Window
}

enum Mode {
  All, Loved
}

const controlValues = ["All", "Loved"] as const

const controlValueForMode = (m: Mode): typeof controlValues[number] => {
  switch (m) {
    case Mode.All: return "All"
    case Mode.Loved: return "Loved"
    default: throw assertExhaustive(m)
  }
}

const modeFromControlValue = (v: typeof controlValues[number]): Mode => {
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
  const [mode, modeRef, setMode] = useStateRef(Mode.All)

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
        afterChange={(v) => { setMode(modeFromControlValue(v)) }}
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

//   private static controlValueForMode(m: Mode): string {
//     switch (m) {
//       case Mode.All:
//         return "All"
//       case Mode.Loved:
//         return "Loved"
//     }
//     assertExhaustive(m)
//   }

//   private static modeFromControlValue(v: string): Mode {
//     switch (v) {
//       case "All": return Mode.All
//       case "Loved": return Mode.Loved
//       default: return Mode.All // fallback
//     }
//   }

//   // Divisble by 2, 3, and 4. This is appropriate because these are the number
//   // of cards typically displayed per row. Using such a number ensures that
//   // the last row isn't an incomplete row.
//   private static readonly moreIncrement = 36

//   private static determineNextEndIdx(idx: number, nSongs: number): number {
//     // increment, but make sure we don't go over the number of songs itself
//     let b = Math.min(idx + U.moreIncrement, nSongs)
//     // if there aren't sufficient songs left for the next time, just include
//     // them now
//     return nSongs - b < U.moreIncrement ? nSongs : b;
//   }

//   constructor(props: UProps) {
//     super(props)
//     this.state = {
//       fetched: false,
//       songs: [],
//       private: true,
//       endIdx: 0,
//       mode: Mode.All,
//     }
//   }

//   componentDidMount() {
//     NProgress.configure({ showSpinner: false, minimum: 0.1, trickleSpeed: 25, speed: 500 })
//     NProgress.start()

//     const leeway = 250
//     window.addEventListener("scroll", () => {
//       if ((window.innerHeight + window.pageYOffset) >= (document.body.offsetHeight - leeway)) {
//         this.showMore()
//       }
//     })

//     this.setState({ mode: U.modeFromURL(window) })
//     this.fetchSongs()
//   }

//   private onControlToggled(v: string) {
//     const mode = U.modeFromControlValue(v)
//     this.setState({ mode }, () => {
//       window.history.pushState(null, "", U.urlFromMode(this.state.mode, this.props.profileUsername)) // TODO: gross side-effect in this function?
//     })
//   }

//   private fetchSongs() {
//     let success = false // TODO: learn `fetch`

//     fetch("/api/v1/scrobbled?username=" + this.props.profileUsername)
//       .then(r => {
//         if (r.status == 200) {
//           success = true
//           this.setState({ private: false })
//           return r.json()
//         }
//         if (r.status == 404) {
//           this.setState({ fetched: true, private: true })
//           return
//         }
//       }).then(r => {
//         if (!success) { return }
//         let songs = r as Song[];
//         this.setState({
//           fetched: true,
//           songs: songs,
//           endIdx: U.determineNextEndIdx(this.state.endIdx, songs.length),
//         })
//       }, err => {
//         console.error(err)
//       })
//   }

//   private header() {
//     return <Header username={this.props.profileUsername} signedIn={!!this.props.logoutURL} />
//   }

//   private showMore() {
//     this.setState(s => {
//       let newEndIdx = U.determineNextEndIdx(s.endIdx, this.songsForCurrentMode().length)
//       return { endIdx: Math.max(newEndIdx, s.endIdx) }
//     })
//   }

//   private songsForCurrentMode = (): Song[] => {
//     switch (this.state.mode) {
//       case Mode.All:
//         return this.state.songs
//       case Mode.Loved:
//         return this.state.songs.filter(s => s.loved)
//     }
//     assertExhaustive(this.state.mode)
//   }

//   render() {
//     if (!this.state.fetched) {
//       return <div>{this.header()}</div>
//     }

//     NProgress.done()

//     if (this.state.private) {
//       return <div>
//         {this.header()}
//         <div className="info">(This user's scrobbles are private.)</div>
//       </div>
//     }

//     if (this.state.songs.length == 0) {
//       return <div>
//         {this.header()}
//         <div className="info">({this.props.self ? "You haven't" : "This user hasn't"} scrobbled yet.)</div>
//       </div>
//     }

//     return <div>
//       {this.header()}
//       <div className="control">
//         <SegmentedControl
//           afterChange={this.onControlToggled.bind(this)}
//           values={["All", "Loved"]}
//           initialValue={U.controlValueForMode(this.state.mode)}
//         />
//       </div>
//       <div className="songs">
//         <Songs
//           songs={this.songsForCurrentMode().slice(0, this.state.endIdx)}
//           artworkBaseURL={this.props.artworkBaseURL}
//           now={() => new Date()}
//         />
//       </div>
//     </div>
//   }
