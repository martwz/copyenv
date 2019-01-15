[![Build Status](https://travis-ci.com/martinxxd/copyenv.svg?branch=master)](https://travis-ci.com/martinxxd/copyenv)

# Cloud Foundry CLI Copy Env Plugin

Cloud Foundry CLI plugin to export application VCAP_SERVICES and VCAP_APPLICATION.

## Install

```bash
curl -L https://github.com/martinxxd/copyenv/releases/download/v1.1.0/copyenv_1.1.0_darwin_amd64.tar.gz | tar -zxvf copyenv
cf install-plugin copyenv
```

## Usage

```bash
cf copyenv APP_NAME [--all] [--plain]

export VCAP_SERVICES='...'
```
