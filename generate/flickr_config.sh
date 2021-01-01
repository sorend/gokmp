#!/bin/bash

#
# Generates flickr config file with our Api Key and Api Secret
#
# A compromise to keeping the secrets in the source code :o)
#

set -e

if [[ -z "$FLICKR_API_KEY" ]]; then
	echo "FLICKR_API_KEY needed in environment."
	exit 1
fi

if [[ -z "$FLICKR_API_SECRET" ]]; then
	echo "FLICKR_API_SECRET needed in environment."
	exit 1
fi

cat >$1 <<EOF

package cmd

const (
    FlickrApiKey = "${FLICKR_API_KEY}"
    FlickrApiSecret = "${FLICKR_API_SECRET}"
)

EOF
