import { Pencil } from "lucide-react";
import { type ReactNode } from "react";
import { P } from "@/components/typography/P";
import { electronDialog } from "@/lib/electron";

interface Props {
  description: string;
  setSrc: (src: string) => void;
  src: string | undefined;
}

export function IconSelector(props: Readonly<Props>): ReactNode {
  return (
    <div className="border-dotted border-2 border-gray-600 min-w-20 max-w-1/2 max-h-1/2 min-h-20 p-1">
      <button
        className="cursor-pointer"
        onClick={() => {
          electronDialog
            .showOpenDialog({
              buttonLabel: "Ok",
              filters: [
                {
                  extensions: ["jpeg", "png", "jpg", "webp", "gif"],
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

              props.setSrc("file://" + result.filePaths[0]);
            })
            .catch(console.error);
        }}
      >
        <Pencil className="absolute" />
        {props.src ? (
          <img alt="User selected" className="static" src={props.src} />
        ) : (
          <P>No icon selected</P>
        )}
      </button>
    </div>
  );
}
