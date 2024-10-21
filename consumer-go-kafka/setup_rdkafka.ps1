# Install the librdkafka.redist package version 2.6.0
nuget install librdkafka.redist -Version 2.6.0

Copy-Item -Recurse $pwd\librdkafka.redist.2.6.0\runtimes\win-x64\native\* $env:TMP

$goPath = & go env GOPATH
Remove-Item -Force "$goPath/pkg/mod/github.com/confluentinc/confluent-kafka-go/v2@v2.6.0/kafka/librdkafka_vendor/librdkafka_windows.a"

# Copy the librdkafka.dll to the specified directory
Copy-Item .\librdkafka.redist.2.6.0\build\native\lib\win\x64\win-x64-Release\v142\librdkafka.lib "$goPath/pkg/mod/github.com/confluentinc/confluent-kafka-go/v2@v2.6.0/kafka/librdkafka_vendor/librdkafka_windows.a"