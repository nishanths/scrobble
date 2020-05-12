import React, { useEffect, useRef } from "react";
import { Song } from "../../shared/types"
import { dateDisplay } from "../../shared/time"

interface SongCardProps {
  song: Song; // rendering degrades gracefully if properties are missing
  artworkBaseURL: string;
  // album-centric instead of song-centric
  // e.g. use album instead of song title in the title areas
  albumCentric: boolean;

  showDate: boolean
  now?: () => Date // required if showDates is true
}

export const SongCard: React.StatelessComponent<SongCardProps> = ({
  song,
  artworkBaseURL,
  albumCentric,
  showDate,
  now
}) => {
  let trackLinkAreaElem: HTMLDivElement | null = null

  useEffect(() => {
    if (trackLinkAreaElem != null) { trackLinkAreaElem.setAttribute("onclick", "") } // Stack Overflow says onclick enables :hover on iOS
  }, [])

  const artworkURL = song.artworkHash ? artworkBaseURL + "/" + song.artworkHash : "";

  const trackViewURL = (() => {
    if (albumCentric && song.trackViewURL != "") {
      // clear song portion (aka query string), so that link goes to album
      // e.g. https://music.apple.com/us/album/crystalised/329481191?i=329481203&uo=4
      try {
        const u = new URL(song.trackViewURL)
        u.search = ""
        return u.toString()
      } catch (e) {
        return song.trackViewURL
      }
    }
    return song.trackViewURL
  })()

  const trackLinkArea = <a href={trackViewURL} title={trackViewURL} target="_blank">
    <div className="trackLinkArea" ref={r => { trackLinkAreaElem = r }}>
      <div className="trackLink"></div>
    </div>
  </a>

  const tooltip = (() => {
    const s = song

    if (albumCentric) {
      let tooltip = ""
      if (s.artistName || s.albumTitle) {
        if (s.artistName) { tooltip += s.artistName }
        if (s.artistName && s.albumTitle) { tooltip += " — " }
        if (s.albumTitle) { tooltip += s.albumTitle }
      }
      return tooltip
    }

    let tooltip = s.title
    if (s.artistName || s.albumTitle) {
      tooltip += "\n"
      if (s.artistName) { tooltip += s.artistName }
      if (s.artistName && s.albumTitle) { tooltip += " — " }
      if (s.albumTitle) { tooltip += s.albumTitle }
    }
    return tooltip
  })()

  const pict = (() => {
    const imgStyles = artworkURL ? { backgroundImage: `url(${artworkURL})` } : { backgroundColor: "#fff" }
    return <div className="pict" style={imgStyles}>{trackViewURL && trackLinkArea}</div>
  })()

  const meta = (() => {
    const s = song
    return <div className="meta" title={tooltip}>
      <div className="title">
        <span className="titleContent">{albumCentric ? s.albumTitle : s.title}</span>
        {(s.loved && !albumCentric) && <span className="love"></span>}
      </div>
      <div className="other">
        {s.artistName && <span className="artist">{s.artistName}</span>}
      </div>
      {showDate && s.lastPlayed && <time className="date">{dateDisplay(new Date(s.lastPlayed * 1000), now!())}</time>}
    </div>
  })()

  const card = <div className="scaleArea">
    {pict}
    {meta}
  </div>

  return <div className="SongCard">{card}</div>
}