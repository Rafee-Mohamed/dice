---
title: HANDSHAKE
description: HANDSHAKE tells the server the purpose of the connection
---

<!-- This file is automatically generated. Any modifications made directly to this file
  may be overwritten. For more details on how this file is generated and how to use
  the related commands, refer to the documentation available in the `internal/cmd/cmd_*.go` files.
-->

#### Syntax

```
HANDSHAKE client_id execution_mode
```

HANDSHAKE is used to tell the DiceDB server the purpose of the connection. It
registers the client_id and execution_mode.

The client_id is a unique identifier for the client. It can be any string, typically
a UUID.

The execution_mode is the mode of the connection, it can be one of the following:

1. "command" - The client will send commands to the server and receive responses.
2. "watch" - The connection in the watch mode will be used to receive the responses of query subscriptions.

If you use DiceDB SDK or CLI then this HANDSHAKE command is automatically sent when the connection is established
or when you establish a subscription.

#### Examples

```

localhost:7379> HANDSHAKE 4c9d0411-6b28-4ee5-b78a-e7e258afa52f command
OK OK

```
