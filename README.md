
# gokmp -- go keep my photos

```text
                    ___
     _______ ______/   / ___ _  ___________
    / __   // _   //  / /  /  \/ __\    __ \
   /  |/   /  /   /  /_/  /   ______\   |/  \
  /   /   /  /   /     __/       ___/   /   /
  \__    /______/__/\   \ _/\__/  _/   ____/  Go Keep My Photos
==__/   / ========== \___\ == /___/   / =========================
 /_____/                         /___/
```

gokmp is a Go conversion of my previous [keepmyphotos](https://github.com/sorend/keepmyphotos) backup tool in Python.

Written mostly to get experience with Go programming.

## Usage

Download from github releases and run it:

```bash
$ curl


```


### Flickr API key and secret

The binary available under releases has an API key built in. If you wish to bake
a binary with your own API key (for developing on gokmp locally), you need to
have an API key available in the environment

You can get it from the App Garden: https://www.flickr.com/services/apps/create/

When building you need to have them in the environment, as `go generate` will use them:

```bash
$ export FLICKR_API_KEY=myapikey FLICKR_API_SECRET=myapisecret
$ go generate
$ cat cmd/flickr_config

package cmd

const (
    FlickrApiKey = "myapikey"
    FlickrApiSecret = "myapisecret"
)
```
