# Debugging Guide

## Table of Contents

- [Set Up Kibana Error Filters](#set-up-kibana-error-filters)
- [Fix Failed Jobs](#fix-failed-jobs)
- [Fix Issues with Message Queues](#fix-issues-with-message-queues)

## Set Up Kibana Error Filters

If something goes wrong with the servers and you need to check the error logs quickly, setting up a filter in Kibana can save you a lot of time.

Here’s a simple step-by-step guide to make things smoother:

**1. Switch to the Right Context**:

Open your terminal and switch to the correct Kubernetes context with the following command:

```bash
kubectl config use-context <target-k8s-context>
```

**2. Connect Kibana to Port 5601**:

Next, you’ll want to make Kibana accessible on your local machine. Run this command:

```bash
kubectl port-forward -n murm-logging svc/murm-logging-kibana 5601:5601
```

**3. Open Kibana in Your Browser**:

Go to [Kibana discover page](http://localhost:5601/app/discover).

**4. Create a Filter**:

Once you're in, set up a filter. This will help you focus on the specific logs you’re interested in.

![Kibana Add Filter](./assets/images/kibana-add-filter.png)

**5. Pin the Filter**:

Don’t forget to pin your filter. This keeps it active and visible, so you don't lose track of it while you navigate around.

![Kibana Pin Filter](./assets/images/kibana-pin-filter.png)

**6. Repeat as Needed**:

You can set up multiple filters depending on what you need to monitor. Just repeat the steps above for each new filter.

## Fix Failed Jobs

1. **Check the Namespace**: Use k9s to look at the default namespace. Look for any jobs that didn't work out.

2. **View Job Logs**: To get more info, check the logs of the job that didn't work. In k9s, press `l` to see the logs of what you've selected.

3. **See All Jobs**: In k9s, type `:jobs`. This shows all the jobs so you can see details like their status, how old they are, and more.

    ![k9s Search Jobs](./assets/images/k9s-search-jobs.png)

4. **Remove Old or Failed Jobs**: After finding jobs that are old or didn't succeed, you can get rid of them in k9s. Click on the job and press `command + d`. You'll be asked to confirm that you want to delete it.

    ![k9s Delete Jobs](./assets/images/k9s-delete-jobs.png)

## Fix Issues with Message Queues

Sometimes, NATS might not work right. If that happens, these steps can help reset everything by removing and letting Kubernetes (k8s) start the stateful sets again.

1. **Go to Namespace**: Use `:namespaces` in k9s to go to the namespace page.
2. **Remove the Message Queue Namespace**: Look for and select the `murm-queue` namespace.

    ![k9s Namespaces Murm Queue](./assets/images/k9s-namespaces-murm-queue.png)

3. **Delete Each Stateful Set**: Choose each NATS stateful set one by one and use `command + d` to delete them. You don't need to wait for one to restart before deleting the next.

    ![k9s Nats](./assets/images/k9s-nats.png)
