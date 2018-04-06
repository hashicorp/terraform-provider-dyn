if [ -z "$CI_NAME" ]; then
    echo "Skipping publishing since not running in CI"
    exit 0
fi

packages_bucket="cloudworks-tf-providers-builds"
artifacts_dir=$1

aws s3 sync $artifacts_dir s3://$packages_bucket --exclude "*" --include "*.zip"