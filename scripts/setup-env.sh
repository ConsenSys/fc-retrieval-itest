#!/bin/bash
set -e
echo "*******************************************************************************"
echo "*** Set-up the env file  ***"
echo "*******************************************************************************"

cd ..

ITEST_DIR="../fc-retrieval-itest"
ITEST_BRANCH=`git rev-parse --abbrev-ref HEAD`
echo "itest repo branch: $ITEST_BRANCH"


# Check repo:
# $1 is the name of the repo. fc-retrieval-client would be client.
# $2 is the relative directory of the repo
#
check_repo() {
    # Check directory of repo exists.
    if [ ! -d "$2" ]; then
        echo "ERROR: $1 directory $2 does not exist"
        exit 1
    fi
    echo "Found $1 repo: $2"

    cd $2
    OTHER_REPO_BRANCH=`git rev-parse --abbrev-ref HEAD`
    echo "$1 repo branch: $OTHER_REPO_BRANCH"
    BRANCH_EXISTS=`git branch -r --list origin/$ITEST_BRANCH`

    if [ $ITEST_BRANCH != $OTHER_REPO_BRANCH ]; then 
        echo "itest and $1 branch do not match"
        if [ -n "$BRANCH_EXISTS" ]; then
            echo "ERROR: Branch $ITEST_BRANCH exists on the $1 repo, but the $1 repo is currently using branch $OTHER_REPO_BRANCH"
            exit 1
        fi
    fi


    cd $ITEST_DIR
    pwd
    echo ${3}_IMAGE=consensys/fc-retrieval-${1}:develop-${OTHER_REPO_BRANCH} >> .env
}

rm -f .env

check_repo gateway ../fc-retrieval-gateway GATEWAY
check_repo register ../fc-retrieval-register REGISTER
check_repo provider ../fc-retrieval-provider PROVIDER
check_repo itest ../fc-retrieval-itest ITEST

cat .env