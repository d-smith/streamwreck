# streamwreck

Reconcile the records read from a stream from those written to a stream.

# Overview

The following diagram shows the solution components.

![](./Components.png)

In this pattern, a producer application writes events to a Kinesis stream, which are processed by a lambda function. The lambda function publishes an event indicating the record was processed. Another lambda which is subscribed to the 'event processed' topic writes an entry to ElastiCache for each stream record processed, using the sequence number to identity the record. Another lambda implements a simple API that takes an array of sequence numbers and indicates in its response any that have not been recorded as being processed in ElastiCache.

To reconsole the stream events with the processed events, a client can read the records from the stream, call the reconcilliation API, and determine which have events have not been processed.

## Notes

Create the vpc: 

```console
aws cloudformation create-stack --stack-name myvpc --template-body file://cfn/vpc.yml
```

Deploy the memcached stack:

```console
aws cloudformation create-stack --stack-name mycache --template-body file://cfn/cache.yml
```

Build and deploy the lambdas:

```console
make
make deploy
```

Write records to yonder stream:

```console
aws kinesis put-record --stream-name WreakStream-Dev --data Data1234Foobar --partition-key foo
```

Get logs

```console
sam logs -n StreamProcessor --stack-name streamwreck
```

List APIs

```console
aws apigateway get-rest-apis
{
    "items": [
        {
            "id": "0gs925saae",
            "name": "streamwreck",
            "createdDate": 1546272932,
            "version": "1.0",
            "apiKeySource": "HEADER",
            "endpointConfiguration": {
                "types": [
                    "EDGE"
                ]
            }
        }
    ]
}
```

List stages for the API

```console
aws apigateway get-stages --rest-api-id 0gs925saae
```

With the rest id and stage, you can form the endpoint as https://{api id}.execute-api.{region}.amazonaws.com/{stage}/reconcile

For example:

```console
curl https://0gs925saae.execute-api.us-east-1.amazonaws.com/Stage/reconcile -d '["49591585917767627211020937568791696204996217827017883650", "49591585917767627211020937568796531908274676481155661826", "49591585917767627211020937568800158685733520506118733826"]'
```

Run the stream checker to check all records from the TRIM_HORIZON to see if they've been processed:

```console
go run streamchecker/main.go --stream WreakStream-Dev --endpoint https://hyqj92kmud.execute-api.us-east-1.amazonaws.com/Stage/reconcile
```

## TODOs

* redo API to use Api resource not function with http event