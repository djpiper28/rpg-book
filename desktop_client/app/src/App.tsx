import { H1 } from "@/components/typography/H1";
import { Settings } from "lucide-react";
import { logger } from "./lib/grpcClient/client";
import { Button } from "@mantine/core";

function App() {
  logger.warn("render call", {});
  return (
    <div className="flex flex-col gap-3">
      <div id="menu" className="flex flex-row gap-3 justify-between">
        <H1>RPG Book</H1>
        <Button aria-label="Open settings" variant>
          <Settings /> Settings
        </Button>
      </div>
    </div>
  );
}

export default App;
