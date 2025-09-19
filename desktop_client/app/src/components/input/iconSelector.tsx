import { uint8ArrayToBase64 } from "@/lib/utils/base64";
import { Pencil } from "lucide-react";
import { type ReactNode, useState } from "react";
import { P } from "@/components/typography/P";
import { electronDialog } from "@/lib/electron";
import { getLogger, getSystemClient } from "@/lib/grpcClient/client";

interface Props {
  description: string;
  imageB64: string | undefined;
  setImageB64: (src: string) => void;
}

export function IconSelector(props: Readonly<Props>): ReactNode {
  const [ext, setExt] = useState("jpg");

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
            .then(async (result: Electron.OpenDialogReturnValue) => {
              if (result.canceled) {
                return;
              }

              const resp = await getSystemClient().readFile({
                filepath: result.filePaths[0],
              });

              const b64 = uint8ArrayToBase64(resp.response.data);

              const extension =
                result.filePaths[0].split(".").pop()?.toLowerCase() ?? "png";

              setExt(extension);
              props.setImageB64(b64);
            })
            .catch((error: unknown) => {
              getLogger().error("Cannot get base64 for file", {
                error: String(error),
              });
            });
        }}
      >
        <Pencil className="absolute" />
        {props.imageB64 ? (
          <img
            alt="User selected"
            className="static max-w-1/2 max-h-1/2"
            src={`data:image/${ext};base64,${props.imageB64}`}
          />
        ) : (
          <P>No icon selected</P>
        )}
      </button>
    </div>
  );
}
