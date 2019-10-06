import React, { useState, useRef, useEffect } from "react";
import { Account } from "../shared/types"
import "../scss/username.scss";

interface SetUsernameProps {
  accountChange: (a: Account) => void
}

export const SetUsername: React.FC<SetUsernameProps> = ({
  accountChange,
}) => {
  const initialMount = useRef(true)
  const [username, setUsername] = useState("")
  const [error, setError] = useState("")
  let input: HTMLInputElement | null = null

  useEffect(() => {
    if (initialMount.current === false) { return }
    initialMount.current = false
    if (input !== null) { input.focus() }
  })

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    initializeAccount(username)
  }

  const initializeAccount = (u: string): void => {
    let { reason, ok } = validateUsername(u)
    if (!ok) {
      setError(reason)
      return
    }

    setError("")
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
          setError("The username is already taken")
          return res.text()
        }
        setError(genericError)
        return res.blob()
      })
      .then(r => {
        if (!success) { return }
        accountChange(r as Account)
      }, err => {
        console.error(err)
        setError(genericError)
      })
  }

  return <>
    <form onSubmit={(e) => { handleSubmit(e) }}>
      <label>
        Set your username:
        <input type="text" value={username} onChange={() => { setUsername(input!.value) }}
          ref={r => { input = r }}></input>
      </label>
      <input type="submit" value="OK" />
      {error !== "" && <span className="error">{error}</span>}
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
