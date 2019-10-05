import * as React from "react";
import { useEffect, useRef } from "react";
import { Song } from "../src/shared"
import { displayString as dateDisplayString } from "../src/time"

interface SongCardProps {
  song: Song;
  artworkBaseURL: string;
  now: () => Date;
}

export const SongCard: React.StatelessComponent<SongCardProps> = ({
  song,
  artworkBaseURL,
  now
}) => {
  const initialMount = useRef(true);
  let trackLinkAreaElem: HTMLDivElement|null = null

  useEffect(() => {
    if (initialMount.current === false) { return }
    initialMount.current = false;
    trackLinkAreaElem!.setAttribute("onclick", "") // Stack Overflow says onclick enables :hover on iOS
  })

  const artworkURL = song.artworkHash ? artworkBaseURL + "/" + song.artworkHash : "";

  const trackLinkArea = <a href={song.trackViewURL} title={song.trackViewURL} target="_blank">
    <div className="trackLinkArea" ref={r => { trackLinkAreaElem = r }}>
      <div className="trackLink"></div>
    </div>
  </a>

  const tooltip = (() => {
    const s = song
    let tooltip = s.title
    if (s.artistName || s.albumTitle) {
      tooltip += "\n"
      if (s.artistName) { tooltip += s.artistName }
      if (s.artistName && s.artistName) { tooltip += " — " }
      if (s.albumTitle) { tooltip += s.albumTitle }
    }
    return tooltip
  })()

  const pict = (() => {
    const imgStyles = artworkURL ? {backgroundImage: `url(${artworkURL})`} : {backgroundColor: "#fff"}
    return <div className="pict" style={imgStyles}>{song.trackViewURL && trackLinkArea}</div>
  })()

  const meta = (() => {
    const s = song
    return <div className="meta" title={tooltip}>
      <div className="title">
        <span className="titleContent">{s.title}</span>
        {s.loved && <span className="love"></span>}
      </div>
      <div className="other">
        {s.artistName && <span className="artist">{s.artistName}</span>}
      </div>
      {s.lastPlayed && <time className="date">{dateDisplayString(new Date(s.lastPlayed * 1000), now())}</time>}
    </div>
  })()

  const card = <div className="scaleArea">
    {pict}
    {meta}
  </div>

  return <div className="SongCard">{card}</div>
}
