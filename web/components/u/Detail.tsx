import React from "react"
import { SongState } from "../../redux/types/song"
import { NProgress } from "../../shared/types"
import { DetailKind } from "./shared"
import { CloseIcon } from "../CloseIcon"

import "../../scss/detail-modal.scss"
import 'react-responsive-modal/styles.css';
import { Modal } from 'react-responsive-modal';

export const Detail: React.StatelessComponent<{
  song: SongState
  detailKind: DetailKind
  nProgress: NProgress
}> = ({
  song,
  detailKind,
  nProgress,
}) => {

    return null
  }

  // if (detail !== undefined) {
  //   assert(detailSong !== null, "detailSong unexpectedly null")

  //   if (detailSong.fetching) {
  //     NProgress.start()
  //   }
  //   // handle initial redux state
  //   if (detailSong.done === false) {
  //     return null
  //   }
  //   NProgress.done()

  //   const modalContent = <div className="flexContainer">
  //     {detailSong.item!.ident}
  //   </div>

  //   console.log(hexEncode(detailSong.item!.ident))
  //   console.log(hexDecode(detail!.hexIdent))

  //   const modal = <Modal
  //     open={true}
  //     onClose={() => { history.push("/u/" + profileUsername + pathForMode(mode) + pathForColor(color)) }}
  //     center
  //     classNames={{ modal: "detailModal", overlay: "detailOverlay", closeButton: "detailCloseButton" }}
  //     closeOnEsc={true}
  //     animationDuration={500}
  //     closeIcon={CloseIcon}>
  //     {modalContent}
  //   </Modal>

  //   return <>{modal}</>
  // }
