#! /bin/bash
set -e

sudo apt-get update -y
sudo apt-get install -y curl unzip
curl -L "${download-url}" > /tmp/approzium.zip

cd /tmp
sudo unzip approzium.zip
sudo mv authenticator /usr/local/bin/authenticator
sudo chmod 0755 /usr/local/bin/authenticator
sudo chown root:root /usr/local/bin/authenticator

# Setup the configuration
cat <<EOF > /tmp/authenticator-config
${config}
EOF
sudo mv /tmp/authenticator-config /usr/local/etc/approzium.config.yml

# Setup the init script
cat <<EOF > /tmp/upstart
description "Approzium Authenticator server"
start on runlevel [2345]
stop on runlevel [!2345]

respawn

script
  if [ -f "/etc/service/approzium" ]; then
    . /etc/service/approzium
  fi

  exec /usr/local/bin/authenticator \
    --config="/usr/local/etc/" \
    >> "${logs-file}" 2>&1
end script
EOF
sudo mv /tmp/upstart /etc/init/approzium.conf

# Extra install steps (if any)
${extra-install}

sudo start approzium
