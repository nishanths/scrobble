import * as React from "react";
import { Song } from "../src/shared"
import { displayString as dateDisplayString } from "../src/time"

export class SongCard extends React.Component<{song: Song, artworkBaseURL: string, now: Date, modeToggled: () => void}, {}> {
  private static readonly defaultArtworkPath = "/static/img/default-artwork.jpeg"
  private scaleArea: HTMLDivElement|null = null

  componentDidMount() {
    // SO says onclick on an element enables :hover on iOS
    this.scaleArea!.setAttribute("onclick", "");
  }

  private artworkURL(): string {
    if (!this.props.song.artworkHash) {
      return ""
    }
    return this.props.artworkBaseURL + "/" + this.props.song.artworkHash
  }

  private tooltip(): string {
    let s = this.props.song;
    let tooltip = s.title
    if (s.artistName || s.albumTitle) {
      tooltip += "\n"
      if (s.artistName) { tooltip += s.artistName }
      if (s.artistName && s.artistName) { tooltip += " — " }
      if (s.albumTitle) { tooltip += s.albumTitle }
    }
    return tooltip
  }

  private card() {
    let s = this.props.song
    return <div className="scaleArea" ref={r => {this.scaleArea = r}}>
      {this.pict()}
      <div className="meta" title={this.tooltip()}>
        <div className="title">
          <span className="titleContent">{s.title}</span>
          {s.loved && <span className="love" onClick={this.props.modeToggled.bind(this)}></span>}
        </div>
        <div className="other">
          {s.artistName && <span className="artist">{s.artistName}</span>}
        </div>
        {s.lastPlayed && <time className="date">{dateDisplayString(new Date(s.lastPlayed * 1000), this.props.now)}</time>}
      </div>
    </div>
  }

  private pict() {
    let imgStyles = this.artworkURL() ? {backgroundImage: `url(${this.artworkURL()})`} : {backgroundColor: "#fff"}
    let e = <div className="pict" style={imgStyles}></div>
    return this.props.song.trackViewURL ? <a href={this.props.song.trackViewURL} target="_blank">{e}</a> : e
  }

  render() {
    return <div className="SongCard">
      {this.card()}
    </div>
  }
}

