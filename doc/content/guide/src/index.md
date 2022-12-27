# Guide

* [What is Scrobble?](#what-is-scrobble)
* [How do I sign up?](#how-do-i-sign-up)
* [How do I scrobble my songs?](#how-do-i-scrobble-my-songs)
* [How do I see scrobbled songs?](#how-do-i-see-scrobbled-songs)
* [Will my listening history on other devices (eg. iPhone) be scrobbled?](#will-my-listening-history-on-other-devices-eg-iphone-be-scrobbled)
* [Is there a HTTP API?](#is-there-a-http-api)
* [Limitations](#limitations)
* [Contact/Report issues](#contact-report-issues)


# What is _Scrobble_?

Scrobble is a music scrobbling service for Apple Music — with a focus on showcasing
album art.

The term 'scrobbling' refers to sharing the details of songs you listen to, to an online feed.

Scrobble's macOS app can scrobble what you listen to. Your listening
history can then be seen by you, and optionally others, on the Scrobble website.

# How do I sign up?

Visit https://scrobble.growl.space, and click "Sign in with Google".

Follow the prompt to set a username.

# How do I scrobble my songs?

Get the [macOS app](https://github.com/nishanths/scrobble/releases/latest). (It's unintrusive and lives in the menu bar.)

Open the app, and click on "Start scrobbling..." from the menu bar.

Follow the prompt to enter your API key, which can be found in the [dashboard](/dashboard/api-key) when you're signed in.

You're all set to scrobble!

![Start scrobbling...](/doc/guide/macos_start_scrobbling.png)

# How do I see scrobbled songs?

You can see scrobbled songs on the website.

A user's scrobbled songs can be seen at `https://scrobble.growl.space/u/<username>`.

Here's an example: https://scrobble.growl.space/u/nishanth.

# Will my listening history on other devices (eg. iPhone) be scrobbled?

Yes, when you use the same Apple Music account on that device as you do on your Mac.

# Is there a HTTP API?

Yes. Check out the [API documentation](/doc/api/v1).

# Limitations

Scrobbling isn't instantaneous due to limitations in how soon your Apple Music listening history is available
in the Apple Music desktop app. Additionally the Scrobble macOS app by default scrobbles only once every 24 hours.
So your listening history on the Scrobble website may be behind by a day or more.

Songs have to be present in your Apple Music library to be scrobbled and to appear on the website.

# Contact/Report issues

Email [bambooparch@gmail.com](mailto:bambooparch@gmail.com), or create an issue
in the [issue tracker](https://github.com/nishanths/scrobble/issues).
