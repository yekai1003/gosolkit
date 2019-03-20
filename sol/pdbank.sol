pragma solidity^0.5.0;

contract pdbank {
    address public  owner;
    mapping(address=>uint256) public balances;//纪录每个账户余额
    uint256 public totalAmount;
    string public bankName;//银行名称
    //构造函数
    constructor(string memory _bankName) public  {
        owner = msg.sender;
        bankName = _bankName;
    }
    //充值
    function deposit() public payable {
        //do nothing
        totalAmount += msg.value;
        balances[msg.sender] += msg.value;
    }
    //提现
    function withdraw(uint256 _amount) public payable {
        if(balances[msg.sender] > _amount) {
            balances[msg.sender]  -=  _amount;  
            msg.sender.transfer(_amount);
            totalAmount -= _amount;
        }
    }
}
