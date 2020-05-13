import React, { useRef, useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import { RouteComponentProps } from "react-router-dom";
import { NProgress } from "../../shared/types"
import { assert, assertExhaustive } from "../../shared/util"
import { DetailKind, pathForMode, pathForColor, Mode } from "./shared"
import { Header } from "./top"
import { State } from "../../redux/types/u"
import { CloseIcon } from "../CloseIcon"
import { Color } from "../colorpicker"
import { LargeSongCard } from "../songcard"
import { fetchSong } from "../../redux/actions/song"
import { Modal } from 'react-responsive-modal'
import 'react-responsive-modal/styles.css'
import "../../scss/u/detail.scss"

const nounForDetailKind = (k: DetailKind): string => {
  switch (k) {
    case DetailKind.Song:
      return "song"
    case DetailKind.Album:
      return "album"
    default:
      assertExhaustive(k)
  }
}

type History = RouteComponentProps["history"]

export const Detail: React.StatelessComponent<{
  profileUsername: string
  artworkBaseURL: string
  private: boolean
  self: boolean
  detailKind: DetailKind
  songIdent: string
  nProgress: NProgress
  mode: Mode
  color: Color | undefined
  history: History
}> = ({
  profileUsername,
  artworkBaseURL,
  private: priv,
  self,
  detailKind,
  songIdent,
  nProgress,
  mode,
  color,
  history,
}) => {
    const dispatch = useDispatch()

    // redux state
    const song = useSelector((s: State) => {
      const key = songIdent
      return s.songs.getOrDefault(key)
    })
    const songRef = useRef(song)
    useEffect(() => { songRef.current = song }, [song])

    // fetch song detail
    useEffect(() => {
      const song = songRef.current
      if (song === null || (song.done === false && song.fetching === false) || song.error === true) {
        dispatch(fetchSong(profileUsername, songIdent))
      }
    }, [profileUsername, songIdent])

    // ... render ...

    const header = Header(profileUsername, false, false)

    const modal = (content: React.ReactNode) => <Modal
      open={true}
      onClose={() => { nProgress.done(); history.push("/u/" + profileUsername + pathForMode(mode) + pathForColor(color)) }}
      center
      classNames={{ modal: "detailModal", overlay: "detailOverlay", closeButton: "detailCloseButton" }}
      closeOnOverlayClick={false}
      closeOnEsc={true}
      animationDuration={500}
      closeIcon={CloseIcon}>
      {header}
      <div className="flexContainer">
        {content}
      </div>
    </Modal>

    const noun = nounForDetailKind(detailKind)
    const privateContent = <div className="info">(This user's songs are private.)</div>

    if (priv === true && self === false) {
      return modal(privateContent)
    }

    if (song.fetching === true) {
      nProgress.start()
      return modal(null)
    }
    if (song.error === true) {
      nProgress.done()
      return modal(<div className="info">(Failed to fetch scrobbles.)</div>)
    }
    // handle initial redux state
    if (song.done === false) {
      return null
    }
    nProgress.done()

    if (song.private === true) {
      return modal(privateContent)
    }
    if (song.noSuchSong === true) {
      return modal(<div className="info">(Failed to find the specified {noun}.)</div>)
    }

    const item = song.item
    assert(item !== null, "item should not be null")

    return modal(<>
      <LargeSongCard
        song={item}
        artworkBaseURL={artworkBaseURL}
        albumCentric={detailKind === DetailKind.Album}
        now={() => new Date()}
      />
    </>)
  }
