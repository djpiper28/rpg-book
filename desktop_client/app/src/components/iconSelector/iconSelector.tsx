import { Pencil } from "lucide-react";
import { useState } from "react";
import { electronDialog } from "@/lib/electron";

interface Props {
  description: string;
}

export function IconSelector(props: Readonly<Props>) {
  const [selectedFile, setSelectedFile] = useState("");

  return (
    <div className="w-20 border-dotted border-2 border-gray-600">
      <button
        className="z-0"
        onClick={() => {
          electronDialog
            .showOpenDialog({
              buttonLabel: "Ok",
              filters: [
                {
                  extensions: ["image.jpeg"],
                  name: props.description,
                },
              ],
              properties: ["openFile"],
              title: "Chose a project to open",
            })
            .then((result: Electron.OpenDialogReturnValue) => {
              if (result.canceled) {
                return;
              }

              setSelectedFile("file://" + result.filePaths[0]);
            })
            .catch(console.error);
        }}
      >
        <Pencil />
      </button>
      <img
        alt="User selected"
        className="absolute left-0 right-0"
        src={selectedFile}
      />
    </div>
  );
}
