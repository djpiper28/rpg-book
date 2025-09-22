import { Pencil } from "lucide-react";
import { type ReactNode } from "react";
import { P } from "@/components/typography/P";
import { electronDialog } from "@/lib/electron";
import { getLogger } from "@/lib/grpcClient/client";

interface Props {
  description: string;
  imageDataB64?: string;
  imagePath?: string;
  setImagePath: (src: string) => void;
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
                  extensions: [
                    "jpeg",
                    "png",
                    "jpg",
                    "webp",
                    "gif",
                    "bmp",
                    "svg",
                    "tiff",
                    "ccitt",
                  ],
                  name: props.description,
                },
              ],
              properties: ["openFile"],
              title: "Chose a project to open",
            })
            .then((result: Electron.OpenDialogReturnValue): void => {
              if (result.canceled) {
                return;
              }

              props.setImagePath(result.filePaths[0]);
            })
            .catch((error: unknown) => {
              getLogger().error("Cannot get base64 for file", {
                error: String(error),
              });
            });
        }}
      >
        <Pencil className="absolute" />
        {props.imagePath ? (
          <img
            alt="User selected"
            className="static max-w-1/2 max-h-1/2"
            src={`file://${props.imagePath}`}
          />
        ) : props.imageDataB64 ? (
          <img
            alt="User selected"
            className="static max-w-1/2 max-h-1/2"
            src={`data:image/jpg;base64,${props.imageDataB64}`}
          />
        ) : (
          <P>No icon selected</P>
        )}
      </button>
    </div>
  );
}
