# Installation
## Requirements
Python `>=3.10` is required

## Installing directly wih `pipx`

`pipx` will create an isolated environment for **loom** and add it to system `PATH`. This will ensure that project's dependencies aren't going to appear in your system.

### 1. Install `pipx`
If you don't have `pipx` installed in your system, run the following commands to fix that:

```bash
pip install pipx
pipx ensurepath
```

### 2. Install **loom** directly
```bash
pipx install git+https://github.com/ekosachev/loom.git
```

## Installing from source

### 1. Clone the repository

```bash
git clone https://github.com/ekosachev/loom.git
cd loom
```

### 2. Create and activate a virtual environment

#### Linux/MacOS

```bash
python3 -m venv venv
source venv/bin/activate
```

#### Windows

```bash
venv\Scripts\activate
```

### 3. Install **loom** from source
```bash
pip install .
```

## After installation
Ensure that **loom** is installed correctly by running
```bash
loom --help
```
