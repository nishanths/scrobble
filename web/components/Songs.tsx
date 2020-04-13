import React from "react";
import { SongCard } from "./SongCard"
import { MoreSongCard } from "./MoreSongCard"
import { Song } from "../shared/types"

interface SongsProps {
  songs: Song[];
  more: number;
  showMore: boolean
  artworkBaseURL: string;
  now: () => Date;
}

const debug = false;

export const Songs: React.StatelessComponent<SongsProps> = ({ songs, more, showMore, artworkBaseURL, now }) => {
  return <>
    {debug && <MoreSongCard more={more} />}
    {songs.map(s => <SongCard key={s.ident} song={s} artworkBaseURL={artworkBaseURL} now={now} />)}
    {showMore && <MoreSongCard more={more} />}
  </>
}
