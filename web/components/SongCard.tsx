import * as React from "react";
import { Song } from "../src/shared"

export class SongCard extends React.Component<{song: Song, artworkBaseURL: string}, {}> {
  private static readonly defaultArtworkPath = "/static/img/default-artwork.jpeg"

  private artworkURL(): string {
    if (!this.props.song.artworkHash) {
      return SongCard.defaultArtworkPath
    }
    return this.props.artworkBaseURL + "/" + this.props.song.artworkHash
  }

  render() {
    let s = this.props.song;

    return <div className="SongCard">
      <div>
        <img src={this.artworkURL()}></img>
      </div>
      <div>{s.title}</div>
      <div>
        <span>{s.artistName}</span>
        {s.artistName && s.albumTitle && <span>&nbsp;—&nbsp;</span>}
        <span>{s.albumTitle}</span>
      </div>
    </div>
  }
}
