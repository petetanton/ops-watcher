# OPS Watcher
A tool that sends MacOS push notifications for updates to JIRA tickets

## Installation with homebrew
Enable the custom tap
```bash
brew tap petetanton/ops-tools
==> Tapping petetanton/ops-tools
Cloning into '/usr/local/Homebrew/Library/Taps/petetanton/homebrew-ops-tools'...
remote: Enumerating objects: 4, done.
remote: Counting objects: 100% (4/4), done.
remote: Compressing objects: 100% (2/2), done.
remote: Total 4 (delta 0), reused 4 (delta 0), pack-reused 0
Unpacking objects: 100% (4/4), done.
Tapped 1 formula (27 files, 24.7KB).
```
Install ops-watcher
```bash
brew install ops-watcher
==> Installing ops-watcher from petetanton/ops-tools
==> Downloading https://github.com/petetanton/ops-watcher/releases/download/0.0.1/ops-watcher-darwin-amd64.zip
🍺  /usr/local/Cellar/ops-watcher/64: 3 files, 12.2MB, built in 4 seconds
```

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
