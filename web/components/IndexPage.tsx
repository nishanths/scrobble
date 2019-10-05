import * as React from "react";
import { BootstrapArgs, Account } from "../src/shared";
import { SetUsername, SetUsernameProps } from "./SetUsername";
import { AccountDetail, AccountDetailProps } from "./AccountDetail";
import "../scss/index.scss";

const deleteMessage = `Deleting your account will remove your account and the list of scrobbled songs. Artwork you may have uploaded might not be removed, and your username can be reused.

Delete account?`

type IndexPageProps = BootstrapArgs

export class IndexPage extends React.Component<IndexPageProps, {account: Account, deleteFail: boolean}> {
  private static readonly downloadURL = "https://github.com/nishanths/scrobble/releases/latest"

  constructor(props: IndexPageProps) {
    super(props)
    this.state = {
      account: this.props.account,
      deleteFail: false,
    }
  }

  updateAccount(a: Account) {
    this.setState({account: a})
  }

  private doDelete() {
    fetch(`/api/v1/account/delete`, {method: "POST"})
      .then(
        r => {
          if (r.status == 200) {
            window.location.assign(this.props.logoutURL)
            return
          }
          console.log("failed to delete: status=%d", r.status)
          this.setState({deleteFail: true})
        },
        err => {
          console.error(err)
          this.setState({deleteFail: true})
        }
      )
  }

  private onDeleteAccountClick(e: MouseEvent) {
    e.preventDefault()
    let ok = confirm(deleteMessage)
    if (!ok) {
      return
    }
    this.doDelete()
  }

  private signIn() {
    if (this.props.email && this.props.logoutURL) {
      return <p><a href={this.props.logoutURL}>Sign out</a> ({this.props.email})</p>
    }
    if (this.props.loginURL) {
      return <p><a href={this.props.loginURL}>Sign in with Google</a> to get started</p>
    }
    return null;
  }

  private visit() {
    return <p>Profiles can be found at <i>/u/username</i>, e.g., <a href="/u/nishanth">/u/nishanth</a>, <a href="/u/888">/u/888</a></p>
  }

  private profile() {
    return this.state.account && <p><a href={"/u/" + this.state.account.username}>Your scrobbles</a></p>
  }

  private download() {
    return <p><a href={IndexPage.downloadURL}>Download</a> menu bar scrobble client for iTunes (macOS 10.14+)</p>
  }

  render() {
    let errClass = (b: boolean) => b ? "error" : "error hidden"

    return <div>
      <h1>{this.props.host}</h1>
      {this.props.email && !this.state.account.username &&
        <SetUsername host={this.props.host} accountChange={this.updateAccount.bind(this)}/>}
      {this.props.email && this.state.account.username &&
        <AccountDetail account={this.state.account} host={this.props.host}/>}
      {this.signIn()}
      {this.state.account.username ? this.profile() : this.visit()}

      {/* TODO: a little gross that we're relying on logoutURL to indicate "existence of account" */}
      {this.props.logoutURL &&
        <p>
          <a href="" onClick={this.onDeleteAccountClick.bind(this)}>Delete account…</a>
          <span className={errClass(this.state.deleteFail)}>&nbsp;Failed to delete. Try again?</span>
        </p>
      }

      {this.download()}
    </div>
  }
}
