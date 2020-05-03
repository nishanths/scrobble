import React from "react";
import { SongCard } from "./SongCard"
import { Song, ArtworkHash } from "../shared/types"

interface ArtworksProps {
  hashes: ArtworkHash[];
  artworkBaseURL: string;
}

const artworkOnlySong = (artworkHash: string): Song => ({ artworkHash }) as unknown as Song

export const Artworks: React.SFC<ArtworksProps> = ({ hashes, artworkBaseURL }) => {
  const now = () => new Date()
  return <>
    {hashes.map(h => <SongCard key={h} song={artworkOnlySong(h)} artworkBaseURL={artworkBaseURL} now={now} />)}
  </>
}
