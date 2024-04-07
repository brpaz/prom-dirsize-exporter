
# Contributing

We welcome and appreciate contributions from the community. Whether you're fixing a bug, proposing a new feature, or improving the documentation, we're grateful for your efforts.

## Contributing with Issues

By submitting an issue, you can help to improve the quality and usability of the project.

### Reporting a Bug

To report a bug, please create a new issue and include the following information:

- A clear and concise description of the problem
- Steps to reproduce the issue
- The expected behavior
- The actual behavior
- Screenshots or code samples if possible
- The version of the project you are using
- The operating system and version you are using

Please be as detailed as possible in your bug report. This will help us to quickly identify and fix the issue.

### Requesting a Feature

To request a new feature, please create a new issue and include the following information:

- A clear and concise description of the feature you would like to see added
- The use case for the feature (e.g. what problem it solves)
- If applicable, any relevant data or research that supports the need for the feature
- A sample implementation or code snippet if possible

Please be as detailed as possible in your feature request. This will help us to understand the need and potential impact of the feature.

We encourage you to submit issues and provide feedback. Your contributions will help to make the project better for everyone.

For bigger changes, it´s recommended to create an issue in the repository to discuss, before starting any implementation.

## Getting Started

- Fork the repository.
- Clone the forked repository to your local machine.
- Create a new branch for your changes.
- Make your changes
- Write tests for your changes and make sure the tests and linting pass.
- Submit a Pull Request

## Contributing with code

We welcome contributions in the form of pull requests. By contributing code, you can help to add new features, fix bugs, and improve the project.

Before you start working on a new feature or bug fix, please check the existing issues to see if there is already an open issue for the problem you want to solve. If there is, please leave a comment on the issue to let us know you are working on a solution.

If there is no existing issue, please create a new one before starting your work. This will help to ensure that the changes you make are in line with the project's goals and roadmap.

### Pre-Requisites

This project requires Go 1.19 installed on your machine. Alternatively you a Docker container is provided. It´s also recommended to install the following tools, required for developing the project.

- [Task](https://taskfile.dev/) - Task is a task runner / build tool that aims to be simpler and easier to use than, for example, GNU Make.
- [lefthook](https://github.com/evilmartians/lefthook) - Fast and powerful Git hooks manager for any type of projects.
- [golangci-lint](https://golangci-lint.run/usage/install/) - Code Linting tool
- [gotestyourself/gotestsum](https://github.com/gotestyourself/gotestsum) - Test runner

### Clone the repository

```
git clone github.com/{{git_repo}}
```

### Run linting and tests

```sh
task lint
```

```sh
task test
```

## Submitting a Pull Request

We encourage the use of [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) when submitting pull requests. Semantic commits make it easier to understand the purpose of a change and improve the overall Git history.

Here's an example of a semantic commit message:

```
feat: add new feature x
```

The message consists of a type (feat), a colon, and a description of the change. Here are some common types you can use:

- **feat**: for a new feature
- **fix**: for a bug fix
- **docs**: for changes to documentation
- **style**: for code style changes (white-space, formatting, missing semi-colons, etc)
- **refactor**: for refactoring code
- **test**: for changes to tests
- **chore**: for build process, etc; no production code change
- **ci**: for CI related changes.

### Add the appropriate Pull request labels.

We use [Release Drafter](https://github.com/integration/release-drafter) to create releases for our project. Release Drafter automatically generates a draft of the release notes based on the Git history and pull request labels.

To ensure that Release Drafter can accurately generate the release notes, it is important to add pull request labels to your pull requests. The following labels are supported by Release Drafter:

* `feature`: for a new feature
* `enhancement`: for improvements to existing features
* `bug`: for a bug fix
* `documentation`: for changes to documentation
* `chore`: for build process, etc; no production code change
* `maintenance`: for internal code clean-up and refactoring
* `dependencies`: for updates to dependencies

Please add the relevant label to your pull request before submitting it. This will help to provide a clear and organized summary of changes in each release.

### Release process

Once a pull request has been reviewed, approved and merged, a draft release will be created by  [Release Drafter](https://github.com/integration/release-drafter).

After the draft release has been published, a CI job (release) will be triggered, which will use [GO Releaser](https://goreleaser.com) to create the actual release artifacts. GO Releaser will build and package the code and update the existing GitHub release with all the generated artifacts.

By using GO Releaser, the release process is streamlined and less prone to errors. This helps to ensure that new releases are created quickly and efficiently, and that the code is properly packaged and distributed.
