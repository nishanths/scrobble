import React from "react"
import { RouteComponentProps } from "react-router-dom";
import { UArgs, NProgress } from "../../shared/types"
import { Mode, DetailKind } from "./shared"
import { Scrobbles } from "./Scrobbles"
import { Detail } from "./Detail"
import { hexDecode } from "../../shared/util"
import { Color } from "../colorpicker"
import "../../scss/u/u.scss"

type UProps = UArgs & {
  wnd: Window
  mode: Mode
  color?: Color
  detail?: {
    kind: DetailKind
    hexIdent: string // hex-encoded song ident
  }
  nProgress: NProgress
} & RouteComponentProps;

// U is the root component for the username page, e.g.,
// https://scrobble.littleroot.org/u/whatever.
export const U: React.FC<UProps> = ({
  artworkBaseURL,
  profileUsername,
  logoutURL,
  self,
  private: priv,
  wnd,
  mode,
  color,
  detail,
  history,
  nProgress,
}) => {
  if (detail === undefined) {
    return <Scrobbles
      profileUsername={profileUsername}
      signedIn={!!logoutURL}
      artworkBaseURL={artworkBaseURL}
      private={priv}
      self={self}
      mode={mode}
      color={color}
      nProgress={nProgress}
      history={history}
      wnd={wnd}
    />
  } else {
    return <Detail
      profileUsername={profileUsername}
      artworkBaseURL={artworkBaseURL}
      private={priv}
      self={self}
      detailKind={detail.kind}
      songIdent={hexDecode(detail.hexIdent)}
      nProgress={nProgress}
      mode={mode}
      color={color}
      history={history}
    />
  }
}

