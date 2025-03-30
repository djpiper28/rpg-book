import { H1 } from "@/components/typography/H1";
import { Menubar } from "@/components/ui/menubar";

function App() {
  return (
    <div className="flex flex-col gap-3">
      <Menubar id="menu">
        <H1>RPG Book</H1>
      </Menubar>
    </div>
  );
}

export default App;
