# SuccorTrail

A blockchain-based solution to transform food aid distribution. SuccorTrail enhances transparency, efficiency, and accountability—ensuring that aid reaches those who need it most while combating corruption and waste.

## Key Features
- **Transparency:** Immutable, verifiable records of every transaction.
- **Combating Corruption:** Reduced opportunities for fraud by making all distributions traceable.
- **Minimizing Waste:** Optimized supply chains to ensure efficient delivery of aid.
- **Secure Distribution:** Robust security measures safeguard the distribution process.
- **Operational Efficiency:** Streamlined processes reduce delays and maximize aid impact.
- **Food Security with Integrity:** Ensures that aid reaches its intended recipients with full accountability.

## Technology Stack
- **Blockchain & Smart Contracts:**  
  Secure transactions on scalable networks (e.g., Ethereum L2 or another EVM-compatible chain).

- **Front-End:**  
  Developed using **HTML, CSS, and JavaScript** for a smooth, interactive user experience.

- **Back-End & API Services:**  
  Powered by **Go (Golang)** for high-performance processing and seamless API integrations.

- **Integrations:**  
  Real-time data feeds from NGOs, supply chain providers, and IoT sensors.

## Codebase Overview

### Go Backend
- **Initialization & Utilities:**  
  Includes project setup, template management, and UUID generation.
- **Database:**  
  Manages data with an SQLite database.
- **Repositories:**  
  Handles data for receivers, meals, feedback, and donations.
- **Models:**  
  Defines data structures for receivers, donations, and feedback.
- **Handlers:**  
  Manages API endpoints for meals, receivers, feedback, and donations.
- **Router & Middleware:**  
  Configures HTTP routes and implements logging and error recovery.

### JavaScript Frontend
- **Deployment Scripts:**  
  Contains scripts for deploying the SuccorTrail smart contract.
- **Components:**  
  React components for interacting with the blockchain.
- **Static JS Files:**  
  Manages donor and receiver functionalities along with meal-finding features.
- **Web3 Integration:**  
  Initializes Web3 and the smart contract instance.

### Smart Contracts
- **Solidity Contract:**  
  Contains the core logic of the SuccorTrail smart contract.
- **ABI:**  
  Provides the interface for interacting with the smart contract.

### Templates & Configuration
- **HTML Templates:**  
  For the homepage, donor, receiver, and meal finder pages.
- **Configuration Files:**  
  Includes Go modules and Git ignore settings.

## Team
- **Project Lead & Blockchain Architect**  
- [**Kherld**](https://x.com/kh3rld) 
- [**Joseph Okumu:**](https://github.com/JosephOkumu) 
- [**Ouma Ouma:**](https://github.com/oumaoumag)

## Impact & Future Plans
- **Scalability & Collaborations:**  
  Expanding partnerships with global relief organizations to enhance data reliability and distribution efficiency.
- **Community-Driven Governance:**  
  Exploring token-based voting or DAO structures for decentralized decision-making.
- **Expansion Beyond Food Aid:**  
  Adapting the platform for broader humanitarian efforts like disaster relief, medical supply logistics, and emergency response.

## License
This project is licensed under the **MIT License**. For details, please refer to the [LICENSE](LICENSE) file.
