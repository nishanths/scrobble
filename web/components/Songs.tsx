import * as React from "react";
import { useSelector } from "react-redux"
import { SongCard } from "./SongCard"
import { Song } from "../shared/types"
import { State } from "../redux/reducers/username"

interface SongsProps {
  songs: Song[];
  artworkBaseURL: string;
  now: () => Date;
}

export const Songs: React.StatelessComponent<SongsProps> = ({ songs, artworkBaseURL, now }) => {
  useSelector((s: State) => {
    console.log(s)
  })

  return <>
    {songs.map(s => <SongCard key={s.ident} song={s} artworkBaseURL={artworkBaseURL} now={now} />)}
  </>
}
