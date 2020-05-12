import React from "react"
import { Song } from "../../shared/types"
import { SongCardProps } from "./SongCard"
import { LargeMeta, LargePicture } from "./util"

interface LargeSongCardProps extends SongCardProps { }

export const LargeSongCard: React.SFC<LargeSongCardProps> = ({
  song,
  artworkBaseURL,
  albumCentric,
  showDate,
  now,
}) => {
  return <div className="LargeSongCard">
    <LargePicture song={song} artworkBaseURL={artworkBaseURL} albumCentric={albumCentric} />
    <LargeMeta song={song} albumCentric={albumCentric} showDate={showDate} now={now} />
  </div>
}
