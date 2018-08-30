import * as React from "react";
import { UArgs, Song } from "../src/shared"
import { Header } from "./Header"
import "../scss/u.scss"

class Songs extends React.Component<{songs: Song[]}, {}> {
  render() {
    return <div>{this.props.songs.length}</div>
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
    return <Header username={this.props.profileUsername} logoutURL={this.props.logoutURL} />
  }

  private onMoreClick() {
    this.setState({
      endIdx: UsernamePage.determineNextEndIdx(this.state.endIdx, this.state.songs.length)
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

    if (this.state.private) {
      return <div>
        {this.header()}
        <div className="info">(This user's scrobbles are private.)</div>
      </div>
    }

    if (this.state.songs.length == 0) {
      return <div>
        {this.header()}
        <div className="info">({this.props.self ? "You haven't" : "This user hasn't"}) scrobbled yet.</div>
      </div>
    }

    return <div>
      {this.header()}
      <Songs songs={this.state.songs.slice(0, this.state.endIdx)}/>
      {this.state.endIdx < this.state.songs.length && <div className="more" onClick={this.onMoreClick.bind(this)}>More songs</div>}
    </div>
  }
}
