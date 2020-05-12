import React, { useEffect } from "react"
import { Song } from "../../shared/types"
import { dateDisplay } from "../../shared/time"
import { pluralize } from "../../shared/util"

// TrackLink is the track link area of a Picture.
const TrackLink: React.SFC<{ previewURL: string }> = ({ previewURL }) => {
  let trackLinkAreaElem: HTMLDivElement | null = null

  useEffect(() => {
    // Stack Overflow says onclick enables :hover on iOS
    // TODO: is this really necessary?
    if (trackLinkAreaElem !== null) { trackLinkAreaElem.setAttribute("onclick", "") }
  }, [])

  return <a href={previewURL} title={previewURL} target="_blank" rel="noopener noreferrer">
    <div className="trackLinkArea" ref={r => { trackLinkAreaElem = r }}>
      <div className="trackLink"></div>
    </div>
  </a>
}

// Picture is the picture area for a SongCard.
export const Picture: React.SFC<{
  song: Song
  artworkBaseURL: string
  albumCentric: boolean
}> = ({ song, artworkBaseURL, albumCentric }) => {
  const previewURL = trackViewURL(song.trackViewURL, albumCentric)
  const artworkURL = song.artworkHash ? artworkBaseURL + "/" + song.artworkHash : "";
  const imgStyles = artworkURL ? { backgroundImage: `url(${artworkURL})` } : { backgroundColor: "#fff" }

  return <div className="pict" style={imgStyles}>
    {previewURL && <TrackLink previewURL={previewURL} />}
  </div>
}

// MorePicture is the picture area used for a MoreSongCard.
export const MorePicture: React.SFC<{ more: number }> = ({ more }) => {
  const imgStyles = { backgroundColor: "#fff" }

  return <div className="pict" style={imgStyles}>
    <div className="moreContainer">
      <div className="and">∞</div>
      <div className="number">{more.toLocaleString()}</div>
      <div>more</div>
    </div>
  </div>
}

// LargePicture is the picture area for a LargeSongCard.
export const LargePicture: React.SFC<{ song: Song, artworkBaseURL: string }> = ({
  song,
  artworkBaseURL,
}) => {
  const artworkURL = song.artworkHash ? artworkBaseURL + "/" + song.artworkHash : "";
  const imgStyles = artworkURL ? { backgroundImage: `url(${artworkURL})` } : { backgroundColor: "#fff" }

  return <div className="pict" style={imgStyles}>
  </div>
}

// Meta is the metadata area for a SongCard.
export const Meta: React.SFC<{
  song: Song
  albumCentric: boolean
  showDate: boolean
  now?: () => Date
  onClick: () => void
}> = ({
  song: s,
  albumCentric,
  showDate,
  now,
  onClick,
}) => {
    const tooltip = metaTooltip(albumCentric, s.title, s.artistName, s.albumTitle)
    const title = albumCentric ? s.albumTitle : s.title
    const includeLoved = !albumCentric && s.loved
    const includeDate = !albumCentric && showDate && s.lastPlayed

    return <div className="meta" title={tooltip} onClick={onClick}>
      <div className="title">
        <span className="titleContent">{title}</span>
        {includeLoved && <span className="love"></span>}
      </div>

      <div className="other">
        {s.artistName && <span className="artist">{s.artistName}</span>}
      </div>

      {includeDate && <time className="date">{dateDisplay(new Date(s.lastPlayed * 1000), now!())}</time>}
    </div>
  }

// LargeMeta is the metadata section for a LargeSongCard.
export const LargeMeta: React.SFC<{
  song: Song
  albumCentric: boolean
}> = ({
  song: s,
  albumCentric,
}) => {
    const includeLoved = !albumCentric && s.loved
    const includePlayCount = !albumCentric

    const previewURL = trackViewURL(s.trackViewURL, albumCentric)
    const title = albumCentric ? s.albumTitle : s.title
    let other = albumCentric ? `${s.artistName}` : `${s.artistName} – ${s.albumTitle}`
    if (s.year) {
      other += " – " + s.year
    }

    return <a className={"link"} href={previewURL} title={previewURL} target="_blank" rel="noopener noreferrer">
      <div className="meta">
        <div className="title">
          <span className="titleContent">{title}</span>
          {includeLoved && <span className="love"></span>}
        </div>
        <div className="other">
          <span className="otherContent">{other}</span>
        </div>
        {includePlayCount && <div className="lastLine">
          <span className="playCount">Played {s.playCount.toLocaleString()} {pluralize("time", s.playCount)}</span>
        </div>}
      </div>
    </a>
  }

const trackViewURL = (songTrackViewURL: string, albumCentric: boolean): string => {
  if (songTrackViewURL === "" || albumCentric === false) {
    return songTrackViewURL
  }
  // clear song portion (a.k.a. the query string), so that link goes to album.
  // e.g. https://music.apple.com/us/album/crystalised/329481191?i=329481203&uo=4
  try {
    const u = new URL(songTrackViewURL)
    u.search = ""
    return u.toString()
  } catch (e) {
    return songTrackViewURL
  }
}

const metaTooltip = (albumCentric: boolean, title: string, artist: string, album: string): string => {
  if (albumCentric) {
    let tooltip = ""
    if (artist || album) {
      if (artist) { tooltip += artist }
      if (artist && album) { tooltip += " — " }
      if (album) { tooltip += album }
    }
    return tooltip
  }

  let tooltip = title
  if (artist || album) {
    tooltip += "\n"
    if (artist) { tooltip += artist }
    if (artist && album) { tooltip += " — " }
    if (album) { tooltip += album }
  }
  return tooltip
}
