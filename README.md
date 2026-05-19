# loom

A Git-inspired CLI utility for interacting with Large Language Models (LLMs) via the [OpenRouter API](https://openrouter.ai/). 

`loom` brings version control concepts to prompt engineering and AI chat. Instead of linear chat histories, `loom` structures conversations into isolated **workspaces** containing branching **dialogue trees**. Messages act as commits - allowing you to branch off at any point, experiment with different prompts or models, and maintain complete control over your context.

Currently, `loom` supports **text-to-text modality**, with the unique ability to swap or update the active LLM after any message in the tree.

---

## Key Features

- **Branching Dialogue Trees**: Experience conversations like code. Fork your chat at any message to explore alternative AI responses or prompt variations.
- **Isolated Workspaces**: Organize your projects, tasks, or personas into separate workspaces.
- **Model Agnostic**: Swap the active model seamlessly mid-conversation.
- **Streaming Responses**: Token-by-token output directly in your terminal for real-time interaction.
- **Interactive Logs**: Easily navigate the history of your current branch.

---

## Installation

<<<<<<< HEAD
Refer to [INSTALLATION.md](./INSTALLATION.md) for installation instructions.
=======
Refer to [INSTALLATION.md] for installation instructions.
>>>>>>> d040914047358431b71380238d09c99a477e4904

---

## Quick start

### 1. Configure your API key

```bash
loom config --key YOUR_OPENROUTER_API_KEY
```

### 2. Initialize a workspace

```bash
loom init research-project
```

### 3. Add and set a model

```bash
loom models add google/gemma-4-31b-it:free
loom models set google/gemma-4-31b-it:free
```

### 4. Send your first message

```bash
loom send "How do planes fly?"
```

---

## Command reference

| Command                   | Git Analog     | Description                                                             |
| ------------------------- | -------------- | ----------------------------------------------------------------------- |
| `loom config --key <key>` | `git config`   | Sets the OpenRouter API key for requests.                               |
| `loom init <name>`        | `git init`     | Creates a new workspace.                                                |
| `loom ws <name>`          | None           | Switches to an existing workspace.                                      |
| `loom workspace`          | None           | Lists all existing workspaces.                                          |
| `loom status`             | `git status`   | Displays the active workspace and branch.                               |
| `loom send <message>`     | `git commit`   | Sends a message to the active model (streams response to console).      |
| `loom checkout <name>`    | `git checkout` | Switches branches. Use -b flag to create and switch to a new branch.    |
| `loom branch`             | `git branch`   | Lists all branches inside the active workspace.                         |
| `loom log`                | `git log`      | Opens an interactive window with the current branch's dialogue history. |
| `loom models list`        | None           | Lists locally available/configured models.                              |
| `loom models add <id>`    | None           | Automatically adds a model by its official OpenRouter ID.               |
| `loom models set <id>`    | None           | Changes the active model for subsequent messages.                       |

---

## Workflow example: branching a conversation

```bash
# 1. Ask the initial question
loom send "Write a Python script to parse a CSV file."

# 2. Check the status and see you are on the default branch (e.g., 'main')
loom status

# 3. Create a new branch to explore an alternative approach (e.g., using pandas)
loom checkout -b use-pandas

# 4. Change the model to a more powerful one for this specific branch
loom models set anthropic/claude-3-opus

# 5. Continue the conversation on this branch
loom send "Now optimize it using the pandas library."

# 6. Switch back to the main branch anytime to follow the standard library path
loom checkout main
```
