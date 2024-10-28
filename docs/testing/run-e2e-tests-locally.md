# How to Run End-to-End (E2E) Tests in Your Local Development Environment

We use [Newman](https://github.com/postmanlabs/newman) to run the E2E tests.

1. Install NodeJS using [nvm](https://github.com/nvm-sh/nvm?tab=readme-ov-file#install--update-script)
2. Install latest Node LTS and set it as the default

    ```bash
    nvm install --lts
    nvm use --lts
    nvm alias default <type_version_by_yourself>
    ```

3. Install Newman as a global package

    ```bash
    npm i -g newman
    ```

4. Start the Murmurations Services

    ```bash
    make dev
    ```

5. Run Newman locally to initiate the E2E tests

    ```bash
    make newman-test
    ```
