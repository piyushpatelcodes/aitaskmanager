"use client";

import React from "react";
import { Button } from "../components/ui/button";
import { SignInButton, SignedIn, SignedOut, UserButton } from "@clerk/nextjs";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { Moon } from "lucide-react";


const Header = () => {
  const router = useRouter();

  return (
    <div className="fixed top-0 w-full h-16 flex items-center justify-center backdrop-blur-xl bg-black/20 dark:bg-white/10 shadow-md border-b border-white/10 dark:border-black/20 z-50 transition-all">
      {/* Background Blob */}
      <div className="absolute inset-0 flex justify-center items-center -z-10">
        <div className="bg-blue-500 w-96 h-32 rounded-full blur-3xl opacity-40"></div>
      </div>

      {/* Navigation */}
      <div className="flex gap-6 text-white  text-lg font-semibold">
      <Image src="/vercel.svg" className="ml-4" height={20} width={20} alt="logo" />
      
      {["Home", "Dashboard", "About"].map((tab) => (
        <button
          key={tab}
          onClick={() => router.push(tab === "Home" ? "/" : `/${tab.toLowerCase()}`)}
          className="px-4 py-1 transition-all rounded-md hover:bg-white/10 hover:text-green-400"
        >
          {tab}
        </button>
      ))}
    </div>

      {/* Right Section */}
      <div className="ml-auto flex items-center gap-4 mr-6">
        {/* Light/Dark Mode Toggle */}
        <Moon className="stroke-sky-400 drop-shadow-[0_0_9px_#38bdf8] text-lg" />
  
       
        
        {/* Auth Buttons */}
        <SignedOut>
          <Button variant="outline">
            <SignInButton />
          </Button>
        </SignedOut>
        <SignedIn>
          <UserButton />
        </SignedIn>
      </div>
    </div>
  );
};

export default Header;
