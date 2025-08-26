#!/bin/bash
# Commit message validation example hook script

echo "📝 Validating commit message..."

commit_msg_file="$1"
commit_msg=$(cat "$commit_msg_file")

# Skip validation for merge commits
if echo "$commit_msg" | grep -q "^Merge"; then
    echo "✅ Merge commit - skipping validation"
    exit 0
fi

# Check for conventional commits format (optional)
if echo "$commit_msg" | grep -qE "^(feat|fix|docs|style|refactor|test|chore|perf|ci|build|revert)(\(.+\))?: .+"; then
    echo "✅ Conventional commit format detected"
elif echo "$commit_msg" | grep -qE "^(Add|Update|Fix|Remove|Refactor|Improve|Clean|Bump): .+"; then
    echo "✅ Standard commit format detected"
else
    echo "❌ Commit message format validation failed"
    echo ""
    echo "Expected formats:"
    echo "  Conventional: type(scope): description"
    echo "    Example: feat(auth): add login functionality"
    echo "    Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build, revert"
    echo ""
    echo "  Standard: Action: description"
    echo "    Example: Add: user authentication system"
    echo "    Actions: Add, Update, Fix, Remove, Refactor, Improve, Clean, Bump"
    echo ""
    echo "Your commit message:"
    echo "  '$commit_msg'"
    exit 1
fi

# Check minimum length
if [ ${#commit_msg} -lt 10 ]; then
    echo "❌ Commit message too short (minimum 10 characters)"
    echo "Your message: '$commit_msg' (${#commit_msg} characters)"
    exit 1
fi

# Check maximum line length for first line
first_line=$(echo "$commit_msg" | head -n1)
if [ ${#first_line} -gt 72 ]; then
    echo "❌ First line of commit message too long (maximum 72 characters)"
    echo "Your first line: '$first_line' (${#first_line} characters)"
    exit 1
fi

echo "✅ Commit message validation passed!"
exit 0