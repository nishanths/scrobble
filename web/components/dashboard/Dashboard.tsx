import React from "react"
import { SetUsername } from "./SetUsername"
import { BootstrapArgs, Account, Notie } from "../../shared/types"
import { Mode } from "./shared"
import { macOSAppLink, guideLink } from "../../shared/const"
import { assertExhaustive } from "../../shared/util"
import { Link } from "react-router-dom"
import "../../scss/dashboard/dashboard.scss"

type DashboardProps = BootstrapArgs & {
  mode: Mode
  notie: Notie
}

const deletePrompt = `Deleting your account will remove your account and the list of scrobbled songs. \
Artwork you may have uploaded might not be removed, and your username can be reused. \


Delete account?`

const failedDelete = `Failed to delete account. Try again.`

export class Dashboard extends React.Component<DashboardProps, { account: Account }> {
  constructor(props: DashboardProps) {
    super(props)
    this.state = {
      account: this.props.account,
    }
  }

  private doDelete() {
    fetch(`/api/v1/account/delete`, { method: "POST" })
      .then(
        r => {
          if (r.status === 200) {
            window.location.assign(this.props.logoutURL)
            return
          }
          console.log("failed to delete: status=%d", r.status)
          this.props.notie.alert({ type: "error", text: failedDelete, stay: true })
        },
        err => {
          console.error(err)
          this.props.notie.alert({ type: "error", text: failedDelete, stay: true })
        }
      )
  }

  private onDeleteAccountClick(e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) {
    e.preventDefault()
    const ok = window.confirm(deletePrompt)
    if (!ok) {
      return
    }
    this.doDelete()
  }

  render() {
    let content: React.ReactNode
    switch (this.props.mode) {
      case Mode.Base:
        content = <Base account={this.props.account} accountChange={(a) => this.setState({ account: a})} />
        break
      case Mode.Privacy:
        content = null // TODO
        break
      case Mode.APIKey:
        content = null // TODO
        break
      default:
        assertExhaustive(this.props.mode)
    }

    return <div className="Dashboard">
      <div className="start">
        <div className="heading">
          <span className="scrobble"><Link to="/">scrobble</Link></span>·
          <span className="desc">Apple Music scrobbling.</span></div>
      </div>

      <div className="nav">
        {this.state.account.username && <>
          <div className="item"><Link to="/dashboard/api-key">API key</Link></div> ·
          <div className="item"><Link to="/dashboard/privacy">Profile privacy</Link></div> ·
          <div className="item danger"><a href="" onClick={this.onDeleteAccountClick.bind(this)}>Delete account…</a></div>·
        </>}
        <div className="item"><a href={this.props.logoutURL} title={"Signed in as " + this.props.email}>Sign out</a></div>
      </div>

      <div className="main">
        {content}
      </div>

      <div className="footer">
        <div className="item"><a href={guideLink}>Guide</a></div>·
        <div className="item"><a href={macOSAppLink}>macOS app</a></div>
      </div>
    </div>
  }
}

const Base: React.SFC<{ account: Account, accountChange: (a: Account) => void }> = ({ account, accountChange }) => {
  if (account.username) {
    // TODO
    return <div></div>
  }

  return <div className="setUsername">
    <div className="instruction">
      <div className="heading">Set a username</div>
      <div className="desc">The username will be displayed on your profile, and will be present in hyperlinks to your profile.</div>
    </div>
    <SetUsername accountChange={(a) => { accountChange(a) }} />
  </div>
}
