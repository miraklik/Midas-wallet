// SPDX-License-Identifier: MIT 
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract Wallet is Ownable {

    error InvalidAmount(uint256 amount);
    error InvalidTokenAddress();
    error InsufficientBalance(uint256 balance, uint256 amount);
    error InvalidRecipientAddress();
    error TransferFailed();
    error FeeTransferFailed();
    error ApprovalFailed();
    error InvalidSpenderAddress(address spender);
    error NotAllowed(address owner, address spender);

    mapping(address => uint256) public balances;
    mapping(address => mapping(address => uint256)) public allowed;

    uint256 public feePercentage;
    IERC20 public token;

    event Deposit(address indexed from, uint256 amount);
    event Withdraw(address indexed to, uint256 amount, uint256 fee);
    event Transfer(address indexed from, address indexed to, uint256 amount);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    constructor(address _token, uint256 _feePercentage) Ownable(msg.sender) {
        if (_token == address(0)) {
            revert InvalidTokenAddress();
        }
        token = IERC20(_token);
        feePercentage = _feePercentage;
    }

    function deposit(uint256 _amount) public {
        if (_amount <= 0) {
            revert InvalidAmount(_amount);
        }
        if (!token.transferFrom(msg.sender, address(this), _amount)) {
            revert TransferFailed();
        }
        balances[msg.sender] += _amount;
        emit Deposit(msg.sender, _amount);
    }

    function withdraw(address _to, uint256 _amount) public {
        if (_amount <= 0) {
            revert InvalidAmount(_amount);
        }
        if (balances[msg.sender] < _amount) {
            revert InsufficientBalance(balances[msg.sender], _amount);
        }
        if (_to == address(0)) {
            revert InvalidRecipientAddress();
        }

        uint256 fee = (_amount * feePercentage) / 100;
        uint256 amountAfterFee = _amount - fee;

        balances[msg.sender] -= _amount;

        if (!token.transfer(_to, amountAfterFee)) {
            revert TransferFailed();
        }

        if (fee > 0 && !token.transfer(owner(), fee)) {
            revert FeeTransferFailed();
        }

        emit Withdraw(_to, _amount, fee);
    }

    function getBalance() public view returns (uint256) {
        return balances[msg.sender];
    }

    function transfer(address _to, uint256 _amount) public {
        if (_amount <= 0) {
            revert InvalidAmount(_amount);
        }
        if (_amount > balances[msg.sender]) {
            revert InsufficientBalance(balances[msg.sender], _amount);
        }
        if (_to == address(0)) {
            revert InvalidRecipientAddress();
        }

        balances[msg.sender] -= _amount;
        balances[_to] += _amount;

        emit Transfer(msg.sender, _to, _amount);
    }

    function approve(address _spender, uint256 _amount) public {
        if (_spender == address(0)) {
            revert InvalidSpenderAddress(_spender);
        }

        allowed[msg.sender][_spender] = _amount;
        if (!token.approve(_spender, _amount)) {
            revert ApprovalFailed();
        }

        emit Approval(msg.sender, _spender, _amount);
    }

    function allowance(address _owner, address _spender) public view returns (uint256) {
        return allowed[_owner][_spender];
    }

    function transferFrom(address _from, address _to, uint256 _amount) public {
        if (_amount <= 0) {
            revert InvalidAmount(_amount);
        }
        if (_to == address(0)) {
            revert InvalidRecipientAddress();
        }
        if (_amount > allowed[_from][msg.sender]) {
            revert NotAllowed(_from, msg.sender);
        }
        if (_amount > balances[_from]) {
            revert InsufficientBalance(balances[_from], _amount);
        }

        balances[_from] -= _amount;
        balances[_to] += _amount;
        allowed[_from][msg.sender] -= _amount;

        emit Transfer(_from, _to, _amount);
    }

    function setFeePercentage(uint256 _newFeePercentage) public onlyOwner {
        feePercentage = _newFeePercentage;
    }

    receive() external payable {
        emit Deposit(msg.sender, msg.value);
    }
}
