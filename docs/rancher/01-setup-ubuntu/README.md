# How to Set Up and Secure Your Ubuntu Server for Development

## Introduction

This guide will show you how to prepare your Ubuntu server for development by setting up secure access with SSH keys, modifying your SSH configuration for enhanced security, and installing necessary packages through a script.

After completing this guide, you'll be able to:

- Securely log into your server using SSH keys.
- Configure your SSH settings for improved security.
- Install essential packages on your server using a script.

## Table of Contents

1. [Introduction](#introduction)
2. [Prerequisites](#prerequisites)
3. [Step 1 - Generating and Setting Up SSH Keys](#step-1---generating-and-setting-up-ssh-keys)
4. [Step 2 - Adjusting Your SSH Settings for Security](#step-2---adjusting-your-ssh-settings-for-security)
5. [Step 3 - Installing Packages with a Script](#step-3---installing-packages-with-a-script)
6. [Conclusion](#conclusion)

## Prerequisites

Before starting, you need:

1. An Ubuntu server with root access.
2. A computer with SSH client software.

## Step 1 - Generating and Setting Up SSH Keys

Create a secure SSH key pair on your computer:

```bash
ssh-keygen -t rsa -b 4096
```

Copy your public SSH key to your clipboard for easy access:

```bash
cat ~/.ssh/id_rsa.pub | pbcopy
```

Log into your server as the root user:

```bash
ssh root@<ip_address>
```

Set up the `.ssh` directory and secure it:

```bash
mkdir -p ~/.ssh && chmod 700 ~/.ssh
```

Paste your SSH key into the `authorized_keys` file:

```bash
vim ~/.ssh/authorized_keys
```

Then secure the file:

```bash
chmod 600 ~/.ssh/authorized_keys
```

**Important:** Test your SSH key by logging in from another terminal tab or session before proceeding to ensure you don't lose access.

```bash
ssh root@<ip_address>
```

## Step 2 - Adjusting Your SSH Settings for Security

Enhance your SSH login security by modifying a couple of settings:

Open the SSH configuration file:

```bash
vim /etc/ssh/sshd_config
```

Include these lines to disable password authentication and deny empty passwords:

```
PasswordAuthentication no
PermitEmptyPasswords no
```

Save your changes and exit the editor.

Apply your new SSH configuration by restarting the service:

```bash
systemctl restart ssh
```

## Step 3 - Installing Packages with a Script

In your Ubuntu server, run the following script to install necessary packages. This script fetches and applies an installation script directly:

```bash
curl -sSL https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsServices/main/scripts/ubuntu_setup.sh | bash
```

## Conclusion

Your Ubuntu server is now set up with enhanced security measures and ready for development. You've implemented SSH key access, updated SSH configurations for improved security, and installed important packages.
