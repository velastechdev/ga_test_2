# HuliPractice Metrics Metabase Dashboards

The goal of this repository is to provide the necessary mechanisms in order to process the following actions in the HuliPractice metrics Metabase instances ([staging](https://hp-metrics.hulilabs.xyz/) and [production](https://hp-metrics.huli.io/)):

- [Export a serialization](#export-a-serialization)
- [Import a serialization](#import-a-serialization)
- Update the dashboard version associated to a user

[Serialization](https://www.metabase.com/docs/latest/installation-and-operation/serialization) is a **Metabase** feature which lets you create an **export** of **settings** of a Metabase that can then be **imported** into one or more **Metabases instances**. The serialization is done through exporting and importing collections; a [collection in Metabase](https://www.metabase.com/docs/latest/exploration-and-organization/collections) is basically a folder that holds:

- More folders
- Metabase dashboards
- Metabase metrics
- Metabase questions
- Metabase models 

## Introduction

[HuliPractice](https://practice.hulilabs.xyz) now has a new [metrics section](https://practice.hulilabs.xyz/es#/metrics) that displays a Metabase [dashboard](https://www.metabase.com/docs/latest/dashboards/introduction) that holds metrics related to the signed in doctor. Since the dashboard already exists, the goal of this repository is to be used to automate the process of deploying dashboards from [staging](https://hp-metrics.hulilabs.xyz/) to [production](https://hp-metrics.huli.io/) (this via workflows in [github actions](https://github.com/features/actions)).

## Commands

In order to automate the process of exporting and importing Metabase collections, a Go script has been created. The script is capable of exporting a collection from a Metabase instance and saving it as a file in the `serialization/` directory.

### Export a serialization
---

```bash
go run main.go export endpoint_url=[METABASE_API_URL] api_key=[METABASE_API_KEY] collection=[1|2|3|...|n|all]
```

#### Parameters

- `endpoint_url`: The Metabase API URL to export the collection from.
- `api_key`: The API key required to authenticate with the Metabase instance.
- `collection`: The collection to export or import. Options are `1`, `2`, `3`, ..., `n`, or `all`. The `all` option exports all collections.

> Since **Export serialization** overwrites the content of the serialization folder, so be careful when exporting an specific collection and **ONLY COMMIT** the changes of the serialization folder when we **want to update the dashboards** used in **staging and production** environments for **deploying purposes**.

### Examples

Export a serialization from staging:

```bash
go run main.go export endpoint_url=https://hp-metrics.hulilabs.xyz api_key=ABCD1234 collection=12
```

### Import a serialization
---

```bash
go run main.go import endpoint_url=[METABASE_API_URL] api_key=[METABASE_API_KEY]
```

#### Parameters

- `endpoint_url`: The Metabase API URL to import the collection to.
- `api_key`: The API key required to authenticate with the Metabase instance.

#### Examples

Import a serialization in production:

```bash
go run main.go import endpoint_url=https://hp-metrics.huli.io/ api_key=ABCD1234
```

## Workflow for deploying Metabase

1. **Deploy Staging Metabase**: This step involves deploying the serialization from the staging Metabase instance and saving it in the `serialization/` directory. The deployment process is automated through GitHub Actions workflows. This workflow is triggered when a branch with the keyword DEPLOY-VERSION in its name is pushed to the repository. The workflow configuration can be found in the `.github/workflows/deploy-metabase-staging.yml` file.

2. **Deploy Production Metabase**: This step involves deploying the serialization saved in the `serialization/` directory to the production Metabase instance. The deployment process is automated through GitHub Actions workflows. This workflow is triggered when a release of this repository is done. The workflow configuration can be found in the `.github/workflows/deploy-metabase-production.yml` file.

3. **Serialization Test**: This step involves testing the serialization process. The test process is automated through GitHub Actions workflows. This workflow is triggered when the Deploy Staging Metabase workflow is executed. The workflow configuration can be found in the `.github/workflows/test-serialization.yml` file.

### Workflow Permissions

Make sure the workflow has sufficient permissions to push changes:

1. Go to the repository settings. `Settings > Actions > General > Workflow permissions`.

2. Ensure the Read and write permissions option is selected.

3. Save the changes.

### Configure Github Action Variables

Set up the required secrets and variables for your repository:

Repository Secrets:

- `STAGING_API_KEY`: The API key required to authenticate with the staging Metabase instance.

- `PRODUCTION_API_KEY`: The API key required to authenticate with the production Metabase instance.

Repository Variables:

- `DEFAULT_COLLECTION`: Specify the default collection used in the workflow.

- `STAGING_ENDPOINT`: The Metabase API URL to export the collection from.

---

This section provides a clear and concise guide for anyone setting up the repository to use your workflow. Let me know if you’d like to adjust the wording or add further details!