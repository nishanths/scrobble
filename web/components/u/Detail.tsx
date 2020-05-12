import React from "react"
import { SongState } from "../../redux/types/song"
import { NProgress } from "../../shared/types"
import { assert, assertExhaustive } from "../../shared/util"
import { DetailKind } from "./shared"
import { CloseIcon } from "../CloseIcon"
import { Mode, pathForMode, pathForColor } from "./shared"
import { Color } from "../colorpicker"
import { RouteComponentProps } from "react-router-dom";
import "../../scss/u/detail.scss"

import 'react-responsive-modal/styles.css';
import { Modal } from 'react-responsive-modal';

type History = RouteComponentProps["history"]

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

export const Detail: React.StatelessComponent<{
  song: SongState
  profileUsername: string
  mode: Mode
  color: Color | undefined
  private: boolean
  self: boolean
  detailKind: DetailKind
  nProgress: NProgress
  history: History
}> = ({
  song,
  profileUsername,
  mode,
  color,
  private: priv,
  self,
  detailKind,
  nProgress,
  history,
}) => {
    const modal = (content: React.ReactNode) => <Modal
      open={true}
      onClose={() => { history.push("/u/" + profileUsername + pathForMode(mode) + pathForColor(color)) }}
      center
      classNames={{ modal: "detailModal", overlay: "detailOverlay", closeButton: "detailCloseButton" }}
      closeOnEsc={true}
      animationDuration={500}
      closeIcon={CloseIcon}>
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

    return <>{modal(null)}</> // TODO
  }
