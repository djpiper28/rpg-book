import { H1 } from "@/components/typography/H1";
import { Menubar } from "@/components/ui/menubar";
import { Button } from "./components/ui/button";
import { Settings } from "lucide-react";

function App() {
  return (
    <div className="flex flex-col gap-3">
      <Menubar id="menu" className="flex flex-row gap-3 justify-between">
        <H1>RPG Book</H1>
        <Button aria-label="Open settings">
          <Settings /> Settings
        </Button>
      </Menubar>
    </div>
  );
}

export default App;
