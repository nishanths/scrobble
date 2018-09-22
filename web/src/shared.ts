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

export function base64Encode(s: string) {
  let bytes = new TextEncoder().encode(s);
  return base64js.fromByteArray(bytes);
}

export function base64Decode(s: string) {
  let bytes = base64js.toByteArray(s);
  return new TextDecoder("utf-8", { fatal: true }).decode(bytes);
}
