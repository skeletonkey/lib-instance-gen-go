name: CHANGELOG.md has been updated

on:
  push:

jobs:
  check-file-change:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Fetch latest changes from main
        run: git fetch origin main

      - name: CHANGELOG.md has been updated
        run: |
          changeLogFileName=CHANGELOG.md
          changed_files=$(git diff --name-only --diff-filter=d origin/main..HEAD)
          if [[ ! "$changed_files" =~ "$changeLogFileName" ]]; then
            echo "Pull request does not contain changes to '$changeLogFileName'. Please update '$changeLogFileName' and try again."
            exit 1
          fi
