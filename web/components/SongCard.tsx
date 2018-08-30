import * as React from "react";
import { Song } from "../src/shared"

export class SongCard extends React.Component<Song, {}> {
  render() {
    return <div>{this.props.title}</div>
  }
}
