package kafka

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../../bin/minimock -i Producer -o ./mocks/ -s "_minimock.go"
