# How to run end-to-end tests

1. Install Node with [nvm](https://github.com/nvm-sh/nvm)
2. Install latest node LTS

```bash
nvm install --lts
nvm use --lts
```

3. (Optional) Set the LTS version as the default
```bash
nvm alias default <type_version_by_yourself>
```

4. Install Newman as a global package

```bash
npm i -g newman
# (Optional) Check the npm packages
npm ls -g --depth=0
```

5. Run Newman locally

```bash
make newman-test
```
