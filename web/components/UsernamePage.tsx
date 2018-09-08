import * as React from "react";
import { UArgs, Song } from "../src/shared"
import { Header } from "./Header"
import { SongCard } from "./SongCard"
import "../scss/u.scss"

declare var NProgress: {
  start(): void
  done(): void
  configure(opts: {[k: string]: any}): void
}

class Songs extends React.Component<{songs: Song[], artworkBaseURL: string}, {}> {
  private static key(s: Song): string {
    return [s.title, s.albumTitle, s.artistName, s.year].join("|")
  }

  render() {
    let now = new Date()
    return this.props.songs.map(s => {
      return <SongCard key={Songs.key(s)} song={s} artworkBaseURL={this.props.artworkBaseURL} now={now}/>
    })
  }
}

type UsernamePageProps = UArgs

interface UsernamePageState {
  fetched: boolean
  songs: Song[]
  private: boolean
  endIdx: number
}

export class UsernamePage extends React.Component<UsernamePageProps, UsernamePageState> {
  // Divisble by 2, 3, and 4. This is appropriate because these are the number
  // of cards typically displayed per row. Using such a number ensures that
  // the last row isn't an incomplete row.
  private static readonly moreIncrement = 48

  constructor(props: UsernamePageProps) {
    super(props)
    this.state = {
      fetched: false,
      songs: [],
      private: true,
      endIdx: 0,
    }
  }

  componentDidMount() {
    NProgress.configure({ showSpinner: false, minimum: 0.1, trickleSpeed: 25, speed: 500 })
    NProgress.start()
    const leeway = 200
    window.addEventListener("scroll", () => {
      if ((window.innerHeight + window.pageYOffset) >= (document.body.offsetHeight - leeway)) {
        this.showMore()
      }
    })
    this.fetchSongs()
  }

  private fetchSongs() {
    let success = false // TODO: learn to use `fetch`

    fetch("/api/v1/scrobbled?username=" + this.props.profileUsername)
      .then(r => {
        if (r.status == 200) {
          success = true
          this.setState({private: false})
          return r.json()
        }
        if (r.status == 404) {
          this.setState({fetched: true, private: true})
          return
        }
      }).then(r => {
        if (!success) { return }
        let s = r as Song[];
        this.setState({
          fetched: true,
          songs: s,
          endIdx: UsernamePage.determineNextEndIdx(this.state.endIdx, s.length),
        })
      }, err => {
        console.error(err)
      })
  }

  private header() {
    return <Header username={this.props.profileUsername} signedIn={!!this.props.logoutURL}/>
  }

  private showMore() {
    this.setState(s => {
      return {endIdx: UsernamePage.determineNextEndIdx(s.endIdx, s.songs.length)}
    })
  }

  private static determineNextEndIdx(idx: number, nSongs: number): number {
    // increment, but make sure we don't go over the number of songs itself
    let b = Math.min(idx + UsernamePage.moreIncrement, nSongs)
    // if there aren't sufficient songs left for the next time, just include
    // them now
    return nSongs - b < UsernamePage.moreIncrement ? nSongs : b;
  }

  render() {
    if (!this.state.fetched) {
      return <div>{this.header()}</div>
    }

    NProgress.done()

    if (this.state.private) {
      return <div>
        {this.header()}
        <div className="info">(This user's scrobbles are private.)</div>
      </div>
    }

    if (this.state.songs.length == 0) {
      return <div>
        {this.header()}
        <div className="info">({this.props.self ? "You haven't" : "This user hasn't"} scrobbled yet.)</div>
      </div>
    }

    return <div>
      {this.header()}
      <div className="songs">
        <Songs songs={this.state.songs.slice(0, this.state.endIdx)} artworkBaseURL={this.props.artworkBaseURL}/>
      </div>
    </div>
  }
}
