# gh-branch-protection

`gh-branch-protection` is a GitHub CLI extension that lets you manage branch protection rules on your repositories directly from the command line.

## Features

- ✅ List protected branches in a repository  
- 🔍 Get detailed protection rules for a specific branch  
- ⚙️ Apply branch protection rules via a JSON file  
- 🗑️ Delete branch protection rules from a branch  
- 📦 Works seamlessly with GitHub CLI (`gh`)  
- 💻 Simple installation and usage  

## Build

To build the extension from source:

```sh
go build
```

## Installation

To install the extension locally:

```sh
gh extensions install .
```

> ℹ️ Requires [GitHub CLI](https://cli.github.com/) and [Go](https://golang.org/) installed.

## Usage

After installation, use the extension with the `gh branch-protection` command:

```sh
gh branch-protection [command] [args]
```

### Commands

- `list`: List all protected branches in the current repository.
- `get <branch>`: Get the branch protection rules for a given branch.
- `set <branch>`: Apply branch protection rules from stdin (expects a valid JSON structure).
- `delete <branch>`: Delete branch protection rules from a given branch.

## Examples

List all protected branches:

```sh
gh branch-protection list
```

Get protection rules for the `master` branch:

```sh
gh branch-protection get master
```

Set protection rules for the `master` branch from a file:

```sh
cat protection.json | gh branch-protection set master
```

Delete protection rules from the `master` branch:

```sh
gh branch-protection delete master
```

Example `protection.json`:

```json
{
  "AllowsDeletions": false,
  "AllowsForcePushes": false,
  "RequiredApprovingReviewCount": 1,
  "RequiresApprovingReviews": true,
  "RequiresCodeOwnerReviews": false,
  "RequiresCommitSignatures": true,
  "RequiresLinearHistory": true,
  "RequiresConversationResolution": true,
  "IsAdminEnforced": true,
  "RestrictsPushes": false,
  "RestrictsReviewDismissals": false
}
```

## Bugs and limitations
* The `set` command does not overwrite existing protection rules. To update a rule, first delete the current protection, then use the `set` command to apply the new rules.
* The argument for the `set`, `get`, and `delete` commands is a pattern that can match zero, one, or multiple branches. Currently, these commands expect an existing branch name; if the branch does not exist, the commands will not work. Future improvements should allow these commands to operate regardless of whether the branch exists.
* Not all the protections are supported yet.

## License

This project is licensed under the [MIT License](LICENSE).
