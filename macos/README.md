![screenshot-menu](https://i.imgur.com/P210m2M.png)

### Making a new release

1. Increment `Version` in General > Identity
1. Product > Archive, then choose Distribute App > Copy App
1. Export the app to `macos/archive/<defaultname>`
1. Zip the `itunes-scrobble` app found at the export directory
1. `git tag $(date "+%s")` and `git push --tags`
1. Create a new release at https://github.com/nishanths/scrobble/releases.
   Use the tag as the title for the release. Include the zip file.
