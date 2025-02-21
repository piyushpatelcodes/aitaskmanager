import Rivecomponent from "../components/rivecomponent";
import ShinyText from "../components/shinytext";
import Featurecard from "../components/feature"
import TextPressure from "../components/pressuretext";

export default function Home() {
  return (
    <>
      <section className="relative w-full h-screen bg-black">
        <Rivecomponent />
      </section>

      <div className="bg-black" style={{position: 'relative', height: '350px'}}>
  <TextPressure
    text="Todo.ai!"
    flex={true}
    alpha={false}
    stroke={false}
    width={true}
    weight={true}
    italic={true}
    textColor="#ffffff"
    strokeColor="#ff0000"
    minFontSize={10}
  />
</div>

      <section className="py-16 bg-black text-white">
        <div className="max-w-6xl mx-auto px-6 items-center justify-center text-center">
          <ShinyText
            text="Our Features"
            disabled={false}
            speed={3}
            className="text-4xl font-bold items-center text-center mb-11"
          />
          {/* <h2 className="text-3xl font-bold text-center mb-8">Our Features</h2> */}

          {/* Feature Grid */}
          <div className="flex gap-4 items-center justify-between">

            <Featurecard />
          </div>
        </div>
      </section>
    
    </>
  );
}
