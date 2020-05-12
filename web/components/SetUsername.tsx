import React, { useState, useEffect } from "react";
import { Account, Notie } from "../shared/types"
import { cookieAuthErrorMessage } from "../shared/util"
import "../scss/set-username.scss";

interface SetUsernameProps {
  accountChange: (a: Account) => void
}

declare let notie: Notie

export const SetUsername: React.FC<SetUsernameProps> = ({ accountChange }) => {
  const [username, setUsername] = useState("")
  let input: HTMLInputElement | null = null

  useEffect(() => {
    if (input !== null) { input.focus() }
  }, [])

  const showError = (text: string) => {
    notie.alert({ type: "error", text })
  }

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    initializeAccount(username)
  }

  const initializeAccount = (u: string): void => {
    const { reason, ok } = validateUsername(u)
    if (!ok) {
      showError(reason)
      return
    }

    const genericError = "Something went wrong. Try again?"

    // TODO: clean up the control flow
    let success = false;
    fetch(`/initializeAccount?username=${u}`, { method: "POST" })
      .then(res => {
        if (res.status == 200) {
          success = true
          return res.json()
        }
        if (res.status == 406) {
          showError("The username is already taken")
          return res.text()
        }
        if (res.status == 401) {
          showError(cookieAuthErrorMessage)
          return res.text()
        }
        showError(genericError)
        return res.blob()
      })
      .then(r => {
        if (!success) { return }
        notie.alert({ type: "success", text: "Username set!" })
        accountChange(r as Account)
      }, err => {
        console.error(err)
        showError(genericError)
      })
  }

  return <>
    <form onSubmit={(e) => { handleSubmit(e) }}>
      <label>
        Set your username:
        <input type="text" size={15} value={username} onChange={() => { setUsername(input!.value) }}
          ref={r => { input = r }}></input>
      </label>
      <input type="submit" value="OK" />
    </form>
  </>
}

const usernameRe = /^[a-z0-9]*$/

const validateUsername = (u: string): ({ reason: "", ok: true } | { reason: string, ok: false }) => {
  if (!(u.length > 2)) {
    return { reason: "Username must be at least 2 characters", ok: false }
  }
  if (!(u.length < 24)) {
    return { reason: "Username must not be more than 24 characters", ok: false }
  }
  if (!usernameRe.test(u)) {
    return { reason: "Username must contain only a-z and 0-9", ok: false }
  }
  if (u.indexOf("scrobble") != -1) {
    return { reason: "Username cannot contain 'scrobble'", ok: false }
  }
  if (u === "username") {
    return { reason: "Username cannot be 'username'", ok: false }
  }
  return { reason: "", ok: true }
}
