# Codebase standup
Standup summary based on changes to Codebase tickets

## Usage

Your updates:

    codebase-standup

Updates for everyone working on the project

    codebase-standup --all

Updates for list of users (based on first names)

    codebase-standup karol john

## How to install
This one is simple just run:

    go get github.com/mrfuxi/codebase-standup

## Configuration
Program requires config file `config.yaml` to work.

Example:
```
auth:
  username: company/your-username
  apikey: your-api-key

general:
  company: company-name
  project: project-name

mapping:
  status:
    New -> In Progress:                    Working on
```

