# prolice-slack-bot

Slack bot that tracks and reminds of pull requests

## How to use

Once PRolice bot is active in your channel, urls posted that fit the matching strings in the .env file and don't contain reserve words (remove) will be tracked.

These Pull requests will be reposted 3 times daily to remind members to review. The posts will include the URL, Author, Posted Date, and Approvals (if any).

### Important Notes

Only published (non drafts), and active (not complete or abandoned) pull requests will be tracked. If a Pull request reaches a terminal state, or is deleted it will be removed from the list of pull requests.

## Commands

### Mention commands (prefaced with @PRolice)

- list - Lists all Pull Requests currently being tracked
- remove (url) - removes pull request matching specified url from list of tracked pull requests
- (empty) - lists available commands

## Updates
### 12/12/2020

#### Bug fixes
- Fixes bug where slack a shortened url (indicated by 3 dots in the url) was not being tracked
- Fixes bug where an already tracked PR was being retracked if it was posted again with a modified url
- Fidex bug where a message with multiple PR's would only have the first PR tracked

#### Enhancements
- Added alternate color coding for messages
- Added functionality to handle PR's posted as named links.
