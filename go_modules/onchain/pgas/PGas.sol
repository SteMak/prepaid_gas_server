// SPDX-License-Identifier: GPL-3.0-only

pragma solidity 0.8.25;

struct Message {
  address from;
  uint256 nonce;
  uint256 order;
  uint256 start;
  address to;
  uint256 gas;
  bytes data;
}

enum Validation {
  None,
  StartInFuture,
  NonceExhaustion,
  BalanceCompliance,
  OwnerCompliance,
  TimelineCompliance
}

// solc --abi --bin PGas.sol -o build
// abigen --abi build/PGas.abi --type PGas --pkg pgas --out pgas.go
// rm -rf build

contract PGas {
  function messageValidate(Message calldata message) external view returns (Validation) {}
}
