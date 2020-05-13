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
  private: boolean
}

export type ArtworkHash = string

export interface Song {
  albumTitle: string
  artistName: string
  title: string
  totalTime: number
  year: number
  releaseDate: number
  lastPlayed: number
  playCount: number
  artworkHash: ArtworkHash
  trackViewURL: string
  loved: boolean

  ident: string
}

export interface ScrobbledResponse {
  total: number
  songs: Song[]
}

export type NProgress = {
  start(): void
  done(): void
  configure(opts: { [k: string]: any }): void
}

export type Notie = {
  alert(options: {
    type?: 'success' | 'warning' | 'error' | 'info' | 'neutral' // optional, default = 'info'
    text: string
    stay?: boolean // optional, default = false
    time?: number // optional, default = 3, minimum = 1
    position?: 'top' | 'bottom' // optional, default = 'top'
  }): void
}
