import * as React from "react";
import { Song } from "../src/shared"
import { displayString as dateDisplayString } from "../src/time"

export class SongCard extends React.Component<{song: Song, artworkBaseURL: string, now: Date}, {}> {
  private static readonly defaultArtworkPath = "/static/img/default-artwork.jpeg"
  private trackLinkAreaElem: HTMLDivElement|null = null

  componentDidMount() {
    // SO says onclick on an element enables :hover on iOS
    this.trackLinkAreaElem!.setAttribute("onclick", "");
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
    return <div className="scaleArea">
      {this.pict()}
      {this.meta()}
    </div>
  }

  private meta() {
    const s = this.props.song
    return <div className="meta" title={this.tooltip()}>
      <div className="title">
        <span className="titleContent">{s.title}</span>
        {s.loved && <span className="love"></span>}
      </div>
      <div className="other">
        {s.artistName && <span className="artist">{s.artistName}</span>}
      </div>
      {s.lastPlayed && <time className="date">{dateDisplayString(new Date(s.lastPlayed * 1000), this.props.now)}</time>}
    </div>
  }

  private pict() {
    let imgStyles = this.artworkURL() ?
      {backgroundImage: `url(${this.artworkURL()})`} :
      {backgroundColor: "#fff"}

    return <div className="pict" style={imgStyles} onClick={() => {
      console.log("clicked pict")
    }}>
      {this.props.song.trackViewURL && this.trackLinkArea()}
    </div>
  }

  private trackLinkArea() {
    return <a href={this.props.song.trackViewURL} title={this.props.song.trackViewURL} target="_blank">
      <div className="trackLinkArea" ref={r => {this.trackLinkAreaElem = r}}>
        <div className="trackLink"></div>
      </div>
    </a>
  }

  render() {
    return <div className="SongCard">
      {this.card()}
    </div>
  }
}

