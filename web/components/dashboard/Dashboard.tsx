import React from "react"
import { SetUsername } from "./SetUsername"
import { SetPrivacy } from "./SetPrivacy"
import { NewAPIKey } from "./NewAPIKey"
import { Base as BaseComponent } from "./Base"
import { BootstrapArgs, Account, Notie } from "../../shared/types"
import { Mode } from "./shared"
import { macOSAppLink, guideLink } from "../../shared/const"
import { assertExhaustive } from "../../shared/util"
import { Link, RouteComponentProps, Redirect } from "react-router-dom"
import "../../scss/dashboard/dashboard.scss"

type History = RouteComponentProps["history"]

type DashboardProps = BootstrapArgs & {
  mode: Mode
  notie: Notie
  history: History
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
    // disallow other modes (besides base) if username isn't set
    if (this.props.mode !== Mode.Base && !this.state.account.username) {
      return <Redirect to="/" />
    }

    let content: React.ReactNode
    switch (this.props.mode) {
      case Mode.Base:
        content = <Base
          notie={this.props.notie}
          account={this.state.account}
          nSongs={this.props.totalSongs}
          lastScrobbleTime={this.props.lastScrobbleTime}
          accountChange={(a) => { this.setState({ account: a }) }}
        />
        break
      case Mode.Privacy:
        content = <Privacy notie={this.props.notie} privacy={this.state.account.private} privacyChange={(v) => {
          this.setState(s => {
            return {
              account: { ...s.account, private: v },
            }
          })
        }} />
        break
      case Mode.APIKey:
        content = <APIKey notie={this.props.notie} apiKey={this.state.account.apiKey} apiKeyChange={apiKey => {
          this.setState(s => {
            return {
              account: { ...s.account, apiKey },
            }
          })
        }} />
        break
      default:
        assertExhaustive(this.props.mode)
    }

    return <div className="Dashboard">
      <div className="start">
        <Link to="/">
          <div className="heading">
            <span className="scrobble">scrobble</span>·
            <span className="desc">Apple Music scrobbling.</span>
          </div>
        </Link>
      </div>

      <div className="nav">
        {this.state.account.username && <>
          <div className="item"><Link to="/dashboard/privacy">Profile privacy</Link> <span className="privacyHint">({this.state.account.private ? "private" : "public"})</span></div>
          <div className="item"><Link to="/dashboard/api-key">API key</Link></div>
          <div className="item danger"><a href="" onClick={this.onDeleteAccountClick.bind(this)}>Delete account…</a></div>
        </>}
        <div className="item"><a href={this.props.logoutURL} title={"Signed in as " + this.props.email}>Sign out</a></div>
      </div>

      <div className="main">
        {content}
      </div>

      <div className="footer">
        <div className="item"><a href={guideLink}>Guide</a></div>
        <div className="item"><a href={macOSAppLink}>macOS app</a></div>
        <div className="item"><a href="/">Home</a></div>
      </div>
    </div>
  }
}

const Base: React.SFC<{
  notie: Notie
  account: Account
  nSongs: number
  lastScrobbleTime: number
  accountChange: (a: Account) => void
}> = ({ notie, account, nSongs, lastScrobbleTime, accountChange }) => {
  if (account.username) {
    return <div className="base">
      <BaseComponent username={account.username} nSongs={nSongs} lastScrobbleTime={lastScrobbleTime} />
    </div>
  }

  return <div className="setUsername">
    <div className="instruction">
      <div className="heading">Set a username</div>
      <div className="desc">The username will be displayed on your profile, and will be present in hyperlinks to your profile.</div>
    </div>
    <SetUsername accountChange={(a) => { accountChange(a) }} notie={notie} />
  </div>
}

const Privacy: React.SFC<{ notie: Notie, privacy: boolean, privacyChange: (v: boolean) => void }> = ({ notie, privacy, privacyChange }) => {
  return <div className="privacy">
    <div className="instruction">
      <div className="heading">Profile Privacy</div>
      <div className="desc">
        <p>
          Your profile's privacy determines whether others can see the songs you scrobble on your profile.&nbsp;
          Your profile is currently <i>{privacy ? "private" : "public"}</i>
          {privacy ? " — only you may see your scrobbled songs when signed in." : " — others can see your scrobbled songs."}
        </p>
        <SetPrivacy privacy={privacy} privacyChange={privacyChange} notie={notie} />
      </div>
    </div>
  </div>
}

const APIKey: React.SFC<{ notie: Notie, apiKey: string, apiKeyChange: (s: string) => void }> = ({ notie, apiKey, apiKeyChange }) => {
  return <div className="apiKey">
    <div className="instruction">
      <div className="heading">API Key</div>
      <div className="desc">
        <p>
          Your API key is:
        </p>
        <pre><code>{apiKey}</code></pre>
        <p></p>
        <p>
          Enter the API key in the <a href={macOSAppLink} target="_blank" rel="noopener noreferrer">macOS application</a> to scrobble your Apple Music listening history.&nbsp;
        </p>
        <p>
          Keep it safe — anyone with the API key can scrobble songs on your behalf and view your scrobbled songs, even if your profile private.
        </p>
        <NewAPIKey apiKeyChange={apiKeyChange} notie={notie} />
      </div>
    </div>
  </div>
}
