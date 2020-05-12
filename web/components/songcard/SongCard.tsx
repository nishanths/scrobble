import React from "react";
import { Song } from "../../shared/types"
import { Picture, Meta } from "./util"

export interface SongCardProps {
  song: Song // rendering degrades gracefully if song properties are missing
  artworkBaseURL: string
  // album-centric instead of song-centric
  // e.g. use album title instead of song title in the title area
  albumCentric: boolean

  showDate: boolean
  now?: () => Date // required if showDates is true
}

export const SongCard: React.SFC<SongCardProps> = ({
  song,
  artworkBaseURL,
  albumCentric,
  showDate,
  now
}) => {
  return <div className="SongCard">
    <div className="scaleArea">
      <Picture song={song} artworkBaseURL={artworkBaseURL} albumCentric={albumCentric} />
      <Meta song={song} albumCentric={albumCentric} showDate={showDate} now={now} />
    </div>
  </div>
}
