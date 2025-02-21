"use client"; // ✅ Required for Next.js (since Rive runs on the client)

import React  from "react";
import  { useRive } from "rive-react";

export default function RiveAnimation() {
  const {  RiveComponent } = useRive({
    src: "/eclipse2.riv", // ✅ Path to your Rive file
    stateMachines: "State Machine 1", // ⚡ Name of the state machine (from Rive editor)
    autoplay: true, // ✅ Starts playing automatically
  });

  // 🎯 Example: Get an input to control animation states

  return (
    <div className="flex justify-center items-center h-screen">
      <RiveComponent className="w-screen h-screen" /> {/* Adjust size */}
    </div>
  );
}
