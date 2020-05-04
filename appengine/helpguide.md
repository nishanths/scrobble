## What is this?

_Scrobble_ is a music scrobbling service that works with Apple Music.

The term scrobbling refers to the act of publishing the songs you listen to to an online feed.

The service consists of two parts:

1. a macOS application which uploads your Apple Music listening history; and
2. this web application which displays your list of scrobbled songs.

## How do you get started?

To begin scrobbling songs:

1. visit the web app at https://scrobble.allele.cc;
1. sign in with Google and pick your username;
1. copy the API key;
1. download the [macOS app][macos];
1. enter the copied API key in the macOS app.

Your scrobbled songs can be seen at `https://scrobble.allele.cc/u/username`.

## The profile page

The profile page displays scrobbled songs.

You can filter by "All" (which displays all scrobbled songs), "Loved" (which displays songs liked on Apple Music), or "By color" (which displays scrobbled songs by artwork color).

![segmented control](http://imgur.com/Cz2v327l.png)

## The macOS menu bar application

The macOS app, which scrobbles your Apple Music listening history, can be downloaded [here][macos].

Click on the app icon in the menu bar, and choose "Start scrobbling...". Enter your API key in the dialog to begin scrobbling.

<img src="https://imgur.com/VTOFh02.png" height=300 />

## Can you scrobble without a macOS computer?

There is no easy way currently to scrobble your songs without a macOS computer. Please [create an issue][issues] to add scrobbling support to non-macOS computers.

You may, however, use the HTTP API to scrobble your songs. The API is not documented yet (!), but the [source code][source] is available.

You may still use the web application to browse songs that your friends have scrobbled.

## Will music you listen to on other devices (eg. iPhone) be scrobbled?

Yes, if you use the same Apple Music account on your iPhone and on your macOS computer. See the section below for details on how this works.

## Delay in displaying your latest songs

This is a known issue due to a limitation in when your latest listening history is available in the Apple Music desktop app and the iTunes app (pre-macOS Catalina).

The delay depends on a few factors:

* the Apple Music macOS app and iTunes app take roughly 12h to update play counts, including local play counts;
* when you've last opened the Apple Music macOS app (the app must be open to update play counts); and
* the macOS scrobble app's upload frequency of 24h.

## Is last.fm or similar support on the roadmap?

No. But please [create an issue][issues] if you're interested.

[issues]: https://github.com/nishanths/scrobble/issues/new
[source]: https://github.com/nishanths/scrobble
[macos]: https://github.com/nishanths/scrobble/releases/latest
