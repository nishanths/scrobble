import React from "react";
import { Account } from "../shared/types"
import { cookieAuthErrorMessage } from "../shared/util"
import "../scss/account"

// TODO: refactor using hooks

interface AccountDetailProps {
  account: Account
  accountChange: (a: Account) => void
}

interface AccountDetailState {
  apiKey: string
  keyGenerateErr: string
  private: boolean
  showPrivacySaved: boolean
  setPrivacyErr: string
}

export class AccountDetail extends React.Component<AccountDetailProps, AccountDetailState> {
  private privacyCheckbox: HTMLInputElement | null = null;
  private generating = false // inflight request to generate API key?
  private savedTimeout: number | undefined = undefined;

  constructor(props: AccountDetailProps) {
    super(props)
    this.state = {
      apiKey: this.props.account.apiKey,
      keyGenerateErr: "",
      private: this.props.account.private,
      showPrivacySaved: false,
      setPrivacyErr: "",
    }
  }

  private newAPIKey() {
    this.setState({ keyGenerateErr: "" })
    this.generating = true
    let success = false

    const genericErr = "Failed to generate API Key. Try again?"

    fetch(`/newAPIKey`, { method: "POST" })
      .then(res => {
        if (res.status == 200) {
          success = true
          return res.json()
        }
        if (res.status == 401) {
          this.setState({ keyGenerateErr: cookieAuthErrorMessage })
          return res.blob()
        }
        console.log("failed to generate API key: status=%d", res.status)
        this.setState({ keyGenerateErr: genericErr })
        return res.blob()
      })
      .then(
        r => {
          if (!success) { return }
          const apiKey = r as string
          this.props.accountChange({ ...this.props.account, apiKey })
          this.setState({ apiKey })
        },
        err => {
          console.error(err)
          this.setState({ keyGenerateErr: genericErr })
        }
      ).then(() => { // TODO: finally() would be nicer
        this.generating = false
      })
  }

  private setPrivacy(v: boolean) {
    this.setState({ showPrivacySaved: false, setPrivacyErr: "" })
    let success = false
    const genericErr = "Failed to toggle privacy. Try again?"

    fetch(`/setPrivacy?privacy=${v.toString()}`, { method: "POST" })
      .then(
        r => {
          if (r.status == 200) {
            window.clearTimeout(this.savedTimeout)
            this.setState({ showPrivacySaved: true })
            this.savedTimeout = window.setTimeout(() => {
              this.setState({ showPrivacySaved: false })
            }, 1000)
          } else {
            console.log("failed to update privacy: status=%d", r.status)
            this.props.accountChange({ ...this.props.account, private: !v })
            this.setState({
              private: !v, // revert
              setPrivacyErr: (r.status == 401) ? cookieAuthErrorMessage : genericErr,
            })
          }
        },
        err => {
          console.error(err)
          this.props.accountChange({ ...this.props.account, private: !v })
          this.setState({
            private: !v, // revert
            setPrivacyErr: genericErr,
          })
        }
      )
  }

  private onRegenerateClick(e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) {
    e.preventDefault()
    if (this.generating) { return }
    this.newAPIKey()
  }

  private onPrivacyClick(e: React.MouseEvent<HTMLInputElement, MouseEvent>) {
    let v = this.privacyCheckbox!.checked
    this.props.accountChange({ ...this.props.account, private: v })
    this.setState({ private: v })
    this.setPrivacy(v)
  }

  render() {
    let errClass = (b: string) => b ? "error" : "error hidden"
    let savedClass = (b: boolean) => this.state.showPrivacySaved ? "success" : "success hidden"

    return <div className="account">
      <table>
        <tbody>
          <tr><td>Username:</td><td>{this.props.account.username}</td></tr>
          <tr>
            <td>API Key:</td>
            <td className="mono apiKey">{this.state.apiKey}</td>
            <td><a href="" title="Note: Old API keys will no longer work after regenerating" onClick={this.onRegenerateClick.bind(this)}>Regenerate</a></td>
            <td><span className={errClass(this.state.keyGenerateErr)}>{this.state.keyGenerateErr}</span></td>
          </tr>
          <tr>
            <td>Private:</td>
            <td>
              <input type="checkbox" defaultChecked={this.state.private} onClick={this.onPrivacyClick.bind(this)} ref={r => { this.privacyCheckbox = r }} />
            </td>
            <td className={savedClass(this.state.showPrivacySaved)}>Saved</td>
            <td></td>
            <td className={errClass(this.state.setPrivacyErr)}>{this.state.setPrivacyErr}</td>
          </tr>
        </tbody>
      </table>
    </div>
  }
}
