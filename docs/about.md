# What is Pair?

<img style="height: 130px; width: auto; margin-left: 10px;" src="/favicon.svg" alt="Pair" align="right" />

Pair is a command-line tool designed to simplify Git co-author management for collaborative development. When multiple developers work on the same code, Pair makes it easy to give proper credit to everyone by adding standardized `Co-authored-by:` trailers to commit messages.

::: info
Although named "Pair," the tool supports any number of co-authors, making it perfect for both pair programming and larger mob programming sessions.
:::

## Why use Pair?

When collaborating on code, especially during pair or mob programming sessions, it's important to maintain a clear record of who contributed to the work. Git natively supports co-authorship through commit trailers, but managing these manually can be tedious and error-prone.

Manually adding co-authors to commit messages is particularly annoying because:

* It requires remembering exact name and email formats
* You need to type the same information repeatedly
* It's easy to make typos in emails, breaking attribution
* The syntax must be precise for platforms like GitHub to recognize co-authors
* It interrupts your workflow when switching between different pairing partners

Pair streamlines this process with intuitive commands that eliminate these frustrations.

## Key features

* Maintain a roster of collaborators: Store your frequent collaborators with easy-to-remember aliases
* Simple command-line interface: Add or remove co-authors with straightforward commands
* Interactive selection: Use fuzzy-search to quickly find and select co-authors
* Visibility: See active co-authors at any time with a simple command
* Git integration: Works seamlessly with Git's commit template mechanism
* Flexible configuration: Support for both global and project-specific co-author lists

## How it works

Pair manages Git's commit templates to include "Co-authored-by:" trailers in your commits. These trailers are recognized by GitHub, GitLab, and other Git hosting platforms to properly attribute commits to multiple authors.

## Benefits

* Recognition: Ensure all contributors get proper credit for their work
* Accountability: Track who participated in which code changes
* Collaboration: Encourage more pair programming by making attribution simple
* Consistency: Maintain standardized co-author formatting across your commits
* Efficiency: Save time with quick commands instead of typing out full names and emails
* Pair is built to be lightweight, easy to use, and integrate seamlessly into your existing Git workflow, making collaborative coding more pleasant and organized.
