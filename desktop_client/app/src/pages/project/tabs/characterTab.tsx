import { Button, Table, TextInput } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useEffect, useState } from "react";
import { Modal } from "@/components/modal/modal";
import { projectClient } from "@/lib/grpcClient/client";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";

export function CharacterTab() {
  const [opened, { close, open }] = useDisclosure(false);
  const [characterName, setCharacterName] = useState("");

  useEffect(() => {
    setCharacterName("");
  }, [opened]);

  const { setError } = useGlobalErrorStore((x) => x);
  const projectHandle = useTabStore((x) => x.selectedTab);
  const projectStore = useProjectStore((x) => x);

  if (!projectHandle) {
    return "No project selected";
  }

  const thisProject = projectStore.getProject(projectHandle);

  return (
    <>
      <Modal close={close} opened={opened} title="Create Character">
        <div className="flex flex-col gap-3">
          <TextInput
            label="Character Name"
            onChange={(x) => {
              setCharacterName(x.target.value);
            }}
            placeholder="John Smith"
            required={true}
            value={characterName}
          />
          <Button
            onClick={() => {
              projectClient
                .createCharacter({
                  name: characterName,
                  project: projectHandle,
                })
                .then((resp) => {
                  console.log("Created a player", resp.response);
                  close();
                })
                .catch((error: unknown) => {
                  setError({
                    body: String(error),
                    title: "Cannot create character",
                  });
                });
            }}
          >
            Create
          </Button>
        </div>
      </Modal>
      <Button
        onClick={() => {
          open();
        }}
      >
        Create Character
      </Button>

      <Table variant="vertical">
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>Factions</Table.Th>
          </Table.Tr>
          {thisProject.project.characters.map((character) => (
            <Table.Tr key={character.handle?.id}>
              <Table.Th>{character.name}</Table.Th>
              <Table.Th>TODO: Change me</Table.Th>
            </Table.Tr>
          ))}
        </Table.Thead>
      </Table>
    </>
  );
}
