# Development Guide

This document provides guidelines for setting up the development environment, running tasks, and contributing to the project.

## Prerequisites

* [VS Code](https://code.visualstudio.com/)
* [Task](https://taskfile.dev/)
* [Go](https://go.dev/)
* [Node.js](https://nodejs.org/) with [Yarn](https://yarnpkg.com/)

## Setting up Development Environment

1. Clone and open the repository

    ```bash
    git clone git@github.com:Stellatarum/docker-machine-driver-pve.git
    ```

    ```bash
    code docker-machine-driver-pve
    ```

1. Install recommended VS Code extensions: `F1 > Extensions: Show Recommended Extensions`

1. Install tools:

    ```bash
    task tools
    ```

1. Install dependencies:

    ```bash
    task dependencies
    ```

## Running tasks

We use [Task](https://taskfile.dev/) for common development tasks. Run following command to list available tasks:

```bash
task --list
```

## Commit Message Guidelines

We follow the [Conventional Commits v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) specification. Allowed commit types are defined in [`.commitlintrc.yaml`](../.commitlintrc.yaml).

You can lint commits in relation to the `main` branch using `task lint:commits`.
