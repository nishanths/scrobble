import React from "react"
import { Song } from "../../shared/types"
import { LargeMeta, LargePicture } from "./util"

interface LargeSongCardProps {
  song: Song
  artworkBaseURL: string
  albumCentric: boolean
  now: () => Date
}

export const LargeSongCard: React.SFC<LargeSongCardProps> = ({
  song,
  artworkBaseURL,
  albumCentric,
  now,
}) => {
  return <div className="LargeSongCard">
    <LargePicture song={song} artworkBaseURL={artworkBaseURL} />
    <LargeMeta song={song} albumCentric={albumCentric} now={now} />
  </div>
}
