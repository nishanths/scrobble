import React, { useState, useEffect } from "react"
import { useSelector, useDispatch } from "react-redux";
import { UArgs, Song } from "../shared/types"
import { trimPrefix, assertExhaustive, pathComponents } from "../shared/util"
import { Header } from "./Header"
import { Songs } from "./Songs"
import { SegmentedControl } from "./SegmentedControl"
import { State } from "../redux/reducers/u"
import "../scss/u.scss"

declare var NProgress: {
  start(): void
  done(): void
  configure(opts: { [k: string]: any }): void
}

type UProps = UArgs

interface UState {
  fetched: boolean
  songs: Song[]
  private: boolean
  endIdx: number
  mode: Mode
}

enum Mode {
  All, Loved
}

export const U: React.FC<UProps> = ({
  host,
  artworkBaseURL,
  profileUsername,
  logoutURL,
  account,
  self,
}) => {
  const [endIdx, setEndIdx] = useState(-1)
  const [mode, setMode] = useState(Mode.All)
  const dispatch = useDispatch()

  useEffect(() => {
    // dispatch(scrobblesRequest(profileUsername))
  })

  useSelector((s: State) => {
    console.log(s)
  })

  return <></>
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
//   private static readonly moreIncrement = 48

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
