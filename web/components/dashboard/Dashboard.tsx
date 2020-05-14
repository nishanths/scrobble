import React from "react"
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
    return <div className="Dashboard">
      <div className="start">
        <div className="heading"><span className="scrobble">
          <a href="/">scrobble</a>
        </span> · <span className="desc">Apple Music scrobbling</span></div>
      </div>
      <div className="nav">
        {this.state.account.username && <>
          <div className="item"><a>API key</a></div> ·
          <div className="item"><a>Profile privacy</a></div> ·
          <div className="item danger"><a>Delete account</a></div> ·
        </>}
        <div className="item"><a href={this.props.logoutURL}>Sign out</a></div>
      </div>

      <div className="footer">
        <div className="item"><a href={guideLink}>Guide</a></div> ·
        <div className="item"><a href={macOSAppLink}>macOS app</a></div>
      </div>
    </div>
  }
}

