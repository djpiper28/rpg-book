import { Button, FileInput, TextInput } from "@mantine/core";
import { File as FileIcon, Plus, Text } from "lucide-react";
import { useEffect, useState } from "react";
import { H2 } from "@/components/typography/H2";
import { H3 } from "@/components/typography/H3";
import { DbExtension } from "@/lib/databaseTypes";
import { randomPlace } from "@/lib/random/randomProjectName";

function validateProjectName(rawName: string): string {
  const name = rawName.trim().normalize("NFC");

  if (name === "") {
    return "Please chose a name for your project";
  }

  return "";
}

export function CreateProjectPage() {
  const [projectName, setProjectName] = useState<string>(randomPlace());
  const [saveLocation, setSaveLocation] = useState<File | null>();

  useEffect(() => {
    if (validateProjectName(projectName) !== "") {
      setSaveLocation(undefined);
      return;
    }

    const filteredName = projectName.trim().normalize("NFC").slice(0, 30);
    setSaveLocation(new File([], `${filteredName}${DbExtension}`));
  }, [projectName]);

  return (
    <>
      <H2>Create A New Project</H2>
      <TextInput
        description="A cool, user friendly name for your project."
        error={validateProjectName(projectName)}
        label="Project Name"
        onChange={(ev) => {
          const name = ev.target.value;
          setProjectName(name);
        }}
        placeholder="Project name"
        rightSection={<Text />}
        value={projectName}
        withAsterisk={true}
      />

      <H3>Advanced Settings</H3>
      <FileInput
        description="Where to save the file, you usually don't need to touch this."
        error={saveLocation ? "" : "Please chose a location to save as"}
        label="Save Location"
        onChange={(f) => {
          setSaveLocation(f);
        }}
        placeholder="Save location"
        rightSection={<FileIcon />}
        value={saveLocation}
        withAsterisk={true}
      />

      <Button>
        <Plus />
        Create
      </Button>
    </>
  );
}
