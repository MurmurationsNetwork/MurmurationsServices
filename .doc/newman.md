# Installation Setup
1. Install [nvm](https://github.com/nvm-sh/nvm)
2. `nvm install --lts`
3. `nvm use --lts`
4. (Optional) Set the lts version as the default.
   `nvm alias default <type_version_by_yourself>`
5. `npm i -g newman`
6. (Optional) Check the npm packages
   `npm ls -g --depth=0`

# Run locally
`make newman-test`