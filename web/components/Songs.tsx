import * as React from "react";
import { SongCard } from "./SongCard"
import { Song } from "../shared/types"

interface SongsProps {
  songs: Song[];
  artworkBaseURL: string;
  now: () => Date;
}

export const Songs: React.StatelessComponent<SongsProps> = ({ songs, artworkBaseURL, now }) => {
  return <>
    {songs.map(s => <SongCard key={s.ident} song={s} artworkBaseURL={artworkBaseURL} now={now} />)}
  </>
}
