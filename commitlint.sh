#!/bin/sh

set -e

COMMIT_MSG=$(git log -1 --pretty=%B)

echo "🔍 Validating commit message: '$COMMIT_MSG'"

# Valid types based on Conventional Commits.
TYPES="feat|fix|docs|style|refactor|perf|test|chore|build|ci|revert"

# Regular expression for valid semantic commit.
echo "$COMMIT_MSG" | grep -qE "^($TYPES)(\([a-z0-9_-]+\))?: .+"

if [ $? -ne 0 ]; then
  echo "❌ Invalid commit message."
  echo ""
  echo "👉 Use Conventional Commits format:"
  echo "   type(scope): description"
  echo ""
  echo "✅ Examples:"
  echo "   feat(auth): add login support"
  echo "   fix(api): handle null response"
  echo "   docs(readme): update installation steps"
  exit 1
fi

echo "✅ Commit message is valid!"
