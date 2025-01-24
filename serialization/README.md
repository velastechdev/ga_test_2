# Serialization

The goal of the `serialization/` folder is to store the serialized settings from a Metabase instance (settings generated using the export command). These collections are typically compressed files (e.g., .tar.gz) that are decompressed and extracted into this directory. These collections (yaml files) can be used to be imported in another Metabase instance.

> Do not modify files manually: Changes to files in this folder might cause inconsistencies with automated processes.

## Structure
	•	Each collection is extracted into its own subfolder within this directory.
	•	The naming of the subfolders corresponds to the collection’s ID or other relevant identifiers.

## Purpose

This folder serves as a centralized location for managing downloaded data collections, ensuring they are available for further processing, analysis, or integration with the application.

## Workflow

	1.	Downloading: Collections are downloaded using the Export command on [/commands/export.go].
	2.	Decompression: Compressed files (e.g., .tar.gz) are decompressed, and their contents are extracted into this directory.
	3.	Accessing: Extracted files can be accessed or processed directly from their respective subfolders.
