# gh-branch-protection

`gh-branch-protection` is a GitHub CLI extension that lets you manage branch protection rules on your repositories directly from the command line.

## Features

- ✅ List all branch protection rules in a repository
- 🔍 Get detailed protection rule by rule ID
- ⚙️ Apply branch protection rules via a JSON file and pattern
- 🗑️ Delete branch protection rules by rule ID
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

- `list`: List all branch protection rules in the current repository.
- `get <rule-id>`: Get the branch protection rule for a given rule ID.
- `set <pattern>`: Apply branch protection rules to branches matching the pattern from stdin (expects a valid JSON structure).
- `delete <rule-id>`: Delete a branch protection rule by rule ID.

## Examples


List all branch protection rules:

```sh
gh branch-protection list
```

Get a branch protection rule by rule ID:

```sh
gh branch-protection get <rule-id>
```

Set protection rules for branches matching a pattern (e.g., `main`):

```sh
cat protection.json | gh branch-protection set main
```

Delete a branch protection rule by rule ID:

```sh
gh branch-protection delete <rule-id>
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
* Not all the protections are supported yet.

## License

This project is licensed under the [MIT License](LICENSE).
