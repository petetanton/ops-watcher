# OPS Watcher
A tool that sends MacOS push notifications for updates to JIRA tickets

## Configuration
The app looks for an `ops-watcher.yaml` file which should contain the following:
```yaml
jira_username: my.email@mydomain.com
jira_password: jiraapikey
jira_baseurl: https://example.atlassian.net
jira_enabled: true
jira_query:
  - <as many JQL queries as you like without and filters by date>
  - watcher = currentUser()
```

When the app runs, it queries the JIRA API every minute for each JQL query in the configuration.  
If a new issue has been raised or updated for any of those queries since the last run, a Mac OS push notification is sent using the `terminal-notification` command.
