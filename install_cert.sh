#!/usr/bin/env bash

# Install a self-signed cert for the test email server in localhost.
# This is obviously OS dependent, so this is for OSX only. Use it under your own risk.
sudo security add-trusted-cert -d -r trustRoot -k ~/Library/Keychains/login.keychain internal/email/testdata/server.crt