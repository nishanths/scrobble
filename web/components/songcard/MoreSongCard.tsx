import React from "react"
import { MorePicture } from "./util"

interface MoreSongCardProps {
  more: number
}

export const MoreSongCard: React.StatelessComponent<MoreSongCardProps> = ({ more }) => {
  return <div className="SongCard MoreSongCard">
    <div className="scaleArea">
      <MorePicture more={more} />
    </div>
  </div>
}
