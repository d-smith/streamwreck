# streamwreck

Reconcile the records read from a stream from those written to a stream.

## Notes

Create the vpc: 

```console
aws cloudformation create-stack --stack-name myvpc --template-body file://cfn/vpc.yml
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



## TODOs

* redo API to use Api resource not function with http event