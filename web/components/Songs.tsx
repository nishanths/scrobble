import React from "react";
import { SongCard, MoreSongCard } from "./songcard"
import { Song } from "../shared/types"

interface SongsProps {
  songs: Song[]
  artworkBaseURL: string
  albumCentric: boolean
  onSongClick: (s: Song) => void

  showMoreCard: boolean
  more?: number // required if showMoreCard is true

  showDates: boolean
  now?: () => Date // required if showDates is true
}

const debug = false;

export const Songs: React.SFC<SongsProps> = ({ albumCentric, songs, more, showMoreCard, artworkBaseURL, showDates, now, onSongClick }) => {
  if (showMoreCard && more == null) {
    throw "bad props"
  }
  return <>
    {debug && showMoreCard && <MoreSongCard more={more!} />}
    {songs.map(s => <SongCard
      key={s.ident}
      song={s}
      artworkBaseURL={artworkBaseURL}
      albumCentric={albumCentric}
      showDate={showDates}
      now={now}
      onSongClick={onSongClick}
    />)}
    {showMoreCard && <MoreSongCard more={more!} />}
  </>
}
