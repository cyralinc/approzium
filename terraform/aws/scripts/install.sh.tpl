#! /bin/bash
set -e

sudo apt-get update -y
sudo apt-get install -y curl unzip
curl -o approzium.zip -LO "${download-url}"

sudo unzip approzium.zip
sudo mv authenticator /usr/local/bin/authenticator
sudo chown root:root /usr/local/bin/authenticator

# Setup the configuration
cat <<EOF > approzium.config
${config}
EOF
sudo mv approzium.config /usr/local/etc/approzium.config

# Setup the init script
cat <<EOF > approzium.conf
description "Approzium Authenticator server"
start on runlevel [2345]
stop on runlevel [!2345]

respawn

script
  if [ -f "/etc/service/approzium" ]; then
    . /etc/service/approzium
  fi

  exec /usr/local/bin/authenticator \
    --config="/usr/local/etc/approzium.config" \
    >> "${logs-file}" 2>&1
end script
EOF
sudo mv approzium.conf /etc/init/approzium.conf

# Extra install steps (if any)
${extra-install}

sudo start approzium
