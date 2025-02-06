// SPDX-License-Identifier: MIT 
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract Wallet is Ownable {
    mapping(address => uint256) public balances;
    mapping(address => mapping(address => uint256)) public allowed;

    uint256 public feePercentage;
    IERC20 public token;

    event Deposit(address indexed from, uint256 amount);
    event Withdraw(address indexed to, uint256 amount, uint256 fee);
    event Transfer(address indexed from, address indexed to, uint256 amount);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    constructor(address _token, uint256 _feePercentage) Ownable(msg.sender) {
        require(_token != address(0), "Invalid token address");
        token = IERC20(_token);
        feePercentage = _feePercentage;
    }

    function deposit(uint256 _amount) public {
        require(_amount > 0, "Invalid amount");
        require(token.transferFrom(msg.sender, address(this), _amount), "Transfer failed");
        balances[msg.sender] += _amount;
        emit Deposit(msg.sender, _amount);
    }

    function withdraw(address _to, uint256 _amount) public {
        require(_amount > 0, "Invalid amount");
        require(_amount <= balances[msg.sender], "Insufficient balance");
        require(_to != address(0), "Invalid recipient address");

        uint256 fee = (_amount * feePercentage) / 100;
        uint256 amountAfterFee = _amount - fee;

        balances[msg.sender] -= _amount;

        require(token.transfer(_to, amountAfterFee), "Transfer failed");

        if (fee > 0) {
            require(token.transfer(owner(), fee), "Fee transfer failed");
        }

        emit Withdraw(_to, _amount, fee);
    }

    function getBalance() public view returns (uint256) {
        return balances[msg.sender];
    }

    function transfer(address _to, uint256 _amount) public {
        require(_amount > 0, "Invalid amount");
        require(_amount <= balances[msg.sender], "Insufficient balance");
        require(_to != address(0), "Invalid recipient address");

        balances[msg.sender] -= _amount;
        balances[_to] += _amount;

        emit Transfer(msg.sender, _to, _amount);
    }

    function approve(address _spender, uint256 _amount) public {
        require(_spender != address(0), "Invalid spender address");

        allowed[msg.sender][_spender] = _amount;
        require(token.approve(_spender, _amount), "Approval failed");

        emit Approval(msg.sender, _spender, _amount);
    }

    function allowance(address _owner, address _spender) public view returns (uint256) {
        return allowed[_owner][_spender];
    }

    function transferFrom(address _from, address _to, uint256 _amount) public {
        require(_amount > 0, "Invalid amount");
        require(_to != address(0), "Invalid recipient address");
        require(_amount <= allowed[_from][msg.sender], "Not allowed to spend this amount");
        require(_amount <= balances[_from], "Insufficient balance");

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
