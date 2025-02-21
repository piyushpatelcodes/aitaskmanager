import { FaTasks, FaRegLightbulb, FaRegClock, FaBrain } from "react-icons/fa";

const features = [
  {
    title: "AI Task Prioritization",
    description: "Gemini AI smartly prioritizes tasks based on deadlines, importance, and habits.",
    icon: <FaRegLightbulb className="text-sky-400 text-4xl mb-4" />,
  },
  {
    title: "Auto Task Suggestions",
    description: "AI analyzes your workflow and suggests new tasks to optimize productivity.",
    icon: <FaBrain className="text-purple-400 text-4xl mb-4" />,
  },
  {
    title: "Smart Reminders",
    description: "Get AI-powered smart reminders that notify you at the right moment.",
    icon: <FaRegClock className="text-yellow-400 text-4xl mb-4" />,
  },
  {
    title: "Collaborative Tasks",
    description: "Easily share, assign, and track team progress with real-time collaboration.",
    icon: <FaTasks className="text-green-400 text-4xl mb-4" />,
  },
];

export default function FeaturesSection() {
  return (
    <section className="py-16 bg-black text-white">
      <div className=" mx-auto px-6 flex flex-col items-center text-center ml-24">
        


        {/* Feature Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-8 w-full max-w-4xl">
          {features.map((feature, index) => (
            <div
              key={index}
              className="bg-white/5 border border-white/10 p-6 rounded-xl shadow-md flex flex-col items-center justify-center text-center 
              transition-all duration-300 hover:scale-105 hover:shadow-[0px_0px_15px_#ffffff]"
            >
              {feature.icon}
              <h3 className="text-xl font-semibold text-white mb-2">{feature.title}</h3>
              <p className="text-gray-400">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
