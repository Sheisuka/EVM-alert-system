pragma solidity ^0.8.30;

contract Store {
  event ItemSet(bytes32 key, uint256 value);
  mapping (bytes32 => uint256) public items;

  function setItem(bytes32 key, uint256 value) external {
    items[key] = value;
    emit ItemSet(key, value);
  }

  function retrieve(bytes32 key) public view returns (uint256) {
    return items[key];
  }
}