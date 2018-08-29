import * as React from "react";
import { UArgs, Song } from "../src/shared"
import { Header } from "./Header"
import "../scss/u.scss"

class Songs extends React.Component<{songs: Song[]}, {}> {
}

type UsernamePageProps = UArgs

interface UsernamePageState {
  fetched: boolean
  songs: Song[]
  private: boolean
  idx: number
}

export class UsernamePage extends React.Component<UsernamePageProps, UsernamePageState> {
  constructor(props: UsernamePageProps) {
    super(props)
    this.state = {
      fetched: false,
      songs: [],
      private: true,
      idx: 0,
    }
  }

  componentDidMount() {
    this.fetchSongs()
  }

  private fetchSongs() {
    let success = false

    fetch("https://" + this.props.host + "/api/v1/scrobbled?username=" + this.props.profileUsername, {method: "GET"})
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
        this.setState({fetched: true, songs: r as Song[]})
      }, err => {
        this.setState({fetched: false})
      })
  }

  header() {
    return <Header username={this.props.profileUsername} logoutURL={this.props.logoutURL} />
  }

  render() {
    if (!this.state.fetched) {
      return <div>{this.header()}</div>
    }

    if (this.state.private) {
      return <div>
        {this.header()}
        <div>This user's scrobbles are private.</div>
      </div>
    }

    return <div>
      {this.header()}
      <Songs songs={this.state.songs}/>
      <div>more</div>
    </div>
  }
}