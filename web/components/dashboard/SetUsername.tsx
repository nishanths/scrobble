import React, { useState, useEffect } from "react";
import { Account, Notie } from "../../shared/types"
import { cookieAuthErrorMessage } from "../../shared/util"
import "../../scss/dashboard/set-username.scss"

interface SetUsernameProps {
  accountChange: (a: Account) => void
}

declare const notie: Notie

export const SetUsername: React.FC<SetUsernameProps> = ({ accountChange }) => {
  const [username, setUsername] = useState("")
  const [error, setError] = useState("")

  let inputRef: HTMLInputElement | null = null

  useEffect(() => {
    if (inputRef !== null) { inputRef.focus() }
  }, [inputRef])

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    initializeAccount(username)
  }

  const initializeAccount = (u: string): void => {
    setError("")

    const { reason, ok } = validateUsername(u)
    if (!ok) {
      setError(reason)
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
          setError("The username is already taken.")
          return res.text()
        }
        if (res.status == 401) {
          setError(cookieAuthErrorMessage)
          return res.text()
        }
        setError(genericError)
        return res.blob()
      })
      .then(r => {
        if (!success) { return }
        notie.alert({ type: "success", text: "Username set!" })
        accountChange(r as Account)
      }, err => {
        console.error(err)
        setError(genericError)
      })
  }

  return <div className="SetUsername">
    <form spellCheck="false" onSubmit={(e) => { handleSubmit(e) }}>
      <label>

        <input type="text"
          value={username}
          onChange={() => { setError(""); setUsername(inputRef!.value) }}
          ref={r => { inputRef = r }}>
        </input>

        {error === "" ?
          <><div className="label">Allowed characters: lowercase a-z and 0-9.</div></> :
          <><div className="label error">{error}</div></>}
      </label>
      <input type="submit" value="⮐" />
    </form>
  </div>
}

const usernameCharRe = /^[a-z0-9]*$/

const validateUsername = (u: string): ({ reason: "", ok: true } | { reason: string, ok: false }) => {
  if (u.length < 2) {
    return { reason: "Username must be at least 2 characters long.", ok: false }
  }
  if (u.length > 24) {
    return { reason: "Username cannot be more than 24 characters long.", ok: false }
  }
  if (!usernameCharRe.test(u)) {
    return { reason: "Username can contain only lowercase a-z and 0-9.", ok: false }
  }
  if (u.indexOf("scrobble") != -1) {
    return { reason: "Username cannot contain the word 'scrobble'.", ok: false }
  }
  if (u === "username") {
    return { reason: "Username cannot be the word 'username'.", ok: false }
  }
  return { reason: "", ok: true }
}
