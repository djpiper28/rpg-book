import { Button, Table, TextInput } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useEffect, useState } from "react";
import { Modal } from "@/components/modal/modal";
import { projectClient } from "@/lib/grpcClient/client";
import { useTabStore } from "@/stores/tabStore";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";

export function CharacterTab() {
  const [opened, { close, open }] = useDisclosure(false);
  const [name, setName] = useState("");
  const project = useTabStore((x) => x.selectedTab);
  const { setError } = useGlobalErrorStore((x) => x);

  useEffect(() => {
    setName("");
  }, [project, opened]);

  return (
    <>
      <Modal close={close} opened={opened} title="Create Character">
        <TextInput
          label="Character Name"
          onChange={(e) => {
            setName(e.target.value);
          }}
          value={name}
        />
        <Button
          onClick={() => {
            projectClient
              .createCharacter({
                name,
                project,
              })
              .then((x) => {
                // TODO: something
              })
              .catch(setError);
          }}
        >
          Create
        </Button>
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
        </Table.Thead>
      </Table>
    </>
  );
}
