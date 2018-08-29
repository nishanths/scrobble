import * as React from "react";
import { Account } from "../src/shared"
import "../scss/account"

export interface AccountDetailProps {
  account: Account
  host: string
}

interface AccountDetailState {
  apiKey: string
  keyGenerateErr: string
  private: boolean
  showPrivacySaved: boolean // "saved" text for privacy input
  privacySaveErr: string
}

export class AccountDetail extends React.Component<AccountDetailProps, AccountDetailState> {
  private privacyCheckbox: HTMLInputElement | null = null;
  private generating = false // inflight request to generate API key?
  private savedTimeout: number|undefined = undefined;

  constructor(props: AccountDetailProps) {
    super(props)
    this.state = {
      apiKey: this.props.account.apiKey,
      keyGenerateErr: "",
      private: this.props.account.private,
      showPrivacySaved: false,
      privacySaveErr: ""
    }
  }

  private newAPIKey() {
    this.setState({keyGenerateErr: ""})
    this.generating = true
    let success = false

    const genericErr = "Failed to generate API key. Try again?"

    fetch(`https://${this.props.host}/newAPIKey`, {method: "POST"})
      .then(res => {
        if (res.status == 200) {
          success = true
          return res.json()
        }
        return res.blob()
      })
      .then(
        r => {
          if (!success) { return }
          this.setState({apiKey: r as string})
        },
        err => {
          console.error(err)
          this.setState({keyGenerateErr: genericErr})
        }
      ).then(() => { // TODO: finally() would be nicer
        this.generating = false
      })
  }

  private setPrivacy(v: boolean) {
    this.setState({showPrivacySaved: false, privacySaveErr: ""})
    let success = false

    fetch(`https://${this.props.host}/setPrivacy?privacy=${v.toString()}`, {method: "POST"})
      .then(
        r => {
          if (r.status == 200) {
            window.clearTimeout(this.savedTimeout)
            this.setState({showPrivacySaved: true})
            this.savedTimeout = window.setTimeout(() => {
              this.setState({showPrivacySaved: false})
            }, 1500)
          }
        },
        err => {
          console.error(err)
          this.setState({
            private: !v, // revert
            privacySaveErr: "Failed to toggle privacy. Try again?",
          })
        }
      )
  }

  private onRegenerateClick(e: MouseEvent) {
    e.preventDefault()
    if (this.generating) { return }
    this.newAPIKey()
  }

  private onPrivacyClick(e: Event) {
    let v = this.privacyCheckbox!.checked
    this.setState({private: v})
    this.setPrivacy(v)
  }

  render() {
    let errClass = this.state.keyGenerateErr ? "error" : "error hidden"
    let savedClass = this.state.showPrivacySaved ? "success" : "success hidden"
    let failedSaveClass = this.state.privacySaveErr ? "error" : "error hidden"

    return <div className="account">
      <table>
        <tbody>
          <tr><td>Username</td><td>{this.props.account.username}</td></tr>
          <tr>
            <td>API Key</td>
            <td className="mono">{this.state.apiKey}</td>
            <td><a href="" onClick={this.onRegenerateClick.bind(this)}>Regenerate</a></td>
            <td><span className={errClass}>{this.state.keyGenerateErr}</span></td>
          </tr>
          <tr>
            <td>Private</td>
            <td>
              <input type="checkbox" checked={this.state.private} onClick={this.onPrivacyClick.bind(this)} ref={r => {this.privacyCheckbox = r}}/>
            </td>
            <td className={savedClass}>Saved</td>
            <td></td>
            <td className={failedSaveClass}>{this.state.privacySaveErr}</td>
          </tr>
        </tbody>
      </table>
    </div>
  }
}
