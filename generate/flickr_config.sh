#!/bin/bash

#
# Generates flickr config file with our Api Key and Api Secret
#
# A compromise to keeping the secrets in the source code :o)
#

cat >$1 <<EOF

package cmd

const (
    FlickrApiKey = "${FLICKR_API_KEY}"
    FlickrApiSecret = "${FLICKR_API_SECRET}"
)

EOF
