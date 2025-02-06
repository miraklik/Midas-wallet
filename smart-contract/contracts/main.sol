//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleContract {
    address public MyAddress;

    function setAddr() public {
        MyAddress = msg.sender;
    }

    function getAddr() public view  returns (address) {
        return MyAddress;
    }
}