import * as React from "react";
import { Song } from "../src/shared"
import { displayString as dateDisplayString } from "../src/time"

export class SongCard extends React.Component<{song: Song, artworkBaseURL: string, now: Date}, {}> {
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

  render() {
    let s = this.props.song
    let imgStyles = this.artworkURL() ? {backgroundImage: `url(${this.artworkURL()})`} : {backgroundColor: "#fff"};

    return <div className="SongCard">
      <div className="scaleArea" ref={r => {this.scaleArea = r}}>
        <div className="pict" style={imgStyles}></div>
        <div className="meta" title={this.tooltip()}>
          <div className="title">{s.title}</div>
          <div className="other">
            {s.artistName && <span className="artist">{s.artistName}</span>}
          </div>
          {s.lastPlayed && <time className="date">{dateDisplayString(new Date(s.lastPlayed * 1000), this.props.now)}</time>}
        </div>
      </div>
    </div>
  }
}

