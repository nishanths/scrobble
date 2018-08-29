interface Account {
  username: string
  apiKey: string
}

export interface BootstrapArgs {
  host: string
  email: string
  loginURL: string
  logoutURL: string
  downloadURL: string
  account: Account
}
