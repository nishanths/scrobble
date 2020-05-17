import React from "react"
import { Song } from "../../shared/types"
import { LargeMeta, LargePicture } from "./util"
import { Loupe } from 'loupe-js'

interface LargeSongCardProps {
  song: Song
  artworkBaseURL: string
  albumCentric: boolean
  now: () => Date
  loupe: Loupe
}

export const LargeSongCard: React.SFC<LargeSongCardProps> = ({
  song,
  artworkBaseURL,
  albumCentric,
  now,
  loupe,
}) => {
  return <div className="LargeSongCard">
    <LargePicture song={song} artworkBaseURL={artworkBaseURL} loupe={loupe} />
    <LargeMeta song={song} albumCentric={albumCentric} now={now} />
  </div>
}
