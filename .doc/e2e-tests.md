# How to run end-to-end tests

1. Install node with nvm

```bash
# https://github.com/nvm-sh/nvm#install--update-script
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
# follow instructions at link above to ensure nvm is added to your PATH

# Install latest node LTS
nvm install 16.13.1
```

2. Install newman as a global package

```bash
npm i -g newman
# Confirm it is installed
npm ls -g --depth=0
```

3. Run Newman

```bash
newman run e2e-tests.json -e e2e-dev-env.json
```
