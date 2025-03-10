# ace

"ace" is short for "\[Azure\] **A**pp**C**onfiguration **E**nvironment variables"

## Codespaces setup

This would apply to any generic debian/ubuntu based system, probably:

```sh
# install azure cli (might want/need to do this in a venv, or via system package manager)
pip install azure-cli
# setup dagger for ci testing
curl -fsSL https://dl.dagger.io/dagger/install.sh | BIN_DIR=$HOME/.local/bin sh
```