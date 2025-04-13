#!/bin/bash

# Base directory where all repositories are located
BASE_DIR="$GOPATH/src"
FASTXML_DIR="$BASE_DIR/github.com/PubMatic-OpenWrap/fastxml"
VASTUNWRAP_DIR="$BASE_DIR/git.pubmatic.com/PubMatic/vastunwrap"
PREBID_SERVER_DIR="$BASE_DIR/github.com/PubMatic-OpenWrap/prebid-server"
HEADER_BIDDING_DIR="$BASE_DIR/header-bidding"

# Function to get latest commit from repository
get_latest_commit() {
    local repo_dir=$1
    local current_dir=$(pwd)
    
    cd "$repo_dir" || exit 1
    local commit_id=$(git rev-parse HEAD)
    cd "$current_dir" || exit 1
    
    echo "$commit_id"
}

# Function to update go.mod and commit changes
update_repo() {
    local repo_dir=$1
    local repo_name=$2
    local dep_repo=$3
    local commit_id=$4
    local branch_name=$5
    
    echo "Updating $repo_name with $dep_repo@$commit_id..."
    cd "$repo_dir" || exit 1
    
    # If branch name is not provided, use current branch
    if [ -z "$branch_name" ]; then
        branch_name=$(git branch --show-current)
    else
        # Check if branch exists
        if git show-ref --verify --quiet "refs/heads/$branch_name"; then
            # Branch exists, switch to it
            git checkout "$branch_name"
        else
            # Branch doesn't exist, create and switch to it
            git checkout -b "$branch_name"
        fi
    fi
    
    # Update dependencies and tidy
    if [ "$dep_repo" = "git.pubmatic.com/PubMatic/vastunwrap" ]; then
        # For vastunwrap, we need to update the replace directive
        sed -i '' "s|=> git.pubmatic.com/PubMatic/vastunwrap v.*|=> git.pubmatic.com/PubMatic/vastunwrap $commit_id|" go.mod
        go mod tidy
    elif [ "$dep_repo" = "github.com/PubMatic-OpenWrap/prebid-server/v3" ]; then
        # For prebid-server, update the version in go.mod
        if grep -q "github.com/PubMatic-OpenWrap/prebid-server/v3" go.mod; then
            # If module already exists, update its version
            sed -i '' "s|github.com/PubMatic-OpenWrap/prebid-server/v3 v.*|github.com/PubMatic-OpenWrap/prebid-server/v3 $commit_id|" go.mod
        else
            # If module doesn't exist, add it
            echo "require github.com/PubMatic-OpenWrap/prebid-server/v3 $commit_id" >> go.mod
        fi
        go mod tidy
    else
        go get -u "$dep_repo@$commit_id"
        go mod tidy
    fi
    
    # Commit changes if there are any
    if git diff --quiet go.mod go.sum; then
        echo "No changes in $repo_name"
        git push origin "$branch_name"
        return 1
    else
        git add go.mod go.sum
        git commit -m "chore: update $dep_repo dependency to $commit_id"
        git push origin "$branch_name"
        echo "Changes committed in $repo_name on branch $branch_name"
        return 0
    fi
}

# Function to update fastxml dependencies
update_fastxml_deps() {
    local commit_id=$1
    local branch_name=$2
    local vastunwrap_updated=false
    local prebid_updated=false
    
    echo "Phase 1: Updating fastxml dependencies in all repositories..."
    # Update dependent repositories in order
    update_repo "$VASTUNWRAP_DIR" "vastunwrap" "github.com/PubMatic-OpenWrap/fastxml" "$commit_id" "$branch_name"
    if [ $? -eq 0 ]; then
        vastunwrap_updated=true
    fi
    
    update_repo "$PREBID_SERVER_DIR" "prebid-server" "github.com/PubMatic-OpenWrap/fastxml" "$commit_id" "$branch_name"
    if [ $? -eq 0 ]; then
        prebid_updated=true
    fi
    
    update_repo "$HEADER_BIDDING_DIR" "header-bidding" "github.com/PubMatic-OpenWrap/fastxml" "$commit_id" "$branch_name"
    
    # Phase 2: If vastunwrap was updated, update its dependencies
    if [ "$vastunwrap_updated" = true ]; then
        echo "Phase 2: Updating vastunwrap dependencies..."
        local vastunwrap_commit=$(get_latest_commit "$VASTUNWRAP_DIR")
        update_repo "$PREBID_SERVER_DIR" "prebid-server" "git.pubmatic.com/PubMatic/vastunwrap" "$vastunwrap_commit" "$branch_name"
        if [ $? -eq 0 ]; then
            prebid_updated=true
        fi
        update_repo "$HEADER_BIDDING_DIR" "header-bidding" "git.pubmatic.com/PubMatic/vastunwrap" "$vastunwrap_commit" "$branch_name"
    fi
    
    # Phase 3: If prebid-server was updated in any phase, update its dependencies
    if [ "$prebid_updated" = true ]; then
        echo "Phase 3: Updating prebid-server dependencies..."
        local prebid_commit=$(get_latest_commit "$PREBID_SERVER_DIR")
        update_repo "$HEADER_BIDDING_DIR" "header-bidding" "github.com/PubMatic-OpenWrap/prebid-server/v3" "$prebid_commit" "$branch_name"
    fi
}

# Function to update vastunwrap dependencies
update_vastunwrap_deps() {
    local commit_id=$1
    local branch_name=$2
    local prebid_updated=false
    
    echo "Phase 1: Updating vastunwrap dependencies..."
    # Update dependent repositories in order
    update_repo "$PREBID_SERVER_DIR" "prebid-server" "git.pubmatic.com/PubMatic/vastunwrap" "$commit_id" "$branch_name"
    if [ $? -eq 0 ]; then
        prebid_updated=true
    fi
    
    update_repo "$HEADER_BIDDING_DIR" "header-bidding" "git.pubmatic.com/PubMatic/vastunwrap" "$commit_id" "$branch_name"
    
    # Phase 2: If prebid-server was updated, update its dependencies
    if [ "$prebid_updated" = true ]; then
        echo "Phase 2: Updating prebid-server dependencies..."
        local prebid_commit=$(get_latest_commit "$PREBID_SERVER_DIR")
        update_repo "$HEADER_BIDDING_DIR" "header-bidding" "github.com/PubMatic-OpenWrap/prebid-server/v3" "$prebid_commit" "$branch_name"
    fi
}

# Function to update prebid-server dependencies
update_prebidserver_deps() {
    local commit_id=$1
    local branch_name=$2
    
    # Update header-bidding with latest prebid-server
    update_repo "$HEADER_BIDDING_DIR" "header-bidding" "github.com/PubMatic-OpenWrap/prebid-server/v3" "$commit_id" "$branch_name"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 -r <repository> [-c <commit-id>] [-b <branch-name>]"
    echo "Options:"
    echo "  -r  Repository to update (required). Options: fastxml, vastunwrap, prebidserver"
    echo "  -c  Commit ID to update to (optional). If not provided, uses latest commit"
    echo "  -b  Branch name to use (optional). If not provided, uses current branch"
    echo ""
    echo "Examples:"
    echo "  $0 -r fastxml -c abc123def456 -b feature/update-deps"
    echo "  $0 -r vastunwrap -b develop"
    echo "  $0 -r prebidserver"
    exit 1
}

# Parse command line arguments
while getopts "r:c:b:h" opt; do
    case $opt in
        r)
            REPO_TO_UPDATE="$OPTARG"
            ;;
        c)
            COMMIT_ID="$OPTARG"
            ;;
        b)
            BRANCH_NAME="$OPTARG"
            ;;
        h)
            show_usage
            ;;
        \?)
            echo "Invalid option: -$OPTARG"
            show_usage
            ;;
        :)
            echo "Option -$OPTARG requires an argument"
            show_usage
            ;;
    esac
done

# Validate repository argument
if [ -z "$REPO_TO_UPDATE" ]; then
    echo "Error: Repository (-r) is required"
    show_usage
fi

# If commit ID is not provided, get the latest commit from the repository
if [ -z "$COMMIT_ID" ]; then
    case "$REPO_TO_UPDATE" in
        "fastxml")
            COMMIT_ID=$(get_latest_commit "$FASTXML_DIR")
            echo "Using latest fastxml commit: $COMMIT_ID"
            ;;
        "vastunwrap")
            COMMIT_ID=$(get_latest_commit "$VASTUNWRAP_DIR")
            echo "Using latest vastunwrap commit: $COMMIT_ID"
            ;;
        "prebidserver")
            COMMIT_ID=$(get_latest_commit "$PREBID_SERVER_DIR")
            echo "Using latest prebid-server commit: $COMMIT_ID"
            ;;
        *)
            echo "Unsupported repository: $REPO_TO_UPDATE"
            echo "Supported repositories: fastxml, vastunwrap, prebidserver"
            exit 1
            ;;
    esac
fi

case "$REPO_TO_UPDATE" in
    "fastxml")
        update_fastxml_deps "$COMMIT_ID" "$BRANCH_NAME"
        ;;
    "vastunwrap")
        update_vastunwrap_deps "$COMMIT_ID" "$BRANCH_NAME"
        ;;
    "prebidserver")
        update_prebidserver_deps "$COMMIT_ID" "$BRANCH_NAME"
        ;;
    *)
        echo "Unsupported repository: $REPO_TO_UPDATE"
        echo "Supported repositories: fastxml, vastunwrap, prebidserver"
        exit 1
        ;;
esac

echo "Update process completed!"
echo "Please review the changes in each repository and push the branches if everything looks good."
