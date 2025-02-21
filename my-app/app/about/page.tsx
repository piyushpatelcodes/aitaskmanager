
"use client"
import BlurText from "../../components/blurtext";
import Tiles from "../../components/tiles"; // Adjust import if needed
import ShinyText from "@/components/shinytext";

const About = () => {
    const handleAnimationComplete = () => {
        console.log('blur text Animation completed!');
      };
  return (
    <div className="relative mt-10 w-full min-h-screen bg-black text-white overflow-hidden">
      
      {/* Tile Background */}
      <div className="absolute inset-0 opacity-45">
        <Tiles
          speed={0.5}
          squareSize={40}
          direction="diagonal"
          borderColor="#444" // Darker for subtle effect
          hoverFillColor="#222"
        />
      </div>

      {/* About Section */}
      <section className="relative max-w-5xl mx-auto px-6 py-20 text-center">
        
        {/* Section Title */}
        <h1>
        ABOUT
        </h1>
     
        <ShinyText
            text="Todo.Ai"
            disabled={false}
            speed={3}
            className="text-5xl font-bold items-center text-center mb-11"
          />
          <BlurText
                  text="Isn't this so cool?!"
                  delay={150}
                  animateBy="words"
                  direction="top"
                  onAnimationComplete={handleAnimationComplete}
                  className="text-2xl font-bold text-zinc-600 mb-8" animationFrom={undefined} animationTo={undefined}/>
        

        {/* Description */}
        <p className="text-gray-300 text-lg max-w-3xl mx-auto mb-12">
          Our Gemini-powered task manager revolutionizes productivity. With AI-driven task prioritization, smart reminders, and workflow optimization, 
          it helps you focus on what truly matters. Experience the future of task management, powered by intelligent automation and real-time insights.
        </p>

        {/* Feature Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mt-10">
          <div className="bg-white/5 border  p-6 rounded-xl  backdrop-blur-md transition-all duration-300  ">
            <h3 className="text-xl font-semibold text-white">🚀 AI-Driven Task Management</h3>
            <p className="text-gray-400 mt-2">
              Automate your workflow with intelligent task suggestions and prioritization.
            </p>
          </div>

          <div className="  border  p-6 rounded-xl  backdrop-blur-md transition-all duration-300  ">
            <h3 className="text-xl font-semibold text-white">🔔 Smart Reminders & Notifications</h3>
            <p className="text-gray-400 mt-2">
              AI-powered notifications ensure you never miss an important deadline.
            </p>
          </div>
        </div>
      </section>

    </div>
  );
};

export default About;
