import { Button } from "@mantine/core";
import { type ReactNode, useEffect, useState } from "react";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import {
  type BasicCharacterDetails,
  type CharacterHandle,
} from "@/lib/grpcClient/pb/project_character";
import { uint8ArrayToBase64 } from "@/lib/utils/base64";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";
import { CharacterEdit } from "./CharacterEdit";

interface Props {
  characterHandle: CharacterHandle;
  closeDialog: () => void;
}

export default function EditCharacterModal(props: Readonly<Props>): ReactNode {
  const [characterDetails, setCharacterDetails] =
    useState<BasicCharacterDetails>({
      description: "",
      handle: {
        id: "",
      },
      name: "",
    });

  const projectStore = useProjectStore((x) => x);
  const [icon, setIcon] = useState("");
  const [iconPath, setIconPath] = useState<string>("");
  const [dirtyIcon, setDirtyIcon] = useState(false);
  const { setError } = useGlobalErrorStore((x) => x);
  const projectHandle = useTabStore((x) => x.selectedTab);

  useEffect(() => {
    if (!projectHandle) {
      return;
    }

    getProjectClient()
      .getCharacter({
        character: props.characterHandle,
        project: projectHandle,
      })
      .then((resp) => {
        if (!resp.response.details) {
          return;
        }

        setCharacterDetails(resp.response.details);
        setIcon(uint8ArrayToBase64(resp.response.icon));
        setDirtyIcon(false);
      })
      .catch((error: unknown) => {
        setError({
          body: String(error),
        });
      });
  }, [props.characterHandle, projectHandle, setError]);

  if (!projectHandle) {
    return "No project selected";
  }

  return (
    <div className="flex flex-col gap-3">
      <CharacterEdit
        description={characterDetails.description}
        iconPath={iconPath}
        imageDataB64={icon}
        name={characterDetails.name}
        onDescriptionChange={(description) => {
          setCharacterDetails((prev) => ({ ...prev, description }));
        }}
        onNameChange={(name) => {
          setCharacterDetails((prev) => ({ ...prev, name }));
        }}
        setIconPath={(path) => {
          setIconPath(path);
          setDirtyIcon(true);
        }}
      />
      <Button
        onClick={() => {
          const f = async (): Promise<void> => {
            try {
              await getProjectClient().updateCharacter({
                details: characterDetails,
                handle: props.characterHandle,
                iconPath,
                project: projectHandle,
                setIcon: dirtyIcon,
              });

              projectStore.addCharacter(projectHandle, characterDetails);

              props.closeDialog();
            } catch (error: unknown) {
              getLogger().error("Cannot edit character", {
                error: String(error),
              });

              setError({
                body: String(error),
              });
            }
          };

          // eslint-disable-next-line @typescript-eslint/no-floating-promises
          f();
        }}
      >
        Save
      </Button>
    </div>
  );
}
