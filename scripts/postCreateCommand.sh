#!/bin/bash
# This script is executed after the creation of a new project.

apt-get install -y direnv

cp scripts/pre-commit.sh .git/hooks/pre-commit
cp scripts/pre-push.sh .git/hooks/pre-push
chmod 755 .git/hooks/pre-commit
chmod 755 .git/hooks/pre-push
