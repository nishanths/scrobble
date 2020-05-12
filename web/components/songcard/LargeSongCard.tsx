import React from "react"
import { Song } from "../../shared/types"
import { SongCardProps, SongCard } from "./SongCard"

interface LargeSongCardProps extends SongCardProps { }

export const LargeSongCard: React.SFC<LargeSongCardProps> = (props) => {
  return <SongCard {...props} />
}
