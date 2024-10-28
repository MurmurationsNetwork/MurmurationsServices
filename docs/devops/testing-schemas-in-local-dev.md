# Testing Schemas in Your Local Development Environment

You can test out new fields and schemas in your local development environment before submitting them for inclusion to a library repo.

## Getting Started with Local Setup

1. Initialization: Execute the `make dev` command to run your local environment.
2. Schema Fetching: By default, the `schemaparser` cronjob is configured to fetch schemas simultaneously from both local directories and remote git repositories, with schemas from local directories having priority over remote repositories.

## Adding Custom Schemas

1. Place your custom schema files within the `library/schemas` directory. This is the designated location for all custom schema files.
2. If your schema includes fields with `$ref` references, ensure that these referenced files are located in the `library/fields` directory. This organization is crucial to avoid errors during the cronjob execution.

## `schemaparser` Cronjob Execution

1. Frequency: The `schemaparser` runs automatically [every minute](/charts/murmurations/charts/schemaparser/templates/schemaparser/cronjob.yaml) in the development environment, and every two minutes in all other environments. This frequent execution ensures that your schemas are continuously updated as you make changes to and test them.
2. Manual Triggering: If you prefer not to wait for the automatic cycle, you can manually trigger the `schemaparser` job. Use the command below, but remember to modify the number (0001) each time, as cronjob names must be unique:

   ```bash
   kubectl create job --from=cronjob/schemaparser-app schemaparser-app-0001
   ```

3. Delete Job: After the job has completed, you can delete it using the command below:

   ```bash
   kubectl delete job schemaparser-app-0001
   ```
