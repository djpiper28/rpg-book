import { Pencil } from "lucide-react";
import { type ReactNode } from "react";
import { P } from "@/components/typography/P";
import { electronDialog } from "@/lib/electron";
import { getLogger } from "@/lib/grpcClient/client";

interface Props {
  description: string;
  imagePath?: string;
  imageUrl?: string;
  setImagePath: (src: string) => void;
}

export function IconSelector(props: Readonly<Props>): ReactNode {
  const getSrc = (path: string): string => {
    if (path.startsWith("http") || path.startsWith("file://")) {
      return path;
    }

    return `file://${path}`;
  };

  const imageSrc = props.imagePath
    ? getSrc(props.imagePath)
    : props.imageUrl
      ? getSrc(props.imageUrl)
      : undefined;

  return (
    <button
      className="relative flex items-center justify-center cursor-pointer border-dotted border-2 border-gray-600 min-w-20 max-w-1/2 max-h-1/2 min-h-40 p-1 overflow-hidden"
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
      <Pencil className="absolute top-2 right-2 drop-shadow-lg drop-shadow-black" />
      {imageSrc ? (
        <img
          alt="User selected"
          className="h-full w-full object-contain"
          src={imageSrc}
        />
      ) : (
        <P>No icon selected</P>
      )}
    </button>
  );
}
