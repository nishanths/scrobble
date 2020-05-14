import React from "react"
import { SetUsername } from "./SetUsername"
import { BootstrapArgs, Account } from "../../shared/types"
import { macOSAppLink, guideLink } from "../../shared/const"
import "../../scss/dashboard/dashboard.scss"

type DashboardProps = BootstrapArgs

export class Dashboard extends React.Component<DashboardProps, { account: Account }> {
  constructor(props: DashboardProps) {
    super(props)
    this.state = {
      account: this.props.account,
    }
  }

  render() {
    let content: React.ReactNode

    if (this.props.account.username) {
      // TODO
    } else {
      content = <div className="setUsername">
        <div className="instruction">
          <div className="heading">Set a username</div>
          <div className="desc">The username will be displayed on your profile, and will be present in hyperlinks to your profile.</div>
        </div>
        <SetUsername accountChange={(a) => { this.setState({ account: a }) }} />
      </div>
    }

    return <div className="Dashboard">

      <div className="start">
        <div className="heading">
          <span className="scrobble"><a href="/">scrobble</a></span>·
          <span className="desc">Apple Music scrobbling.</span></div>
      </div>

      <div className="nav">
        {this.state.account.username && <>
          <div className="item"><a href="TODO">API key</a></div> ·
          <div className="item"><a href="TODO">Profile privacy</a></div> ·
          <div className="item danger"><a href="TODO">Delete account</a></div> ·
        </>}
        <div className="item"><a href={this.props.logoutURL}>Sign out</a></div>
      </div>

      <div className="main">
        {content}
      </div>

      <div className="footer">
        <div className="item"><a href={guideLink}>Guide</a></div> ·
        <div className="item"><a href={macOSAppLink}>macOS app</a></div>
      </div>

    </div>
  }
}

