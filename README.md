
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

Download from github releases:

```bash
# download binary
curl ... -o gokmp
chmod +x gokmp
```

Login to Flickr with the tool. You can complete the flow on your PC if you are running from a headless server.
Login creates a file `.gokmp.yaml` with your tokens. If you loose this file, you need to run `gokmp login` again.

```bash
# login to Flickr
./gokmp login

Opening browser for URL:
   https://www.flickr.com/services/oauth/authorize?oauth_token=72REDACTED0684227-74d9REDACTEDff66&perms=read
   (please do it manually if opening browser fails)

Enter verifier PIN here: 842-000-810

Access Token : 721REDACTED251363-7c7REDACTED2e1fd
Access Secret: ffecREDACTED4624
```

After login you can start the backup:

```bash
# start backup
./gokmp backup -i "59362368@N00" -d ~/flickr_backup
Backing up ...
921 photos found in /home/sorend/flickr_backup
Looking for photos to backup in 390 photosets...
 .. Denmark 2020 November .. queue 1
 .. Denmark 2020 October .. queue 1
 .. Denmark 2020 September .. queue 1
 .. Denmark 2020 August .. queue 112

...

All done :-)
```

If gokmp fails (usually due to Flickr timeout or application error), you can simply start it again.
You should not run multiple gokmp instances simultaneously as they will compete over completing the
same images.


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
