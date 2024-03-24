#!/bin/bash

# Script Name: ubuntu_setup.sh
# Description: This script will update Ubuntu packages and install essential tools
# including vim, git, and htop.

# Ensure the script is run as root
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

# Update and Upgrade the System
echo "Updating and upgrading system packages..."
DEBIAN_FRONTEND=noninteractive apt-get update && DEBIAN_FRONTEND=noninteractive apt-get upgrade -y

# Install essential tools
echo "Installing vim, git, and htop..."
apt-get install -y vim git htop

# Installation open-iscsi for longhorn.
apt-get install -y open-iscsi

echo "Ubuntu setup is complete. Essential packages have been installed."
