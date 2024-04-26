// SPDX-License-Identifier: GPL-3.0-only

pragma solidity 0.8.25;

struct GasPayment {
  address token;
  uint256 perUnit;
}

struct Order {
  address manager;
  uint256 gas;
  uint256 expire;
  uint256 start;
  uint256 end;
  uint256 txWindow;
  uint256 redeemWindow;
  GasPayment gasPrice;
  GasPayment gasGuarantee;
}

struct FilteredOrder {
  uint256 id;
  Order order;
  OrderStatus status;
  uint256 gasLeft;
  address executor;
}

enum OrderStatus {
  None,
  Pending,
  Accepted,
  Active,
  Inactive,
  Untaken,
  Closed
}

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
  function domainSeparator() external view returns (bytes32) {}

  function nonce(address, uint256) external view returns (bool) {}

  function gasOrder(uint256) external view returns (Order memory) {}

  function orderAccept(uint256) external {}

  function messageValidate(Message calldata) external view returns (Validation) {}

  function getExecutorOrders(address, bool, uint256, uint256) external view returns (FilteredOrder[] memory) {}

  function execute(Message calldata, bytes calldata) external {}
}
