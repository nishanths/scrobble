import React from "react";
import { BootstrapArgs, Account, Notie } from "../shared/types";
import { SetUsername } from "./SetUsername";
import { AccountDetail } from "./AccountDetail";
import "../scss/root.scss";

const deleteMessage = `Deleting your account will remove your account and the list of scrobbled songs. Artwork you may have uploaded might not be removed, and your username can be reused.

Delete account?`

type RootProps = BootstrapArgs

declare const notie: Notie

export class Root extends React.Component<RootProps, { account: Account }> {
  private static readonly downloadURL = "https://github.com/nishanths/scrobble/releases/latest"

  constructor(props: RootProps) {
    super(props)
    this.state = {
      account: this.props.account,
    }
  }

  private updateAccount(a: Account) {
    this.setState({ account: a })
  }

  private doDelete() {
    fetch(`/api/v1/account/delete`, { method: "POST" })
      .then(
        r => {
          if (r.status == 200) {
            window.location.assign(this.props.logoutURL)
            return
          }
          console.log("failed to delete: status=%d", r.status)
          notie.alert({ type: "error", text: "Failed to delete account.", stay: true })
        },
        err => {
          console.error(err)
          notie.alert({ type: "error", text: "Failed to delete account.", stay: true })
        }
      )
  }

  private onDeleteAccountClick(e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) {
    e.preventDefault()
    const ok = window.confirm(deleteMessage)
    if (!ok) {
      return
    }
    this.doDelete()
  }

  private getStarted() {
    return <>
      <p>To get started with your own profile, <a href={this.props.loginURL}>sign in with Google</a>.</p>
      <p>See an <a href="/u/nishanth">example user profile</a>.</p>
    </>
  }

  private profile() {
    return <p><a href={"/u/" + this.state.account.username}>See your scrobbles.</a></p>
  }

  private download() {
    return <p><a href={Root.downloadURL}>Menu bar application</a> for macOS to scrobble Apple Music history.</p>
  }

  private privacyPolicy() {
    return <p><a href="/privacy-policy">Privacy.</a></p>
  }

  render() {
    return <div className="Root">
      <div className="heading">
        <h1><strong>scrobble</strong></h1>
        <p className="larger">Music scrobbling for Apple Music.</p>
      </div>

      <div className="content larger">
        {this.state.account.username && <section>
          {this.profile()}
        </section>}

        {this.props.loginURL && <section>
          <h2>Get Started</h2>
          {this.getStarted()}
        </section>}

        {this.props.email && <section>
          <h2>Account</h2>
          <p><a href={this.props.logoutURL}>Sign out.</a> (You are signed in using Google as <span className="meta">{this.props.email}</span>.)</p>
          {this.state.account.username && <>
            <AccountDetail account={this.state.account} />
            <p>Your username is <a href={"/u/" + this.state.account.username}>{this.state.account.username}</a>.</p>
            <p><a href="" className="danger" onClick={this.onDeleteAccountClick.bind(this)}>Delete account…</a></p>
          </>}
          {!this.state.account.username && <>
            <p>To create your profile set a username. It can contain numbers and must be lowercase otherwise.</p>
            <SetUsername accountChange={this.updateAccount.bind(this)} />
          </>}
        </section>}

        <section>
          <h2>Helpful Links</h2>
          {<p><a href="/guide">How do I use this?</a></p>}
          {this.download()}
          {this.privacyPolicy()}
        </section>
      </div>
    </div>
  }
}
