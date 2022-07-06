# fan-gopher
A CLI used to scrape information from OnlyFans, written in go.

This is to be used by other utilities that may need to be able to scrape this data, but are unable to handle the SPA / JavaScript implementation of OnlyFans.

## Usage

This can be used to verify that a particular post exists and also to get the details of a particular post.

On first run, you may see output like the following:

```
[launcher.Browser]2022/07/05 20:20:18 try to find the fastest host to download the browser binary
[launcher.Browser]2022/07/05 20:20:18 check https://storage.googleapis.com/chromium-browser-snapshots/Mac/1010848/chrome-mac.zip
[launcher.Browser]2022/07/05 20:20:18 check https://registry.npmmirror.com/-/binary/chromium-browser-snapshots/Mac/1010848/chrome-mac.zip
[launcher.Browser]2022/07/05 20:20:18 check https://playwright.azureedge.net/builds/chromium/1010848/chromium-linux-arm64.zip
[launcher.Browser]2022/07/05 20:20:18 check result: Get "https://registry.npmmirror.com/-/binary/chromium-browser-snapshots/Mac/1010848/chrome-mac.zip": context canceled
[launcher.Browser]2022/07/05 20:20:18 check result: Get "https://playwright.azureedge.net/builds/chromium/1010848/chromium-linux-arm64.zip": context canceled
[launcher.Browser]2022/07/05 20:20:18 Download: https://storage.googleapis.com/chromium-browser-snapshots/Mac/1010848/chrome-mac.zip
[launcher.Browser]2022/07/05 20:20:18 Progress:
[launcher.Browser]2022/07/05 20:20:18 00%
[launcher.Browser]2022/07/05 20:20:19 16%
[launcher.Browser]2022/07/05 20:20:20 42%
[launcher.Browser]2022/07/05 20:20:21 67%
[launcher.Browser]2022/07/05 20:20:22 91%
[launcher.Browser]2022/07/05 20:20:23 100%
[launcher.Browser]2022/07/05 20:20:23 Unzip to: /Users/FansForYou/.cache/rod/browser/chromium-1010848
[launcher.Browser]2022/07/05 20:20:23 Progress:
[launcher.Browser]2022/07/05 20:20:23 00%
[launcher.Browser]2022/07/05 20:20:24 50%
[launcher.Browser]2022/07/05 20:20:25 98%
[launcher.Browser]2022/07/05 20:20:25 100%
```

If you are parsing out the console output JSON, it's recommended you exclude invalid JSON from your parsing.

### Verification

You can use the `-verify` flag to verify if a post exists, like so:

```
./fan-gopher -postID=1234 -creatorName=onlyFansUsername -verify=true`
```

Its output will look like the following if the post is found:

```
{
    "error": false,
    "errorCode": "",
    "errorMessage": ""
}
```

If the post is _not_ found, it will look like:

```
{
    "error": true,
    "errorCode": "POST_NOT_FOUND",
    "errorMessage": "<descriptive error message>"
}
```

### Post Details

You can use the `-get-details` flag to retrieve the details of a post, like so:

```
./fan-gopher -postID=1234 -creatorName=onlyFansUsername -get-details=true
```

It will return the following should the post be found:

```
{
   "error":false,
   "errorCode": "",
   "errorMessage":"",
   "actorDetails":[
      {
         "actorName":"Cool Person",
         "profileImageUrl":"https://public.onlyfans.com/files/thumbs/c50/0/07/07p/07phtwjnst3thcpnkndwittbix6l8yhd1539574274/avatar.jpg"
      }
   ],
   "videoDetails":{
      "videoDescription":"This is the body of the post"
   }
}
```

## Releasing

Use the following to build releases of this project:

```
make version=<version goes here> release
```

This will build artifacts with a version identifier for 64-bit Linux, Mac, and Windows systems in the `dist/` folder.