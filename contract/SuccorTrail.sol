// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SuccorTrail {
    struct Donation {
        address donor;
        uint256 amount;
        uint256 timestamp;
    }

    struct Feedback {
        address user;
        string message;
        uint256 timestamp;
    }

    struct Receiver {
        address receiverAddress;
        string name;
        bool isRegistered;
    }

    address public owner;
    mapping(address => Receiver) public receiverMap;
    Donation[] public donations;
    Feedback[] public feedbacks;

    event DonationReceived(address indexed donor, uint256 amount);
    event FeedbackReceived(address indexed user, string message);
    event MealDistributed(address indexed receiver, string mealDescription);
    event ReceiverRegistered(address indexed receiverAddress, string name);

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    modifier onlyRegisteredReceiver() {
        require(
            receiverMap[msg.sender].isRegistered,
            "Only registered receivers can call this function"
        );
        _;
    }

    constructor() {
        owner = msg.sender;
    }

    // Register a receiver
    function registerReceiver(
        address _receiverAddress,
        string memory _name
    ) public onlyOwner {
        Receiver storage r = receiverMap[_receiverAddress];
        r.receiverAddress = _receiverAddress;
        r.name = _name;
        r.isRegistered = true;
        emit ReceiverRegistered(_receiverAddress, _name);
    }

    // Donate money
    function donate() public payable {
        require(
            msg.value > 0 && receiverMap[msg.sender].isRegistered,
            "Donation amount must be greater than zero and sender must be a registered receiver"
        );

        donations.push(Donation(msg.sender, msg.value, block.timestamp));
        emit DonationReceived(msg.sender, msg.value);
    }

    // Give feedback
    function giveFeedback(string memory _message) public {
        require(bytes(_message).length > 0, "Message cannot be empty.");
        feedbacks.push(Feedback(msg.sender, _message, block.timestamp));
        emit FeedbackReceived(msg.sender, _message);
    }

    // Distribute a meal
    function distributeMeal(
        string memory _mealDescription
    ) public onlyRegisteredReceiver {
        require(
            bytes(_mealDescription).length > 0,
            "Meal description cannot be empty."
        );

        emit MealDistributed(msg.sender, _mealDescription);
    }

    // Remove donations
    function removeDonation(uint256 index) public onlyOwner {
        require(index < donations.length, "Index out of range");
        delete donations[index];
    }

    // Get the count of donations
    function getDonationCount() public view returns (uint256) {
        return donations.length;
    }

    // Get the count of feedbacks
    function getFeedbackCount() public view returns (uint256) {
        return feedbacks.length;
    }
}
