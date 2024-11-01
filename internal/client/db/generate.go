package db

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../../bin/minimock -i github.com/spv-dev/platform_common/pkg/db.TxManager -o ./mocks/ -s "_minimock.go"
