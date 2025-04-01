# contracttoken

This module provides a function to issue a token from CosmWasm contract.

## Denom

The denom of the token issued from this module is:

```plaintext
contract/[contract_address]
```

## BeforeSendHook

If the contract token enables the before send hook, the `BeforeSend` message will be executed before the token is sent.

```rust
enum ExecuteMsg {
    BeforeSend(BeforeSendMsg),
}

struct BeforeSendMsg {
    from: String,
    to: String,
    amount: Uint128,
}
```

- The caller of the contract is the contract address `env.contract.address`.
  - This can be used to verify the authority of the caller.
- The `from` is the sender of the token.
- The `to` is the recipient of the token.
- The `amount` is the amount of the token to be sent.
- If the contract returns an error, the transaction will fail.
- The payment to the contract in `BeforeSend` is always empty.

Also see [bank/README.md](../bank/README.md)
