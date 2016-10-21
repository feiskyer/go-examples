## Sign the Hyper.sh API Requests

For security, most requests to hyper.sh must be signed with the credential of a user, which consists of an access key ID and secret access key. The hyper.sh API signature algorithm are based on [AWS Signature Version 4](http://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html), and there have already been some existing implementations for some popular languages:

- Go Implementation: [hyperhq/hypercli](https://github.com/hyperhq/hypercli/blob/302a6b530148f6a777cd6b8772f706ab5e3da46b/vendor/src/github.com/docker/engine-api/client/sign4.go#L73)
- Python Implementation: [tardyp/hyper_sh](https://github.com/tardyp/hyper_sh/blob/master/hyper_sh/requests_aws4auth/aws4auth.py)
- NodeJS/Javascript Implementation: [npm: hyper-aws4](https://www.npmjs.com/package/hyper-aws4)

### Difference between Hyper.sh Signature and AWS Signature V4

The difference between Hyper.sh Signature and AWS signature V4 includes:

- Use host, region, service name of Hyper instead of AWS.
- Change the HTTP headers `X-AMZ-*` to `X-Hyper-*`
- Change the literatures with `"AWS"` to `"HYPER"`

### Detail steps of API signature generating

#### Step 0: Prepare the requests

The signed requests must includes the following headers:

- `Content-Type`, default value is `application/json`
- `X-Hyper-Date`, the API timestamps
- `Host`, the API endpoint, for example `us-west-1.hyper.sh`

If not present, the `Content-Type` will be initialized as `application/json` and the `X-Hyper-Date` will be UTC time in the format of `20060102T150405Z`.

#### Step 1: Create a Canonical Request

Hash the request body with **SHA256**, and write the hash in the Header `X-Hyper-Content-Sha256`.

Then, collect the headers to be hashed, including `Content-Type`, `Content-Md5`, `Host`, and all headers with `X-Hyper-` prefix. The headers are sorted by alphabet with the header name (lowercase) as key. Note, if the `Host` header contains a port, such as `us-west-1.hyper.sh:443`, the `:port` part will be dropped.

The `headersTobeSign` are joined with colon (`:`) and newline (`\n`), for example:

```
content-type:application/md5\nhost:us-west-1.hyper.sh\nx-hyper-content-sha256:111222333aaabbbcccddde\nx-hyper-date:20060102T150405Z\n
```

Then, join the headers with semicolon, for example

```
content-type;host;x-hyper-content-sha256;x-hyper-date
```

Then we could get the canonical request, which joins the following parts with newline(`\n`): request method, URI path, query string, the above `headersTobeSign`, the header list, and the hash of payload.

And we calculate the SHA256 checksum of Canonical Request.

#### Step 2:  Create a String to Sign

The string to sign contains 4 parts, and joined with newline(`\n`):

- Algorithm: literature `"HYPER-HMAC-SHA256"`
- Request time stamp
- Request scope, includes the following parts joined with slash(`/`)
  - Region: default is `us-west-1`
  - Service: default is `hyper`
  - Date: first 8 bytes of timestamp, e.g. the date part.
  - Literature `"hyper_request"`
- The hex string of hashed canonical request got in step 1

#### Step 3: Calculate the Signature

Use the HMAC SHA256 Algorithm to sign the request, we call it several times to get the signing key firstly:

```
kDate := hmacSHA256((keyPartsPrefix+secretKey), date)
kRegion := hmacSHA256(kDate, region)
kService := hmacSHA256(kRegion, service)
kSigning := hmacSHA256(kService, keyPartsRequest)
```

In the above code,

- `keyPartsPrefix` is `"HYPER"`,
- `secretKey` is the user's secretKey
- `region` and `service` are got in step 2, and
- `keyPartsRequest` is `"hyper_request"`

Having gotten the kSigning, we calculate the Signature of the string in step 2 with another `hmacSHA256`:

`hmacSHA256(signingKey, stringToSign)`

#### Step4: Add the Signing Information to the Request

The signature will be inserted as `Authorization` header, the content are

```
HYPER-HMAC-SHA256  Credential={AccessKey}/{Request Scope}, SignedHeaders={Signed Header}, Signature={Signature}
```

Where the

- `{AccessKey}` is the AccessKey of the user;
- `{Request Scope}` is the request scope in step 2;
- `{Signed Header}` is the semicolon joined header list in Step 1;
- `{Signature}` is the signature we got in Step 3.
