import React from "react";
import { Account, Notie } from "../shared/types"
import { cookieAuthErrorMessage } from "../shared/util"
import "../scss/account"

declare let notie: Notie

interface AccountDetailProps {
  account: Account
}

interface AccountDetailState {
  apiKey: string
  private: boolean
}

export class AccountDetail extends React.Component<AccountDetailProps, AccountDetailState> {
  private generating = false // inflight request to generate API key?

  constructor(props: AccountDetailProps) {
    super(props)
    this.state = {
      apiKey: this.props.account.apiKey,
      private: this.props.account.private,
    }
  }

  private newAPIKey() {
    this.generating = true
    let success = false

    const genericErr = "Failed to regenerate API Key. Try again?"

    fetch(`/newAPIKey`, { method: "POST" })
      .then(res => {
        if (res.status == 200) {
          success = true
          return res.json()
        }
        if (res.status == 401) {
          notie.alert({ type: "error", text: cookieAuthErrorMessage })
          return res.blob()
        }
        console.log("failed to generate API key: status=%d", res.status)
        notie.alert({ type: "error", text: genericErr })
        return res.blob()
      })
      .then(
        r => {
          if (!success) { return }
          const apiKey = r as string
          notie.alert({ type: "success", text: "Regenerated API key." })
          this.setState({ apiKey })
        },
        err => {
          console.error(err)
          notie.alert({ type: "error", text: genericErr })
        }
      ).then(() => { // TODO: finally() would be nicer
        this.generating = false
      })
  }

  private setPrivacy(v: boolean) {
    const genericErr = "Failed to update privacy. Try again?"

    fetch(`/setPrivacy?privacy=${v.toString()}`, { method: "POST" })
      .then(
        r => {
          if (r.status == 200) {
            notie.alert({ type: "success", text: "Updated privacy." })
            this.setState({
              private: v,
            })
          } else {
            console.log("failed to update privacy: status=%d", r.status)
            notie.alert({ type: "error", text: (r.status == 401) ? cookieAuthErrorMessage : genericErr })
          }
        },
        err => {
          console.error(err)
          notie.alert({ type: "error", text: genericErr })
        }
      )
  }

  private onRegenerateClick(e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) {
    e.preventDefault()
    if (this.generating) { return }
    this.newAPIKey()
  }

  private privacyOn() {
    return <>
      <p>Your profile is private. You can <a href="" onClick={(e) => { e.preventDefault(); this.setPrivacy(false) }}>make your profile publicly viewable.</a></p>
    </>
  }

  private privacyOff() {
    return <>
      <p>Your profile is publicly viewable. You can <a href="" onClick={(e) => { e.preventDefault(); this.setPrivacy(true) }}>make your profile private.</a></p>
    </>
  }

  render() {
    return <div className="account">
      <p>
        Your API key is <span className="mono apiKey">{this.state.apiKey}</span>.&nbsp;
        <a href="" title="Note: Old API keys will no longer work after regenerating" onClick={this.onRegenerateClick.bind(this)}>Regenerate API key.</a>
      </p>
      {this.state.private ? this.privacyOn() : this.privacyOff()}
    </div>
  }
}
