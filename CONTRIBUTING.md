# Contributing to Sopsy

First off, thanks for taking the time to contribute! ❤️

All types of contributions are encouraged and valued. See the [Table of Contents](#table-of-contents) for different ways to help and details about how this project handles them. Please make sure to read the relevant section before making your contribution. It will make it a lot easier for us maintainers and smooth out the experience for all involved. The community looks forward to your contributions.

## Table of Contents

- [I Have a Question](#i-have-a-question)
- [I Want To Contribute](#i-want-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Your First Code Contribution](#your-first-code-contribution)

## I Have a Question

> If you want to ask a question, we assume that you have read the available [Documentation](README.md).

Before you ask a question, it is best to search for existing [Issues](https://github.com/enbiyagoral/sopsy/issues) that might help you. In case you have found a suitable issue and still need clarification, you can write your question in this issue. It is also advisable to search the internet for answers first.

## I Want To Contribute

### Reporting Bugs

- **Search Existing Issues:** Before creating a bug report, please check if the issue has already been reported.
- **Use the Template:** When creating a new issue, please use the provided Bug Report template.
- **Provide Details:** Include as much information as possible: steps to reproduce, expected behavior, screenshots, and logs.

### Suggesting Enhancements

- **Explain the Why:** Clearly describe the problem you are trying to solve.
- **Describe the Solution:** Explain how your suggested enhancement will work.

### Your First Code Contribution

1. **Fork the Repository**
2. **Clone the Fork:** `git clone https://github.com/enbiyagoral/sopsy.git`
3. **Create a Branch:** `git checkout -b feat/my-new-feature`
4. **Make Changes:** Write your code and tests.
5. **Run Tests:** Ensure all tests pass.
6. **Commit Changes:** Use conventional commits (e.g., `feat: add new command`).
7. **Push to Fork:** `git push origin feat/my-new-feature`
8. **Create Pull Request:** Go to the original repository and click "New Pull Request".

## Styleguides

### Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/). **Merge to `main` triggers [go-semantic-release](https://github.com/go-semantic-release/semantic-release)**: commit types determine the next version (feat → minor, fix → patch, BREAKING CHANGE → major), a tag is created, and [GoReleaser](https://goreleaser.com/) builds and publishes the release. All in Go, no Node.js.

- `feat:` for new features (minor bump)
- `fix:` for bug fixes (patch bump)
- `docs:` for documentation changes
- `chore:` for build tasks, package manager configs, etc.
- `style:` for formatting changes
- Footer `BREAKING CHANGE:` or `!` after type for major version bump

### Go Code Style

- Use `gofmt` to format your code.
- Run `golangci-lint` before submitting a PR.
