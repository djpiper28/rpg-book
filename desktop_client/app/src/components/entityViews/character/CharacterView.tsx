import { type ReactNode, useEffect, useState } from "react";
import { EditDelete } from "@/components/buttons/editDelete";
import MarkdownRenderer from "@/components/renderers/markdown";
import { H2 } from "@/components/typography/H2";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import { type BasicCharacterDetails } from "@/lib/grpcClient/pb/project_character";
import { uint8ArrayToBase64 } from "@/lib/utils/base64";
import { useTabStore } from "@/stores/tabStore";

export interface CharacterViewProps {
  characterId: string;
  isEditModalOpen: boolean;
  onDelete: () => void;
  onEdit: () => void;
}

export function CharacterView(props: CharacterViewProps): ReactNode {
  const [character, setCharacter] = useState<
    BasicCharacterDetails | undefined
  >();

  const [iconB64, setIconB64] = useState("");
  const projectHandle = useTabStore((x) => x.selectedTab);

  useEffect(() => {
    // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
    const sync = async () => {
      if (!props.characterId) {
        return;
      }

      if (!projectHandle) {
        return;
      }

      try {
        const resp = await getProjectClient().getCharacter({
          character: {
            id: props.characterId,
          },
          project: projectHandle,
        });

        setCharacter(resp.response.details);
        setIconB64(uint8ArrayToBase64(resp.response.icon));
      } catch (error: unknown) {
        getLogger().error("Cannot get character", {
          character: props.characterId,
          error: JSON.stringify(error),
          project: JSON.stringify(projectHandle),
        });
      }
    };

    // eslint-disable-next-line @typescript-eslint/no-floating-promises
    sync();
  }, [projectHandle, props.characterId, props.isEditModalOpen]);

  if (!character) {
    return <div className="flex-1" />;
  }

  return (
    <div className="flex-1 gap-3 overflow-y-auto flex flex-col">
      <div className="flex flex-col gap-3">
        <div className="flex flex-row gap-3 items-center justify-between">
          <H2>{character.name}</H2>
          <EditDelete delete={props.onDelete} edit={props.onEdit} />
        </div>
        {iconB64 && (
          <img
            alt="User selected"
            className="max-h-screen"
            src={`data:image/jpg;base64,${iconB64}`}
          />
        )}
      </div>
      <MarkdownRenderer markdown={character.description} />
    </div>
  );
}
