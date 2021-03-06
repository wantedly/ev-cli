#!/usr/bin/env bash

set -eu
set -o pipefail

# check token
: ${WANTEDLY_HOME:=$HOME/.wantedly}
TOKENPATH=$WANTEDLY_HOME/credentials/github-access-token
: ${TOKEN:=${GITHUB_TOKEN:-${GITHUB_ACCESS_TOKEN:-""}}}
: ${TOKEN:=$([ -f $TOKENPATH ] && cat $TOKENPATH)}
if [ -z "$TOKEN" ]; then
	echo -e "You need to set \$GITHUB_ACCESS_TOKEN.
Go to the link below to create a new compatible token
\033[4mhttps://github.com/settings/tokens/new?scopes=repo,read:org&description=wantedly+ev-cli\033[0m"
	exit 1
fi

# check system
: ${OS:=$(uname -s | tr '[A-Z]' '[a-z]')}
: ${ARCH:=amd64}

# check target
: ${EV_VERSION:=latest}
: ${EXTENSION:=.tar.gz}
: ${DEST:=$WANTEDLY_HOME/bin}
if [ "$EV_VERSION" = "latest" ]; then
	RELEASE_PATH="latest"
else
	RELEASE_PATH="tags/v$EV_VERSION"
fi

# download
echo "Installing ev-cli $EV_VERSION for $OS-$ARCH to $DEST"
query=".assets[] | select(.name | contains(\"$OS\") and contains(\"$ARCH\") and contains(\"$EXTENSION\")) | .id"
id=$(curl -s https://$TOKEN@api.github.com/repos/wantedly/ev-cli/releases/$RELEASE_PATH | jq "$query")
mkdir -p $DEST
curl -sLJ -H 'Accept: application/octet-stream' https://$TOKEN@api.github.com/repos/wantedly/ev-cli/releases/assets/$id | tar xz -C $DEST --strip=1 "$OS-$ARCH/ev"
echo "ev-cli has successfully been installed to $DEST"
echo "You might want to add the line below to the shell profile"
echo "  export PATH=${DEST/$HOME/\$HOME}:\$PATH"
