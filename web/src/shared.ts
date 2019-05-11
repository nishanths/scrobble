import * as base64js from "base64-js";

export interface Account {
  username: string
  apiKey: string
  private: boolean
}

export interface BootstrapArgs {
  host: string
  email: string
  loginURL: string
  logoutURL: string
  account: Account
}

export interface UArgs {
  host: string
  artworkBaseURL: string
  profileUsername: string
  logoutURL: string
  account: Account
  self: boolean
}

export interface Song {
  albumTitle: string
  artistName: string
  title: string
  totalTime: number
  year: number
  lastPlayed: number
  artworkHash: string
  trackViewURL: string
  loved: boolean

  ident: string
}

export function hasPrefix(s: string, prefix: string): boolean {
  return s.length >= prefix.length && s.slice(0, prefix.length) == prefix
}

export function trimPrefix(s: string, prefix: string): string {
  if (hasPrefix(s, prefix)) {
    return s.slice(prefix.length)
  }
  return s
}

export function unreachable(): never {
  throw new Error("unreachable");
}

// base64Encode encodes the string using base64 standard encoding, i.e.,
// the encoding corresponding to base64.StdEncoding in Go.
//
// Verified using https://runkit.com/nishanths/5cd735892538b9001a7e08d5
// and https://gobyexample.com/base64-encoding.
export function base64Encode(s: string): string {
  let bytes = new TextEncoder().encode(s);
  return base64js.fromByteArray(bytes);
}

// base64Decode is the inverse of base64Encode.
export function base64Decode(s: string): string {
  let bytes = base64js.toByteArray(s);
  return new TextDecoder("utf-8", { fatal: true }).decode(bytes);
}

export function pathComponents(p: string): string[] {
  return p.split("/").filter(s => s != "")
}
