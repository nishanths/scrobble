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

export interface ScrobbledResponse {
  total: number
  songs: Song[]
}
