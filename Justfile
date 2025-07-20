PACKAGES:='stream'

test:
    #!/usr/bin/env bash
    set -euo pipefail
    for pkg in {{PACKAGES}}; do
        echo "Testing package: $pkg"
        # 假设有一个测试命令
        cd $pkg && go test ./...
    done

lint:
    #!/usr/bin/env bash
    set -euo pipefail
    for pkg in {{PACKAGES}}; do
        echo "Linting package: $pkg"
        # 假设有一个 lint 命令
        cd $pkg && golangci-lint run
    done

