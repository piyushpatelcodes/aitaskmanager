"use client";

import {  useAuth, useUser } from "@clerk/nextjs";
import { useEffect, useState } from "react";
import axios from "axios";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Dialog, DialogTrigger, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { PencilIcon, Trash2Icon } from "lucide-react";
import { toast } from "sonner";
import Tiles from "@/components/tiles";
import Loadingcomponent from "../../components/ui/loading"

const getPriorityStyles = (priority:unknown) => {
  switch (priority) {
    case "High":
      return {
        border: "border-red-500",
        text: "text-red-600 bg-red-100",
      };
    case "Medium":
      return {
        border: "border-yellow-500",
        text: "text-yellow-600 bg-yellow-100",
      };
    case "Low":
      return {
        border: "border-green-500",
        text: "text-green-600 bg-green-100",
      };
    default:
      return {
        border: "border-gray-300",
        text: "text-gray-600 bg-gray-100",
      };
  }
};

interface Task {
  priority: string;
  id: string;
  title: string;
  description: string;
  completed: boolean;
}

export default function Dashboard() {
  const { userId } = useAuth();
  const { user } = useUser();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [taskTitle, setTaskTitle] = useState(""); // ✅ Added missing state
  const [updateData, setUpdateData] = useState<Partial<Task>>({}); // ✅ Fix: Partial<Task>
  const [currentTaskId, setCurrentTaskId] = useState("");
  const [loading,setloading] = useState(false);

  // Fetch Tasks
  useEffect(() => {
    fetchTasks();
    const message = localStorage.getItem("toastMessage");
    if (message) {
      toast.success(message);  // ✅ Show toast after refresh
      localStorage.removeItem("toastMessage"); // ✅ Prevent duplicate toasts
    }
  },[userId]);
  
  const fetchTasks = async () => {
    if (!userId) return;
    try {
      setloading(true);
      const res = await axios.get(`${process.env.NEXT_PUBLIC_BACKEND_URL}/tasks`, {
        headers: { "X-User-ID": userId },
      });
      setTasks(res.data);
      const message = localStorage.getItem("toastMessage");
      if (message) {
        toast.success(message);  // ✅ Show toast after refresh
        localStorage.removeItem("toastMessage"); // ✅ Prevent duplicate toasts
      }
      setloading(false)
    } catch (error) {
      setloading(false)
      console.error("Error fetching tasks:", error);
    }finally{
      setloading(false)
    }
  };

  const capitalizeFirstLetter = (str:string) => {
    return str.charAt(0).toUpperCase() + str.slice(1);
  };
  

  // Add New Task
  const addTask = async () => {
    if (!taskTitle) return;
    try {
      setloading(true)
      await axios.post(
        `${process.env.NEXT_PUBLIC_BACKEND_URL}/tasks`,
        { title: taskTitle },
        { headers: { "X-User-ID": userId } }
      );
      setTaskTitle("");
      setloading(false)
      fetchTasks(); // ✅ Refresh
      localStorage.setItem("toastMessage", "Task has been created with AI Description and Priority.");
    } catch (error) {
      setloading(false)
      console.error("Error adding task", error);
    }finally{
      setloading(false)
    }
  };

  // Update Task
  const updateTask = async () => {
    if (!currentTaskId) return;
    try {
      setloading(true)
      await axios.put(
        `${process.env.NEXT_PUBLIC_BACKEND_URL}/tasks/${currentTaskId}`,
        updateData,
        { headers: { "X-User-ID": userId } }
      );
      setUpdateData({});
      setCurrentTaskId("");
      localStorage.setItem("toastMessage", "Task Updated.");
      setloading(false)
      fetchTasks(); // ✅ Refresh
    } catch (error) {
      console.error("Error updating task", error);
    }finally{
      setloading(false)
    }
  };

  // Delete Task
  const deleteTask = async (id: string) => {
    console.log("task id: ", id)
    await axios.delete(`${process.env.NEXT_PUBLIC_BACKEND_URL}/tasks/${id}`, {
      headers: { "X-User-ID": userId },
    });
    localStorage.setItem("toastMessage", "Task has been Deleted.");
    toast.success("Task has been Deleted")
    setTasks(tasks.filter((task) => task.id !== id));
  };

  return (
    <div className="p-6 max-w-5xl mx-auto bg-black">
       {loading && <Loadingcomponent />}
        <div className="absolute inset-0 opacity-25 -z-10">
        <Tiles
          speed={0.5}
          squareSize={40}
          direction="diagonal"
          borderColor="#444" // Darker for subtle effect
          hoverFillColor="#222"
        />
      </div>
     
      {/* Header */}
      <h1 className="text-3xl font-bold mb-4">Dashboard</h1>
      <p className="text-gray-500">Welcome, {user?.fullName}!</p>

      {/* Add Task Input */}
      <div className="flex gap-3 mt-4">
        <Input
        className="text-white"
          placeholder="Enter task title..."
          value={taskTitle}
          onChange={(e) => setTaskTitle(e.target.value)}
        />
        <Button onClick={addTask}>Add Task</Button>
      </div>

      {/* Tasks Section */}
      <h2 className="text-xl font-semibold mt-6 mb-3">Your Tasks</h2>
      <div className="grid gap-4 md:grid-cols-2">
      {tasks?.map((task) => {
          const styles = getPriorityStyles(task.priority);
          return (
            <Card key={task.id} className={`p-4 border-2 ${styles.border} shadow-lg w-full`}>
              <CardContent>
                <div className="flex justify-between items-center">
                  <h3 className="font-semibold">{capitalizeFirstLetter(task.title)}</h3>
                  {/* Priority Pill */}
                  <span className={`px-3 py-1 text-xs font-semibold rounded-full ${styles.text}`}>
                    {task.priority}
                  </span>
                </div>
                <p className="text-gray-600 mt-2">{task.description}</p>

                {/* Action Buttons */}
                <div className="flex justify-end  mt-4">
                <Dialog>
                  <DialogTrigger asChild>
                    <Button
                    className="text-blue-600 hover:text-blue-800 bg-black hover:bg-gray-400/10"
                      onClick={() => {
                        setUpdateData({ title: task.title, description: task.description });
                        setCurrentTaskId(task.id);
                      }}
                    >
                       <PencilIcon size={18} />
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="p-4">
                  <DialogTitle>Are you absolutely sure?</DialogTitle>
                    <h3 className="text-lg font-bold">Update Task</h3>
                    <Input
                      value={updateData.title || ""}
                      onChange={(e) => setUpdateData({ ...updateData, title: e.target.value })}
                      placeholder="Title"
                    />
                    <Input
                      value={updateData.description || ""}
                      onChange={(e) => setUpdateData({ ...updateData, description: e.target.value })}
                      placeholder="Description"
                    />
                    <Button onClick={updateTask}>Update</Button>
                  </DialogContent>
                </Dialog>
                <Button className="text-red-600 hover:text-red-800 bg-black hover:bg-gray-400/10" variant="destructive" onClick={() => deleteTask(task.id)}>
                <Trash2Icon size={18} />
                </Button>


                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>
    </div>
  );
}
