import React from "react"
import { assertExhaustive, assert } from "../../shared/util"
import { NProgress } from "../../shared/types"
import { Mode } from "./shared"
import { Color } from "../colorpicker"
import { ScrobblesState } from "../../redux/types/scrobbles"
import { Songs } from "../Songs"

export const Scrobbles: React.StatelessComponent<{
  scrobbles: ScrobblesState | null
  artworkBaseURL: string
  endIdx: number
  private: boolean
  self: boolean
  mode: Mode
  color: Color | undefined
  header: JSX.Element
  top: JSX.Element
  nProgress: NProgress
}> = ({
  scrobbles,
  artworkBaseURL,
  endIdx,
  private: priv,
  self,
  mode,
  color,
  header,
  top,
  nProgress
}) => {
    // Easy case. For private accounts that aren't the current user, render the
    // private info-message.
    if (priv === true && self === false) {
      return <>
        {header}
        <div className="info">(This user's scrobbles are private.)</div>
      </>
    }

    // If in the Color mode and no color is selected, render the top area and
    // the color picker, and we're done.
    if (mode === Mode.Color && color === undefined) {
      return <>
        {top}
      </>
    }

    assert(scrobbles !== null, "scrobbles unexpectedly null")

    if (scrobbles.fetching === true) {
      nProgress.start()
      return <>{top}</>
    }

    if (scrobbles.error === true) {
      nProgress.done()
      return <>
        {header}
        <div className="info">(Failed to fetch scrobbles.)</div>
      </>
    }

    // handle initial redux state
    if (scrobbles.done === false) {
      return null
    }

    nProgress.done()

    // can happen if the privacy was changed after the initial server page load
    if (scrobbles.private) {
      return <>
        {header}
        <div className="info">(This user's scrobbles are private.)</div>
      </>
    }

    if (scrobbles.items.length === 0) {
      return <>
        {top}
        <div className="info">({self ? "You haven't" : "This user hasn't"} scrobbled {mode != Mode.All ? "matching " : ""}songs yet.)</div>
      </>
    }

    const itemsToShow = scrobbles.items.slice(0, endIdx);

    switch (mode) {
      case Mode.All:
      case Mode.Loved: {
        return <>
          {top}
          <div className="songs">
            <Songs
              songs={itemsToShow}
              artworkBaseURL={artworkBaseURL}
              albumCentric={false}
              more={scrobbles.total! - itemsToShow.length}
              // "showing all songs that are available on the client" && "more number of songs present for the user "
              showMoreCard={(itemsToShow.length === scrobbles.items.length) && (scrobbles.total! > scrobbles.items.length)}
              showDates={true}
              now={() => new Date()}
            />
          </div>
        </>
      }

      case Mode.Color: {
        return <>
          {top}
          <div className="songs">
            <Songs
              songs={itemsToShow}
              artworkBaseURL={artworkBaseURL}
              albumCentric={true}
              showMoreCard={false}
              showDates={false}
            />
          </div>
        </>
      }

      default: {
        assertExhaustive(mode)
      }
    }
  }
