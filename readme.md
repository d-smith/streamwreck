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
