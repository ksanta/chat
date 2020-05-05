# Chat

This is a very simple chat application. It's a little different from the usual chat application because:

- It does not identify each person chatting
- It does not keep a history of chats
- Chat messages only go out to people who are connected at the time

This is great for sharing anonymous, honest feedback with friends for example during boring corporate presentations.

## Requirements

* AWS CLI already configured with Administrator permission
* [Golang](https://golang.org)

## Setup process

### Building

```shell
sam build
```

## Packaging and deployment

To deploy your application for the first time, run the following in your shell:

```bash
sam deploy --guided
```

This assumes you own a top level domain on AWS. The deployment package will create a "chat" subdomain for you.
