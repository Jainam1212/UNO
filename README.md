# Full-Stack Project Name

This is a UNO game which works in LAN. Just a public project who donot have games in offices or beginners who want to learn about golang advance concepts like mutex, channels.

---

## 🚀 Features

* **Fast Backend:** Built with **Go** for high-performance and concurrency.
* **Modern UI:** Responsive frontend using **React** and **Fundamental CSS**. (Frontend would be done by claude main project is based on go logic)
* **Real-time Updates:** Integrated **WebSockets** for instant state synchronization.
* **Component Library:** Uses **shadcn/ui** and **Radix UI** for accessible, high-quality components.

---

## 🛠️ Tech Stack

### Backend
* **Language:** Go (1.21+)
* **Framework:** FastHTTP
* **Communication:** REST API & WebSockets

### Frontend
* **Framework:** React (Vite / Next.js)
* **State Management:** Zustand
* **Styling:** Fundamental CSS

---

## 🏃 Getting Started

### Prerequisites
* **Go** (1.21+) installed.
* **Node.js** (v18+) and **npm/pnpm** installed.
* **MySQL** or your preferred database instance running.

### 1. Installation
Clone the repository and install dependencies for both layers:

```bash
# Clone the repo
git clone [https://github.com/Jainam1212/UNO.git](https://github.com/Jainam1212/UNO.git)
open cmd prompt -> ipconfig
open App.jsx in front dir and search lan ip used - format - 192.168.x.x
search it's occurance throughout your front dir and replace everywhere with ipconfig value of LAN

# Install Backend dependencies
cd to backend dir
go mod tidy

# Install Frontend dependencies
cd to frontend dir
npm install