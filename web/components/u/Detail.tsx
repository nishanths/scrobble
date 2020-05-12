import React from "react"
import { SongState } from "../../redux/types/song"
import { NProgress } from "../../shared/types"
import { assert } from "../../shared/util"
import { DetailKind } from "./shared"
import { CloseIcon } from "../CloseIcon"
import { Mode, pathForMode, pathForColor } from "./shared"
import { Color } from "../colorpicker"
import { RouteComponentProps } from "react-router-dom";

import "../../scss/detail-modal.scss"
import 'react-responsive-modal/styles.css';
import { Modal } from 'react-responsive-modal';

type History = RouteComponentProps["history"]

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
    if (priv === true && self === false) {
      return null // TODO
    }

    if (song.fetching === true) {
      nProgress.start()
      return null // TODO
    }
    if (song.error === true) {
      nProgress.done()
      return null // TODO
    }
    // handle initial redux state
    if (song.done === false) {
      return null
    }
    nProgress.done()

    if (song.private === true) {
      return null // TODO
    }
    if (song.noSuchSong === true) {
      return null // TODO
    }

    const item = song.item
    assert(item !== null, "item should not be null")

    const content = <div className="flexContainer">
      {item.ident}
    </div>

    const modal = <Modal
      open={true}
      onClose={() => { history.push("/u/" + profileUsername + pathForMode(mode) + pathForColor(color)) }}
      center
      classNames={{ modal: "detailModal", overlay: "detailOverlay", closeButton: "detailCloseButton" }}
      closeOnEsc={true}
      animationDuration={500}
      closeIcon={CloseIcon}>
      {content}
    </Modal>

    return <>{modal}</>
  }
