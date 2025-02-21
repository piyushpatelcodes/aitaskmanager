# 🚀 Next.js Project - Ai Todo List

## 📌 Overview
This is a **Next.js** project built with **React, TypeScript, and Clerk authentication**. It includes:
- ✅ **Server-side and Client-side rendering**
- ✅ **Authentication with Clerk**
- ✅ **API Integration**
- ✅ **Responsive UI with Tailwind CSS**

   <img width="947" alt="Screenshot 2025-02-20 090739" src="https://github.com/user-attachments/assets/54e38617-8ef3-4882-97da-22b0cb33c41b" />


 
## 🛠 Tech Stack
- **Next.js** - React framework for full-stack applications
- **React** - UI Components
- **TypeScript** - Static Typing
- **Tailwind CSS** - Styling
- **Clerk** - Authentication
- **Axios** - API requests

---

## 🚀 Getting Started

### 1️⃣ Clone the Repository

git clone [https://github.com/amitpal0728/AI-TaskManager-Master.git]
cd my-next-app
----
2️⃣ Install Dependencies
 
npm install
# OR
yarn install
3️⃣ Setup Environment Variables


Create a .env.local file and add:

---
NEXT_PUBLIC_BACKEND_URL=https://your-api.com
CLERK_PUBLISHABLE_KEY=your_clerk_key
CLERK_SECRET_KEY=your_clerk_secret
---
4️⃣ Run the Development Server

npm run dev
# OR
yarn dev
Now open http://localhost:3000 in your browser. 🚀

🔑 Authentication with Clerk
auth() is used in Server Components
Client Components receive authData as props
Session management handled automatically
📌 API Routes
Method	Route	Description
GET	/api/tasks	Fetch user tasks
POST	/api/tasks	Create a new task
PATCH	/api/tasks/:id	Update task
DELETE	/api/tasks/:id	Delete task
📦 Building for Production
bash
Copy
Edit
npm run build
npm run start
This optimizes the Next.js app for best performance in production.

🛠 Deployment
Deploy on Vercel with one click:

Alternatively, deploy on Netlify, AWS, or Docker.

✨ Contributing
Feel free to fork this repo, open issues, or create PRs to improve the project.

📌 Maintainer: Amit Pal amitpal0728@gmail.com

📜 License
This project is MIT licensed. You are free to use, modify, and distribute it. 🚀

