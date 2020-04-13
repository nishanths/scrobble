import React, { useEffect, useRef } from "react";
import { Song } from "../shared/types"
import { dateDisplay } from "../shared/time"

interface MoreSongCardProps {
  more: number
}

export const MoreSongCard: React.StatelessComponent<MoreSongCardProps> = ({ more }) => {
  const pict = (() => {
    const imgStyles = { backgroundColor: "#fff" }
    return <div className="pict" style={imgStyles}>
      <div className="moreContainer">
        <div className="and">∞</div>
        <div className="number">{more.toLocaleString()}</div>
        <div>more</div>
      </div>
    </div>
  })()

  const card = <div className="scaleArea">
    {pict}
  </div>

  return <div className="SongCard MoreSongCard">{card}</div>
}
