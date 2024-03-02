# Debugging Guide

## Table of Contents

- [Troubleshooting Failed Jobs](#troubleshooting-failed-jobs)

## Troubleshooting Failed Jobs

1. **Inspect the Namespace**: Utilize k9s to examine the default namespace. Keep an eye out for any jobs that have failed.

2. **Access Job Logs**: For a deeper investigation, access the logs of the failed job. In k9s, you can press `l` on your keyboard to view the logs of the selected resource.

3. **List All Jobs**: Within k9s, input the command `:jobs`. This will display a list of all jobs, providing you with the ability to inspect each job's details, such as status, age, and more.

    ![k9s Search Jobs](./assets/images/k9s-search-jobs.png)

4. **Delete Old or Failed Jobs**: Once you've identified old or failed jobs, they can be deleted directly from k9s. Select the job and press `command + d`. A prompt will appear for you to confirm the deletion.

    ![k9s Delete Jobs](./assets/images/k9s-delete-jobs.png)
